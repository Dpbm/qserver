from typing import Any


class CanceledJob(Exception):
    """
    An error for when a job is canceled by user.
    """

    def __init__(self, job_id: str):
        super().__init__(f"Job was canceled. JOB ID={job_id}")


class IdNotFound(Exception):
    """
    An error for when a requested job has
    an invalid or not stored in database ID
    """

    def __init__(self, data: Any):
        super().__init__(f"Job ID not found. Got Data: {data}")


class InvalidResultTypes(Exception):
    """
    An error for when a requested job has
    invalid data for result types
    """

    def __init__(self, types: Any):
        super().__init__(f"Invalid Result Types. Got: {types}")


class InvalidQasmFile(Exception):
    """
    An error for when a requested job has
    invalid qasm file
    """

    def __init__(self, qasm: Any):
        super().__init__(
            "Invalid QASM. You must provide the correct path for your code! "
            + f"Got: {qasm}"
        )


class InvalidBackend(Exception):
    """
    An error for when a requested job has
    invalid Backend name
    """

    def __init__(self, backend: Any):
        super().__init__(f"Invalid backend. Got: {backend}")


class InvalidStatus(Exception):
    """
    An error for when a requested job has
    a status different of pending
    """

    def __init__(self, status: Any):
        super().__init__(f"Invalid Status. Got: {status}")
