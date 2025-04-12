import time
import sys
from typing import List
from datetime import datetime
import os

class Time:
    """
    Helper class to handle time usage.  """

    @staticmethod
    def diff_in_days(modifcation_date: float) -> int:
        """
        get the difference between now and the date the file was modified
        and return the amount of days.
        """
        current_timestamp = time.time()
        now = datetime.fromtimestamp(current_timestamp)
        file_date = datetime.fromtimestamp(modifcation_date)

        return (abs(now - file_date)).days


def delete_file(file_path: str, days_passed: int, delete_after_days: int):
    """
    Delete this very file if the days criteria is fullfiled
    """

    if days_passed < delete_after_days:
        return

    print(f"Deleting file: {file_path}")
    try:
        os.remove(file_path)
    except FileNotFoundError:
        print("The file wasnt found, check it out")
    except OSError:
        print("It should be a file not a directory")


def get_file_modification_diff(file: str) -> int:
    """
    Get the difference of time since the last modification
    """
    modifcation_date = os.path.getmtime(file)
    amount_of_days_passed = Time.diff_in_days(modifcation_date)
    return amount_of_days_passed


def delete_files(base_path: str, files: List[str], time_to_delete: int):
    """
    Delete files from a list.
    """
    for file in files:
        full_file_path = os.path.join(base_path, file)
        amount_of_days_passed = get_file_modification_diff(full_file_path)
        delete_file(full_file_path, amount_of_days_passed, time_to_delete)


def clear_qasm(qasm_path: str, time_to_delete: int):
    """
    Clear qasm files.
    """
    files = os.listdir(qasm_path)
    delete_files(qasm_path, files, time_to_delete)


def clear_logs(logs_path: str, time_to_delete: int):
    """
    Clear log files in subdirectories
    """
    for directory in os.listdir(logs_path):
        dir_full_path = os.path.join(logs_path, directory)

        files = os.listdir(dir_full_path)
        delete_files(dir_full_path, files, time_to_delete)


if __name__ == "__main__":
    QASM_PATH = os.environ.get("QASM_PATH")
    LOGS_PATH = os.environ.get("LOGS_PATH")
    TIME_CRITERIA = os.environ.get("TIME_TO_DELETE")

    if None in (QASM_PATH, LOGS_PATH, TIME_CRITERIA):
        print("You must provide all env variables")
        sys.exit(1)

    try:
        TIME_CRITERIA = int(TIME_CRITERIA)  # type: ignore
    except ValueError:
        print(
            "It wasn't possible to convert the env variable TIME_TO_DELETE to int"
        )
        sys.exit(1)

    print("Checking qasm files")
    clear_qasm(QASM_PATH, TIME_CRITERIA)  # type: ignore

    print("Checking log files")
    clear_qasm(LOGS_PATH, TIME_CRITERIA)  # type: ignore
