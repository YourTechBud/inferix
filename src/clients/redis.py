from email import charset
from typing import Optional
from redis.asyncio import Redis, from_url

from utils.logger import logger


class RedisClient:
    client: Optional[Redis] = None

    @classmethod
    def get_client(cls) -> Redis:
        if cls.client is None:
            cls.client = Redis(decode_responses=True)

        return cls.client

    @classmethod
    async def close_client(cls) -> None:
        if cls.client:
            await cls.client.close()
            cls.client = None


async def create_redis_client() -> None:
    logger.info("Initialising async redis client")

    c = RedisClient.get_client()
    logger.info(f"Ping redis result: {await c.ping()}")


async def destroy_redis_client() -> None:
    logger.info("Destroying async redis client")
    await RedisClient.close_client()
