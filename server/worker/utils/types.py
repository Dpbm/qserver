from typing import Dict, Any, List, Optional, Literal, Callable
from enum import Enum
from psycopg2.extras import RealDictRow

Backend = str
Results = Dict[int | str, float] | List[float]
Metadata = Dict[Any, Any]
ResultType = Literal["counts", "quasi_dist", "expval"]
QasmFilePath = str

DBRow = RealDictRow | None
HelperMethods = Dict[ResultType, Callable]


def port_to_int(port: Optional[str]) -> Optional[int]:
    """
    Convert incoming port env to an int.
    """

    if port is None:
        return None

    try:

        int_port = int(port)

        if int_port < 0:
            raise ValueError("Invalid Port")

        return int_port
    # pylint: disable=broad-exception-caught
    except Exception:
        return None


class Statuses(Enum):
    """
    Possible status that a job can have.
    """

    PENDING = "pending"
    RUNNING = "running"
    FINISHED = "finished"
    CANCELED = "canceled"
    FAILED = "failed"
