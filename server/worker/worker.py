import os
import sys
import datetime
import uuid
import logging
import pika
from utils import DB, Plugin
from utils.types import port_to_int, Statuses
from utils.exceptions import (
    CanceledJob,
    IdNotFound,
    InvalidResultTypes,
    InvalidQasmFile,
    InvalidBackend,
    InvalidStatus,
)
from utils.checks import (
    valid_result_types,
    valid_data_for_id,
    valid_qasm,
    valid_backend,
)
from utils.log_files import create_path

logger = logging.getLogger(__name__)
logging.getLogger("pika").setLevel(logging.WARNING)


# pylint: disable=too-many-locals,too-many-branches,too-many-statements
def callback(ch, method, body, db_instance):
    """
    Handles the incoming data from queue
    """
    try:
        job_id = body.decode()

        logger.debug("Processing job %s", job_id)

        data = db_instance.get_job_data(job_id)
        if not valid_data_for_id(data):
            raise IdNotFound(data)

        logger.debug("Got data: %s", data)

        result_types = data["selected_result_types"]
        if not valid_result_types(result_types):
            raise InvalidResultTypes(result_types)

        qasm_file = data["qasm"]
        if not valid_qasm(qasm_file):
            raise InvalidQasmFile(qasm_file)

        target_backend = data["target_simulator"]
        if not valid_backend(target_backend):
            raise InvalidBackend(target_backend)

        metadata = {}
        if data.get("metadata") is not None:
            metadata = data["metadata"]

        status = data["status"]
        if status == Statuses.CANCELED.value:
            raise CanceledJob(job_id)

        if status != Statuses.PENDING.value:
            raise InvalidStatus(status)

        db_instance.update_job_start_time_to_now(job_id)
        db_instance.update_job_status(Statuses.RUNNING.value, job_id)

        # the plugin name is first checked by the api to see if it's official
        # however, the user may try to bypass that
        # so be aware with potential threads here
        row = db_instance.get_plugin(target_backend)
        logger.debug("Got plugin row:  %s", row)
        if len(row) != 1:
            raise ValueError("Failed on get plugin Name")

        plugin_name = row["plugin"]
        logger.debug("using plugin: %s", plugin_name)

        plugin = Plugin(plugin_name)

        logger.debug("mapping through result types: %s", result_types.items())
        for result_type, active in result_types.items():
            if not active:
                continue

            logger.debug("executing for %s results", result_type)
            results = plugin.run(target_backend, qasm_file, metadata, result_type)

            logger.debug("Saving results...")
            db_instance.save_results(result_type, results, job_id)

        db_instance.update_job_finish_time_to_now(job_id)
        db_instance.update_job_status(Statuses.FINISHED.value, job_id)

    except IdNotFound as error:
        logger.error("Job Id Not Found")
        logger.error(str(error))

    except InvalidStatus as error:
        logger.error("Job was already executed")
        logger.error(str(error))

    except CanceledJob:
        db_instance.update_job_finish_time_to_now(job_id)
        logger.warning("Job Was Canceled")

    # pylint: disable=broad-exception-caught
    except Exception as error:
        db_instance.update_job_finish_time_to_now(job_id)
        db_instance.update_job_status(Statuses.FAILED.value, job_id)
        logger.error("failed on worker callback: %s", str(error))

    finally:
        ch.basic_ack(delivery_tag=method.delivery_tag)


if __name__ == "__main__":

    logs_path = os.getenv("LOGS_PATH")
    if logs_path:
        filename = os.path.join(
            logs_path, f"{str(datetime.datetime.now())}-{str(uuid.uuid4())}.log"
        )
        create_path(filename)
        logging.basicConfig(level=logging.DEBUG, filename=filename)
    else:
        logging.basicConfig(level=logging.DEBUG)

    rabbitmq_host = os.getenv("RABBITMQ_HOST")
    rabbitmq_port = port_to_int(os.getenv("RABBITMQ_PORT"))
    rabbitmq_queue_name = os.getenv("RABBITMQ_QUEUE_NAME")
    rabbitmq_user = os.getenv("RABBITMQ_USER")
    rabbitmq_password = os.getenv("RABBITMQ_PASSWORD")

    db_host = os.getenv("DB_HOST")
    db_port = port_to_int(os.getenv("DB_PORT"))
    db_name = os.getenv("DB_NAME")
    db_user = os.getenv("DB_USERNAME")
    db_password = os.getenv("DB_PASSWORD")

    variables = (
        rabbitmq_host,
        rabbitmq_port,
        rabbitmq_queue_name,
        db_host,
        db_port,
        db_name,
        db_user,
        db_password,
    )

    if None in variables:
        logger.error("Invalid environment variables!: %s", str(variables))
        sys.exit(1)

    credentials = pika.PlainCredentials(rabbitmq_user, rabbitmq_password)
    # pylint: disable=invalid-name
    connection = None

    logger.debug("Waiting for connection...")
    while connection is None:
        try:
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(host=rabbitmq_host,  port=rabbitmq_port, credentials=pika.credentials.PlainCredentials(username=rabbitmq_user, password=rabbitmq_password))
            )
        # pylint: disable=broad-exception-caught
        except Exception:
            pass
    logger.debug("Connection Stablished!")

    channel = connection.channel()

    channel.queue_declare(queue=rabbitmq_queue_name, durable=True)
    logger.debug("Waiting for jobs...")

    db = None
    while db is None:
        try:
            db = DB(
                host=db_host,  # type: ignore
                port=db_port,  # type: ignore
                db_name=db_name,  # type: ignore
                user=db_user,  # type: ignore
                password=db_password,  # type: ignore
            )

        # pylint: disable=broad-exception-caught
        except Exception:
            pass
    logger.debug("Connected to DB!")

    channel.basic_qos(
        prefetch_count=1
    )  # ensure that a single message is passed to each idle worker
    channel.basic_consume(
        queue=rabbitmq_queue_name,
        on_message_callback=lambda ch, method, properties, body: callback(
            ch, method, body, db
        ),
    )

    channel.start_consuming()
    db.close()
