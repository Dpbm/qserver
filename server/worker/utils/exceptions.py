class CanceledJob(Exception):
    """
    An error for when a job is canceled by user.
    """

    def __init__(self):
        super().__init__("Job was canceled")


class IdNotFound(Exception):
    """
    An error for when a requested job has
    an invalid or not stored in database ID
    """

    def __init__(self):
        super().__init__("Job ID not found")


class InvalidResultTypes(Exception):
    """
    An error for when a requested job has
    invalid data for result types
    """

    def __init__(self):
        super().__init__("Invalid Result Types")


class InvalidQasmFile(Exception):
    """
    An error for when a requested job has
    invalid qasm file
    """

    def __init__(self):
        super().__init__(
            "Invalid QASM. You must provide the correct path for your code!"
        )


class InvalidBackend(Exception):
    """
    An error for when a requested job has
    invalid Backend name
    """

    def __init__(self):
        super().__init__("Invalid backend")


class InvalidStatus(Exception):
    """
    An error for when a requested job has
    a status different of pending
    """

    def __init__(self):
        super().__init__("Invalid Status")
