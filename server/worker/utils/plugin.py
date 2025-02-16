from .types import Results, Metadata, Backend, QasmFilePath, ResultType
from .build_pip_url import build_pip_url

class Plugin:
    def __init__(self, name:str):
        try:
            self._plugin = __import__(name)
        except ModuleNotFoundError as error:
            print(f"module {name} not found, attempting to install using pip....")

            import pip

            package_url = build_pip_url(name)

            # from: https://stackoverflow.com/questions/12332975/how-can-i-install-a-python-module-with-pip-programmatically-from-my-code
            # it's not the best way, but it's securer once we're not succeptible to code injection directly, once we're not
            # directly spawning terminal commands
            command = pip.main if(hasattr(pip, 'main')) else pip._internal.main
            command(['install', package_url])
            self._plugin = __import__(name)


    def run(self, target_backend: Backend, qasm_file:QasmFilePath, metdata: Metadata, result_type:ResultType) -> Results:
        return self._plugin.execute(target_backend, qasm_file, metadata, result_type)