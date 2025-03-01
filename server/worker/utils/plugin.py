from .types import Results, Metadata, Backend, QasmFilePath, ResultType
from .build_pip_url import build_pip_url, pipfy_name


class Plugin:
    """
    The Plugin class is meant to be a wrapper to the internally
    installed plugin.

    If the requested plugin is not installed, the class must
    check the quantum-plugins github org and verify if it's
    available.
    """

    def __init__(self, name: str):
        try:
            self._plugin = __import__(pipfy_name(name))
        except ModuleNotFoundError:
            print(f"module {name} not found, attempting to install using pip....")

            # pylint: disable=import-outside-toplevel
            import pip

            package_url = build_pip_url(name)

            # from:
            # https://stackoverflow.com/questions/12332975/how-can-i-install-a-python-module-with-pip-programmatically-from-my-code
            # it's not the best way, but it's securer once we're not succeptible to code
            # injection directly, once we're not directly spawning terminal commands
            command = pip.main if (hasattr(pip, "main")) else pip._internal.main  # type: ignore
            command(["install", package_url])
            self._plugin = __import__(pipfy_name(name))

    def run(
        self,
        target_backend: Backend,
        qasm_file: QasmFilePath,
        metadata: Metadata,
        result_type: ResultType,
    ) -> Results:
        """
        Is a wrapper to the installed plugin's execute function.

        It only pass the parameters to the instantiated object.
        """
        return self._plugin.execute(  # type: ignore
            target_backend, qasm_file, metadata, result_type
        )
