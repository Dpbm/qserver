class CanceledJob(Exception):
    """
    An error for when a job is canceled by user.
    """

    def __init__(self):
        super().__init__("Job was canceled")
