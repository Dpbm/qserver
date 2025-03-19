import re


def sanitize_pip_name(name: str) -> str:
    """
    remove any ilegal character from pip's package name
    """

    # check for invalid characters like: & * , ; . ) ( ] [ { }
    pattern = re.compile(r"([^A-Za-z0-9-_])")
    return re.sub(pattern, "", name)
