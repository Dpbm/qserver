import os
import logging

logger = logging.getLogger(__name__)


def create_path(file_path: str):
    """
    Create the whole path to a log file, including subdirectories and .log file.
    """

    folder, file = os.path.split(file_path)

    if not os.path.exists(folder):
        logger.debug("Creating folder: %s", folder)
        os.makedirs(folder)

    if not os.path.exists(file_path):
        logger.debug("Creating file: %s", file_path)
        with open(file, "w", encoding="utf-8") as log_file:
            log_file.write("")
