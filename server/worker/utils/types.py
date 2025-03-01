from typing import Dict, Any, List, Optional
from psycopg2.extras import RealDictRow

Backend = str
Results = Dict[int | str, float] | List[float]
Metadata = Dict[Any, Any]
ResultType = str
QasmFilePath = str

DBRow = RealDictRow | None


def port_to_int(port:str) -> Optional[int]:
    try:
        int_port = int(port)

        if(int_port < 0):
            raise ValueError("Invalid Port")

        return int_port
    except:
        return None