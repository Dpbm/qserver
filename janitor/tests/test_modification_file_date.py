from janitor import get_file_modification_diff


class TestModificationFileDate:
    """
    Test suite to ensure the correct time difference is given
    """

    def test_time_diff(self, file_to_delete_path):
        """
        Should return zero, since the file was just created
        """
        assert get_file_modification_diff(file_to_delete_path) == 0
