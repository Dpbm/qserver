from typing import Tuple, Any
import psycopg2 as pg
import psycopg2.extras as pgextras
import logging
from datetime import datetime
from .types import Results


class DB:
    def __init__(self, host:str, port:str, db_name:str, user:str, password:str):
        self._connection = pg.connect(f"postgres://{user}:{password}@{host}:{port}/{db_name}")
        # cursor_factory: https://www.geeksforgeeks.org/psycopg2-return-dictionary-like-values/
        self._cursor = self._connection.cursor(cursor_factory = pgextras.RealDictCursor)
        logging.info("Initialized DB")

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

    def update_job_status(self, status:str, job_id:str):
        self._cursor.execute("UPDATE jobs SET status=%s WHERE id=%s", (status, job_id,))
        self._commit()

    def _update_time_column(self, column:str, job_id:str):
        self._cursor.execute("UPDATE jobs SET %s=%s WHERE id=%s", (column, datetime.now(), job_id,))
        self._commit()

    def update_job_start_time_to_now(self, job_id:str):
        self._update_time_column('start_time', job_id)

    def update_job_finish_time_to_now(self, job_id:str):
        self._update_time_column('finish_time', job_id)

    def _was_results_table_initialized_by_the_job(self, job_id:str) -> bool:
        self._cursor.execute("SELECT 1 FROM results WHERE job_id=%s", (job_id,))
        query_results = self._cursor.fetchone()

        return query_results is not None and len(query_results) > 0

    def _initialize_results_table_for_job(self,job_id:str):
        self._cursor.execute("INSERT INTO results(job_id) VALUES(%s)", (job_id,))
        self._commit()

    def save_results(self, result_type:str, results:Resul, job_id:str):
        # It may be seen as a SQL Injection vulnerable code.
        # But theorectically, it's safe. Once we retrieved the data directly from database 
        # and the user had no access of what result_type name actually is
        column = result_type

        if(not self._was_results_table_initialized_by_the_job(job_id)):
            self._initialize_results_table_for_job(job_id)

        self._cursor.execute("UPDATE results SET %s=%s WHERE job_id=%s", (column, results, job_id,))
        self._commit()


    def _commit(self):
        """
            Make database changes definitive.
        """
        self._connection.commit()

    def close(self):
        """
            Finish database connection.
        """
        self._cursor.close()
        self._connection.close()
