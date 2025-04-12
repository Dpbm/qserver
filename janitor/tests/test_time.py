import time
from datetime import timedelta, datetime
from janitor import Time


class TestTimeHelperFunction:
    """
    Test Time helper
    """

    def test_zero_days(self):
        """Should return 0"""
        assert Time.diff_in_days(time.time()) == 0

    def test_three_days(self):
        """should return 3"""
        assert Time.diff_in_days(
            datetime.timestamp(timedelta(days=3) + datetime.today())
        )
