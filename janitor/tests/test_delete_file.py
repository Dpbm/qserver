import os
import pytest
from janitor import delete_file, delete_files, clear_qasm, clear_logs


class TestDeleteFile:
    """
    Test suite for deleting files
    """

    def test_delete_file(self, file_to_delete_path, files_path, amount_of_files):
        """Should delete the file with no problems"""
        delete_file(file_to_delete_path, 1, 1)
        assert len(os.listdir(files_path)) == amount_of_files - 1

    def test_file_not_found(self, files_path, amount_of_files):
        """should raise an error internally and do nothing with our files"""
        delete_file(os.path.join("not", "exists"), 1, 1)
        assert len(os.listdir(files_path)) == amount_of_files

    def test_not_in_time_to_delete(
        self, file_to_delete_path, files_path, amount_of_files
    ):
        """should delete nothing"""
        delete_file(file_to_delete_path, 1, 3)
        assert len(os.listdir(files_path)) == amount_of_files

    # ------ DELETE FILES ------

    def test_delete_nothing(self, files_path, amount_of_files):
        """should delete no files"""
        delete_files(files_path, [], 1)
        assert len(os.listdir(files_path)) == amount_of_files

    def test_delete_single_file(self, files_path, files, amount_of_files):
        """should delete a single file"""
        delete_files(files_path, [files[0]], 0)
        assert len(os.listdir(files_path)) == amount_of_files - 1

    def test_delete_all_files(self, files_path, files, amount_of_files):
        """should delete all files"""
        delete_files(files_path, files, 0)
        assert len(os.listdir(files_path)) == 0

    # ------ CLEAR QASM -----

    def test_clear_qasm_files(self, files_path, amount_of_files):
        """should delete all files"""
        clear_qasm(files_path, 0)
        assert len(os.listdir(files_path)) == 0

    # ------ CLEAR LOGS -----

    def test_no_directory_error(self, files_path):
        """should raise an error, since we don't have the required directories for logs"""
        with pytest.raises(Exception):
            clear_logs(files_path, 0)

    def test_delete_logs(
        self, files_path, amount_of_files, logs_path, logs_subdir, amount_of_log_files
    ):
        """should delete all log files"""
        clear_logs(logs_path, 0)
        assert len(os.listdir(files_path)) == amount_of_files
        assert len(os.listdir(logs_subdir)) == 0
