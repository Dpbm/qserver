from typing import Tuple, Any
import psycopg2 as pg
import logging

class DB:
    def __init__(self, host:str, port:str, db_name:str, user:str, password:str):
        self._connection = pg.connect(f"postgres://{user}:{password}@{host}:{port}/{db_name}")
        self._cursor = self._connection.cursor()
        logging.info("Initialized DB")

    # update job status
    # save results
    # update start_time
    # update finish_time

    def get_job_data(self, job_id:str) -> Tuple[Any]:
        self._cursor.execute("""
            SELECT 
                job_data.*, 
                (
                        SELECT row_to_json(data)
                        FROM (
                            SELECT counts, quasi_dist, expval
                            FROM result_types AS rt
                            WHERE rt.job_id = job_data.id
                        ) data
                ) AS selected_result_types
            FROM 
                jobs AS job_data
            WHERE
                job_data.id = %s
        """, (job_id,)) # get all data from job and create a json with the result types the user selected
        return self._cursor.fetchone()


    def _commit(self):
        """
            Make database changes definitive.
        """
        self._connection.commit()

    def _close(self):
        """
            Finish database connection.
        """
        self._cursor.close()
        self._connection.close()
