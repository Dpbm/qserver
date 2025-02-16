import pika, os
from utils.db import DB
from utils.plugin import Plugin


def callback(ch, method, body,db):
    try:
        job_id = body.decode()

        print(f"Processing job {job_id}")

        data = db.get_job_data(job_id)
        
        result_types = data["selected_result_types"]
        qasm_file = data["qasm"]
        target_backend = data["target_simulator"]

        db.update_job_status('running', job_id)
        db.update_job_start_time_to_now(job_id)


        # get plugin (be aware that once the plugin name can be passed by the user, he may try to bypass and run arbitrary code)
        plugin_name = db.get_plugin(target_backend)
        # download plugin if it isn't installed
        plugin = Plugin(plugin_name)

        for result_type, active in result_types.items():
            if(not active):
                continue

            print(f"executing for {result_type} results")
            results = plugin.run(qasm_file, result_type)

            print("Saving results...")
            db.save_results(result_type, results, job_id)

        db.update_job_status('finished', job_id)

    except Exception as error:
        db.update_job_status('failed', job_id)
        print(f"failed on worker callback: {str(error)}")

    finally:
        db.update_job_finish_time_to_now(job_id)
        ch.basic_ack(delivery_tag=method.delivery_tag)

if __name__ == '__main__':
    rabbitmq_host = os.getenv("RABBITMQ_HOST")
    rabbitmq_port = os.getenv("RABBITMQ_PORT")
    rabbitmq_queue_name = os.getenv("RABBITMQ_QUEUE_NAME")

    db_host = os.getenv("DB_HOST")
    db_port = os.getenv("DB_PORT")
    db_name = os.getenv("DB_NAME")
    db_user = os.getenv("DB_USER")
    db_password = os.getenv("DB_PASSWORD")

    credentials = pika.PlainCredentials('guest', 'guest')
    connection = None

    print("Waiting for connection...")
    while(connection is None):
        try:
            connection = pika.BlockingConnection(pika.ConnectionParameters(host=rabbitmq_host))
        except Exception:
            pass
    print("Connection Stablished!")

    channel = connection.channel()

    channel.queue_declare(queue=rabbitmq_queue_name, durable=True)
    print("Waiting for jobs...")

    db = DB(
        host=db_host, 
        port=db_port, 
        db_name=db_name, 
        user=db_user, 
        password=db_password)

    channel.basic_qos(prefetch_count=1) # ensure that a single message is passed to each idle worker
    channel.basic_consume(
        queue=rabbitmq_queue_name, 
        on_message_callback=lambda ch, method, properties, body: callback(ch, method, body, db)
        )

    channel.start_consuming()
    db.close()
