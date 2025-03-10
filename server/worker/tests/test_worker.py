from typing import Dict, Any
from worker import callback
from utils.types import Statuses


class TestWorker:
    """
    Test everything on callback function from worker.py.
    """

    def test_invalid_job_id(self):
        """
        Test if only valid ids are accepted
        """
        db = DB()
        channel = Channel()
        method = Method()
        body = Body("invalid-id")

        callback(channel, method, body, db)

        assert db.get_status() == Statuses.PENDING

    def test_invalid_result_types(self):
        """
        Test if job fails when using invalid result types
        """
        db = DB()
        channel = Channel()
        method = Method()
        body = Body("job-valid-id")

        callback(channel, method, body, db)

        assert db.get_status() == Statuses.FAILED

    def test_invalid_qasm(self):
        """
        Test if job fails when using invalid qasm file
        """
        db = DB(result_types={"counts": True, "quasi_dist": False, "expval": False})
        channel = Channel()
        method = Method()
        body = Body("job-valid-id")

        callback(channel, method, body, db)

        assert db.get_status() == Statuses.FAILED

    def test_invalid_backend(self):
        """
        Test if job fails when using invalid backend
        """
        db = DB(
            result_types={"counts": True, "quasi_dist": False, "expval": False},
            qasm="./tests/qasm_test.qasm",
        )
        channel = Channel()
        method = Method()
        body = Body("job-valid-id")

        callback(channel, method, body, db)

        assert db.get_status() == Statuses.FAILED

    def test_canceled_job(self):
        """
        Test if it stops when user cancels it
        """
        job_id = "job-valid-id"
        db = DB(
            result_types={"counts": True, "quasi_dist": False, "expval": False},
            qasm="./tests/qasm_test.qasm",
            backend="test",
        )
        db.update_job_status(Statuses.CANCELED, job_id)
        channel = Channel()
        method = Method()
        body = Body(job_id)

        callback(channel, method, body, db)

        assert db.get_status() == Statuses.CANCELED

    def test_invalid_status(self):
        """
        Test if only pending status is accepted to start running a job
        """
        job_id = "job-valid-id"
        db = DB(
            result_types={"counts": True, "quasi_dist": False, "expval": False},
            qasm="./tests/qasm_test.qasm",
            backend="test",
        )
        db.update_job_status(Statuses.FINISHED, job_id)
        channel = Channel()
        method = Method()
        body = Body(job_id)

        callback(channel, method, body, db)

        assert db.get_status() == Statuses.FINISHED


class Body:
    """
    A dummy class to act like the incoming body data
    """

    def __init__(self, job_id: str):
        self._id = job_id

    def decode(self) -> str:
        """
        Return decoded id data from body
        """
        return self._id


class Method:
    """
    A dummy class to act like a method parameter
    """

    delivery_tag = 1


class Channel:
    """
    A dummy class to act like a rabbitmq channel
    """

    def basic_ack(self, delivery_tag: Any):
        """
        Dummy ACK
        """


class DB:
    """
    A dummy class to act like a postgres db
    """

    # pylint: disable=dangerous-default-value
    def __init__(self, result_types: Dict = {}, qasm: str = "", backend: str = ""):
        self._data = {
            "id": "job-valid-id",
            "status": Statuses.PENDING,
            "finish_time": None,
            "qasm": qasm,
            "selected_result_types": result_types,
            "target_simulator": backend,
        }

    def _is_the_correct_id(self, job_id: str) -> bool:
        """
        check if the recieved id is correct
        """
        return self._data["id"] == job_id

    def get_job_data(self, job_id: str) -> Dict | None:
        """
        Get data if Id is correct
        """
        return self._data if self._is_the_correct_id(job_id) else None

    def update_job_status(self, status: Statuses, job_id: str):
        """
        Update job status for an arbitrary one
        """
        if not self._is_the_correct_id(job_id):
            return

        self._data["status"] = status

    def update_job_finish_time_to_now(self, job_id: str):
        """
        Update job finish time
        """
        if not self._is_the_correct_id(job_id):
            return

        self._data["finish_time"] = "now"

    def get_status(self) -> Statuses:
        """
        Get current job status
        """
        return Statuses(self._data["status"])
