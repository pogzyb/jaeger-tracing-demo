from functools import lru_cache

from pydantic import BaseSettings


class Config(BaseSettings):
    app_name: str = 'py-service'
    app_version: str = '0.1.0'
    debug: bool = False
    description: str = 'simple API using python'


@lru_cache
def get_config() -> Config:
    return Config()
