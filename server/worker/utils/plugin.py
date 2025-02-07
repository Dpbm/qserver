from .types import Results

class Plugin:
    def __init__(self, name:str):
        self._name = name

    def run(self, qasm_file:str, result_type:str) -> Results:
        pass