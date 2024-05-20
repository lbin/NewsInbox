
from common.log import _reset_logger, logger
from pip._internal import main as pipmain


def install(package):
    pipmain(["install", package])


def install_requirements(file):
    pipmain(["install", "-r", file, "--upgrade"])
    _reset_logger(logger)


def check_dulwich():
    pass
