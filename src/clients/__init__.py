from .redis import create_redis_client, destroy_redis_client
from .http import create_http_client, destroy_http_client


async def create_clients() -> None:
    await create_http_client()
    await create_redis_client()


async def destroy_clients() -> None:
    await destroy_http_client()
    await destroy_redis_client()

