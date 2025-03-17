from typing import Any
import os


def valid_data_for_id(data: Any) -> bool:
    """
    Verify if the returned data from db is valid
    """
    return isinstance(data, dict) and len(data.items()) > 0


def valid_result_types(result_types: Any) -> bool:
    """
    Check if the returned result_types from db are valid
    """
    return (
        isinstance(result_types, dict)
        and len(result_types.items()) > 0
        and len(list(filter(lambda x: x, result_types.values()))) > 0
    )


def valid_qasm(qasm: Any) -> bool:
    """
    check if the qasm file path is valid and the file exists
    """
    return isinstance(qasm, str) and len(qasm) > 10 and os.path.exists(qasm)


def valid_backend(backend: Any) -> bool:
    """
    check if the requested backend is valid
    """
    return isinstance(backend, str) and len(backend) > 0
