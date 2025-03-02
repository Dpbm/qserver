class CanceledJob(Exception):
    def __init__(self):
        super().__init__("Job was canceled")