from email import charset
from typing import Optional
from redis.asyncio import Redis, from_url

from utils.logger import logger


class RedisClient:
    client: Optional[Redis] = None

    @classmethod
    def get_client(cls) -> Redis:
        """
        Returns the Redis client instance.

        If the client instance doesn't exist, it creates a new one with decode_responses set to True.

        Returns:
          Redis: The Redis client instance.
        """
        if cls.client is None:
            cls.client = Redis(decode_responses=True)

        return cls.client

    @classmethod
    async def close_client(cls) -> None:
        """
        Closes the Redis client connection.

        Args:
          None

        Returns:
          None
        """
        if cls.client:
            await cls.client.close()
            cls.client = None

    @classmethod
    async def append_to_sorted_set(cls, key: str, value: str) -> None:
        """
        Appends a value to a sorted set in Redis. It will first derive the max score by querying
        redis. It will then append the value with a score of max_score + 1

        Args:
          key (str): The key of the sorted set.
          value (str): The value to be appended.

        Returns:
          None
        """
        client = cls.get_client()

        # First lets get the max score
        member = await client.zrevrange(key, 0, 0, withscores=True)

        max_score: int = 0 if not bool(member) else int(member[0][1])

        # Now lets add the new member with a score of max_score + 1
        await client.zadd(key, {value: max_score + 1}, nx=True)
        await client.expire(key, 60 * 10) # Expire in 10 minutes

    @classmethod
    async def store_as_sorted_set(cls, key: str, values: list[str]) -> None:
        """
        Stores a list of values as a sorted set in Redis. It will overwrite the existing sorted set.

        Args:
          key (str): The key of the sorted set.
          values (list[str]): The values to be stored.

        Returns:
          None
        """
        client = cls.get_client()

        # First lets delete the key if it exists
        await client.delete(key)

        # Let's add the values as a sorted set to redis
        await client.zadd(key, {v: i + 1 for i, v in enumerate(values)})
        await client.expire(key, 60 * 10) # Expire in 10 minutes


async def create_redis_client() -> None:
    logger.info("Initialising async redis client")

    c = RedisClient.get_client()
    logger.info(f"Ping redis result: {await c.ping()}")


async def destroy_redis_client() -> None:
    logger.info("Destroying async redis client")
    await RedisClient.close_client()
