from .types import Results, Metadata, Backend, QasmFilePath, ResultType

class Plugin:
    def __init__(self, name:str):
        self._plugin = __import__(name)

    def run(self, target_backend: Backend, qasm_file:QasmFilePath, metdata: Metadata, result_type:ResultType) -> Results:
        return self._plugin.execute(target_backend, qasm_file, metadata, result_type)