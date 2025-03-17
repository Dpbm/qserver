import os
import logging

logger = logging.getLogger(__name__)

def create_path(file_path:str):
    folder, file = os.path.split(file_path)

    if not os.path.exists(folder):
        logger.debug(f"Creating folder: {folder}")
        os.makedirs(folder)
    
    if not os.path.exists(file_path):
        logger.debug(f"Creating file: {file_path}")
        with open(file, "w", encoding="utf-8") as log_file:
            log_file.write("")

        
    