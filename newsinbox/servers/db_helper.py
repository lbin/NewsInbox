from pymongo import MongoClient
from pymongo.errors import DuplicateKeyError

from ..common.logger import logger


def get_mongo_collection(config:dict, name: str):
    """
    name in format of `<database>.<collection>`,
    only first dot is considered as separator.
    """

    url = f"mongodb://{config['mongo_user']}:{config['mongo_password']}@{config['mongo_host']}"
    db, table = name.split(".", maxsplit=1)
    return MongoClient(url).get_database(db).get_collection(table)

def add_news(config: dict, name: str, news: dict):
    try:
        get_mongo_collection(config, name).insert_one(news)
    except DuplicateKeyError:
        logger.error("Key already exists, update or handle the error accordingly.")