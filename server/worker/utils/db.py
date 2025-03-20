from datetime import datetime
import logging
import json
import psycopg2 as pg
import psycopg2.extras as pgextras
from .types import Results, ResultType, Backend, DBRow, Statuses, HelperMethods

logger = logging.getLogger(__name__)


class DB:
    """
    DB wrapper class for postgres (psycopg2)
    """

    # pylint: disable=too-many-arguments
    # pylint: disable=too-many-positional-arguments
    def __init__(self, host: str, port: str, db_name: str, user: str, password: str):
        self._connection = pg.connect(
            f"postgres://{user}:{password}@{host}:{port}/{db_name}"
        )
        # cursor_factory: https://www.geeksforgeeks.org/psycopg2-return-dictionary-like-values/
        self._cursor = self._connection.cursor(cursor_factory=pgextras.RealDictCursor)

    def get_job_data(self, job_id: str) -> DBRow:
        """
        Retrieve all job related data on database, including result types.
        """
        self._cursor.execute(
            """
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
        """,
            (job_id,),
        )  # get all data from job and create a json with the result types the user selected
        return self._cursor.fetchone()

    def update_job_status(self, status: Statuses, job_id: str):
        """
        Method used to update a job status by an arbitrary one.
        """

        try:
            self._cursor.execute(
                "UPDATE jobs SET status=%s WHERE id=%s",
                (
                    status,
                    job_id,
                ),
            )
            self._commit()
        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error("failed on update job status: %s for %s", status, job_id)
            logger.error("error: %s", str(error))
            self._rollback()

            raise error

    def update_job_start_time_to_now(self, job_id: str):
        """
        Update job starting time to this very moment.
        """

        try:
            self._cursor.execute(
                "UPDATE jobs SET start_time=%s WHERE id=%s",
                (
                    datetime.now(),
                    job_id,
                ),
            )
            self._commit()
        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error("failed on update start time to now for job %s", job_id)
            logger.error("error: %s", str(error))
            self._rollback()
            raise error

    def update_job_finish_time_to_now(self, job_id: str):
        """
        Update job finishing time to this very moment.
        """
        try:
            self._cursor.execute(
                "UPDATE jobs SET finish_time=%s WHERE id=%s",
                (
                    datetime.now(),
                    job_id,
                ),
            )
            self._commit()
        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error("failed on update finish time to now for job %s", job_id)
            logger.error("error: %s", str(error))
            self._rollback()
            raise error

    def _was_results_table_initialized_by_the_job(self, job_id: str) -> bool:
        """
        Check if the row containing job_id in results table was already populated
        """

        self._cursor.execute("SELECT 1 FROM results WHERE job_id=%s", (job_id,))
        query_results = self._cursor.fetchone()

        return query_results is not None and len(query_results) > 0

    def _initialize_results_table_for_job(self, job_id: str):
        """
        Add a row to results table to this specific job
        """
        try:
            self._cursor.execute("INSERT INTO results(job_id) VALUES(%s)", (job_id,))
            self._commit()
        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error("failed on intialize results table for job %s", job_id)
            logger.error("error: %s", str(error))
            self._rollback()
            raise error

    def _save_counts(self, results: Results, job_id: str):
        """
        Update results row inserting counts for this id
        """

        try:
            self._cursor.execute(
                "UPDATE results SET counts=%s WHERE job_id=%s",
                (
                    json.dumps(results),
                    job_id,
                ),
            )
            self._commit()

        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error(
                "failed on save counts for job: %s; counts=%s", job_id, results
            )
            logger.error("error: %s", str(error))
            self._rollback()
            raise error

    def _save_quasi_dist(self, results: Results, job_id: str):
        """
        Update results row inserting quasi_dist for this id
        """

        try:
            self._cursor.execute(
                "UPDATE results SET quasi_dist=%s WHERE job_id=%s",
                (
                    json.dumps(results),
                    job_id,
                ),
            )
            self._commit()

        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error(
                "failed on save quasi_dist for job: %s; dist=%s", job_id, results
            )
            logger.error("error: %s", str(error))
            self._rollback()
            raise error

    def _save_expval(self, results: Results, job_id: str):
        """
        Update results row inserting expval for this id
        """

        try:
            self._cursor.execute(
                "UPDATE results SET expval=%s WHERE job_id=%s",
                (
                    json.dumps(results),
                    job_id,
                ),
            )
            self._commit()

        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error(
                "failed on save expval for job: %s; values=%s", job_id, results
            )
            logger.error("error: %s", str(error))
            self._rollback()
            raise error

    def save_results(self, result_type: ResultType, results: Results, job_id: str):
        """
        Retrieve results and store them on database
        """
        if not self._was_results_table_initialized_by_the_job(job_id):
            self._initialize_results_table_for_job(job_id)

        helpers: HelperMethods = {
            "counts": self._save_counts,
            "quasi_dist": self._save_quasi_dist,
            "expval": self._save_expval,
        }

        assert helpers.get(result_type) is not None, "Invalid result type"

        save_func = helpers[result_type]
        save_func(results, job_id)

    def _commit(self):
        """
        Make database changes definitive.
        """
        self._connection.commit()

    def _rollback(self):
        """
        Undo database changes.
        """
        self._connection.rollback()

    def close(self):
        """
        Finish database connection.
        """
        self._cursor.close()
        self._connection.close()

    def get_plugin(self, backend: Backend) -> DBRow:
        """
        Retrieve the relative plugin name for the requested backend.
        """

        self._cursor.execute(
            "SELECT plugin FROM backends WHERE backend_name=%s", (backend,)
        )
        return self._cursor.fetchone()
