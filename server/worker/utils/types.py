from typing import Dict, Any, List
from psycopg2.extras import RealDictRow

Backend = str
Results = Dict[int | str, float] | List[float]
Metadata = Dict[Any, Any]
ResultType = str
QasmFilePath = str

DBRow = RealDictRow | None
