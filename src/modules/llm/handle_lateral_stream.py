import asyncio
from typing import AsyncGenerator
from fastapi import Request
from models.http import StandardResponse
from sse_starlette import EventSourceResponse

from clients.redis import RedisClient
from utils.logger import get_logger

logger = get_logger()


async def handle_lateral_stream(ctx_id: str, ctx_name: str, interval: int = 200) -> EventSourceResponse:
    # This is the key we'll use to retrieve the result
    key_name = f"inferix:llm:{ctx_id}:{ctx_name}"

    # Get the redis client
    redis_client = await RedisClient.get_client()

    async def stream_response() -> AsyncGenerator[dict, None]:
        while True:
            # Get the response
            response = await redis_client.hgetall(key_name)  # type: ignore

            # Check is response is empty. This means that inference has probably not started yet.
            # For now we'll just wait till we get a response
            if bool(response):
                # Yield the response before anything else
                yield {"id": key_name, "data": response.get("response", "")}

                # Break if done
                if response.get("done", "False") == "True":
                    break

            # Sleep for interval
            await asyncio.sleep(interval / 1000)

    return EventSourceResponse(stream_response())


async def handle_delete_lateral_stream(ctx_id: str, ctx_name: str) -> StandardResponse:
    # This is the key we'll use to retrieve the result
    key_name = f"inferix:llm:{ctx_id}:{ctx_name}"

    # Get the redis client
    redis_client = await RedisClient.get_client()

    # Delete the key
    await redis_client.delete(key_name)

    return StandardResponse(message=f"Deleted lateral stream for {ctx_id}:{ctx_name}")
