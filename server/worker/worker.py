import pika, os
from utils.db import DB


def callback(ch, method, body,db):
    job_id = body.decode()
    print(f"Processing job {job_id}")

    data = db.get_job_data(job_id)
    print(data)
    # get qasm (maybe)
    # update status to running
    # update start_time
    # get all results types the user wants
    # get the requested plugin
    # execute the job for each result type
    # store the results in memory (or disk)
    # save results on db
    # update status
    # update finish_time
    # note: catch errors, if an error is raised, update the status as well


    ch.basic_ack(delivery_tag=method.delivery_tag)


if __name__ == '__main__':
    host = os.getenv("RABBITMQ_HOST")
    queue_name = os.getenv("RABBITMQ_QUEUE_NAME")

    connection = pika.BlockingConnection(pika.ConnectionParameters(host=host))
    channel = connection.channel()

    channel.queue_declare(queue=queue_name, durable=True)
    print("Waiting for jobs...")

    db = DB(host="0.0.0.0", port="5432", db_name="quantum", user="test", password="test")

    channel.basic_qos(prefetch_count=1) # ensure that a single message is passed to each idle worker
    channel.basic_consume(
        queue=queue_name, 
        on_message_callback=lambda ch, method, properties, body: callback(ch, method, body, db)
        )

    channel.start_consuming()