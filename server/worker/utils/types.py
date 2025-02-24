from typing import Dict, Any
from psycopg2.extras import RealDictRow

Backend = str
Results = Dict[int | str, float] | float
Metadata = Dict[Any, Any]
ResultType = str
QasmFilePath = str

DBRow = RealDictRow | None
