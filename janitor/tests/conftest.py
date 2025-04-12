import os
from typing import List
import pytest

# pylint: disable=redefined-outer-name


@pytest.fixture
def files_path() -> str:
    """Path to files which we gonna delete during our tests"""
    return os.path.join(".", "tests", "files", "qasm")


@pytest.fixture
def file_to_delete_path(files_path) -> str:
    """The path to the main file to be deleted"""
    return os.path.join(files_path, "to-delete.txt")


@pytest.fixture
def files() -> List[str]:
    """list of files to be created"""
    return ["to-delete.txt", "1.txt", "2.txt"]


@pytest.fixture
def amount_of_files(files) -> int:
    """ "return the amount of files we have"""
    return len(files)


@pytest.fixture
def logs_path() -> str:
    """the logs subdirectory"""
    return os.path.join(".", "tests", "files", "logs")


@pytest.fixture
def logs_subdir(logs_path) -> str:
    """the logs inner directory"""
    return os.path.join(logs_path, "server")


@pytest.fixture
def log_files() -> List[str]:
    """the log files we have"""
    return ["a.txt", "b.txt"]


@pytest.fixture
def amount_of_log_files(log_files) -> int:
    """amount of log files we have"""
    return len(log_files)


@pytest.fixture(autouse=True)
def create_file_and_path(files_path, files, logs_subdir, log_files):
    """Generate the files we need"""

    os.makedirs(files_path, exist_ok=True)
    os.makedirs(logs_subdir, exist_ok=True)

    for file in files:
        file_path = os.path.join(files_path, file)

        if not os.path.exists(file_path):
            with open(file_path, "w", encoding="utf-8") as f:
                f.write("a")

    for file in log_files:
        file_path = os.path.join(logs_subdir, file)

        if not os.path.exists(file_path):
            with open(file_path, "w", encoding="utf-8") as f:
                f.write("a")
