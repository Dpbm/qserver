import os
import sys
import pika
from utils import DB, Plugin
from utils.types import port_to_int, Statuses
from utils.exceptions import CanceledJob


# pylint: disable=too-many-locals
def callback(ch, method, body, db_instance):
    """
    Handles the incoming data from queue
    """
    try:
        job_id = body.decode()

        print(f"Processing job {job_id}")

        data = db_instance.get_job_data(job_id)

        result_types = data["selected_result_types"]
        qasm_file = data["qasm"]
        target_backend = data["target_simulator"]
        metadata = data["metadata"]
        status = data["status"]

        if(status == Statuses.CANCELED):
            raise CanceledJob()

        db_instance.update_job_start_time_to_now(job_id)
        db_instance.update_job_status("running", job_id)

        # the plugin name is first checked by the api to see if it's official
        # however, the user may try to bypass that
        # so be aware with potential threads here
        row = db_instance.get_plugin(target_backend)

        if len(row) != 1:
            raise ValueError("Failed on get plugin Name")

        plugin_name = row[0]
        plugin = Plugin(plugin_name)

        for result_type, active in result_types.items():
            if not active:
                continue

            print(f"executing for {result_type} results")
            results = plugin.run(target_backend, qasm_file, metadata, result_type)

            print("Saving results...")
            db_instance.save_results(result_type, results, job_id)

        db_instance.update_job_finish_time_to_now(job_id)
        db_instance.update_job_status("finished", job_id)

    # pylint: disable=broad-exception-caught
    except Exception as error:
        db_instance.update_job_finish_time_to_now(job_id)
        db_instance.update_job_status("failed", job_id)
        print(f"failed on worker callback: {str(error)}")

    finally:
        ch.basic_ack(delivery_tag=method.delivery_tag)


if __name__ == "__main__":
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

    if None in (
        rabbitmq_host,
        rabbitmq_port,
        rabbitmq_queue_name,
        db_host,
        db_port,
        db_name,
        db_user,
        db_password,
    ):
        print("Invalid environment variables!")
        sys.exit(1)

    credentials = pika.PlainCredentials(rabbitmq_user, rabbitmq_password)
    # pylint: disable=invalid-name
    connection = None

    print("Waiting for connection...")
    while connection is None:
        try:
            connection = pika.BlockingConnection(
                pika.ConnectionParameters(host=rabbitmq_host)
            )
        # pylint: disable=broad-exception-caught
        except Exception:
            pass
    print("Connection Stablished!")

    channel = connection.channel()

    channel.queue_declare(queue=rabbitmq_queue_name, durable=True)
    print("Waiting for jobs...")

    db = DB(
        host=db_host,  # type: ignore
        port=db_port,  # type: ignore
        db_name=db_name,  # type: ignore
        user=db_user,  # type: ignore
        password=db_password,  # type: ignore
    )

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
