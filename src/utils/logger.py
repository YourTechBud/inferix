import logging
import sys

logger = logging.getLogger("uvicorn")
logger.setLevel(logging.DEBUG)
logger.addHandler(logging.StreamHandler(sys.stdout))

def get_logger() -> logging.Logger:
    return logger