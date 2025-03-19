import importlib
import subprocess
import logging
from .types import Results, Metadata, Backend, QasmFilePath, ResultType
from .build_pip_url import build_pip_url, pipfy_name
from .sanitize import sanitize_pip_name

logger = logging.getLogger(__name__)


class Plugin:
    """
    The Plugin class is meant to be a wrapper to the internally
    installed plugin.

    If the requested plugin is not installed, the class must
    check the quantum-plugins github org and verify if it's
    available.
    """

    def __init__(self, name: str):
        has_module = True
        imported_plugin = None

        sanitized_name = sanitize_pip_name(name)
        import_name = pipfy_name(sanitized_name)

        try:

            logger.debug("trying to import %s....", import_name)
            imported_plugin = importlib.import_module(import_name)

        except ModuleNotFoundError as error:
            logger.error("module %s not found", import_name)
            logger.error("%s", str(error))
            has_module = False

        except ImportError as error:
            logger.error("failed on first import %s", import_name)
            logger.error("%s", str(error))
            has_module = False

        # pylint: disable=broad-exception-caught
        except Exception as error:
            logger.error("Another error occoured during first import")
            logger.error("%s", str(error))

        if not has_module:
            try:
                logger.debug("attempting to install %s....", import_name)

                package_url = build_pip_url(sanitized_name)
                logger.debug("installing plugin from: %s", package_url)

                subprocess.check_call(["python", "-m", "pip", "install", package_url])
                logger.debug("plugin was installed successfuly!")

                imported_plugin = importlib.import_module(import_name)

            except subprocess.CalledProcessError as error:
                logger.error("Failed on install plugin")
                logger.error(str(error))
                logger.error("stderr: %s", error.stderr)
                logger.error("output: %s", error.output)

            except ModuleNotFoundError as error:
                logger.error("Failed on second import (not found)")
                logger.error(str(error))

            except ImportError as error:
                logger.error("failed on second import (import error)")
                logger.error(str(error))

            # pylint: disable=broad-exception-caught
            except Exception as error:
                logger.error("Another error occoured during installation/second import")
                logger.error("%s", str(error))

        if imported_plugin is None:
            raise ValueError("Invalid plugin")

        logger.debug("got plugin!")
        self._plugin = imported_plugin.Plugin()

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
