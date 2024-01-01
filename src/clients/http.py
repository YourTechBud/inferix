from typing import AsyncGenerator, Optional, Any

from aiohttp import ClientSession
from fastapi import HTTPException

from models.http import StandardResponse
from utils.logger import logger


class HttpClient:
    aiohttp_client: Optional[ClientSession] = None

    @classmethod
    def get_client(cls) -> ClientSession:
        if cls.aiohttp_client is None:
            cls.aiohttp_client = ClientSession()

        return cls.aiohttp_client

    @classmethod
    async def close_client(cls) -> None:
        if cls.aiohttp_client:
            await cls.aiohttp_client.close()
            cls.aiohttp_client = None

    @classmethod
    async def post(cls, url: str, data: str) -> str:
        client = cls.get_client()

        async with client.post(url, data=data) as response:
            if response.status != 200:
                raise HTTPException(
                    status_code=response.status,
                    detail=StandardResponse(message=(await response.json())["error"]).to_json(),
                )

            result = await response.text()

        return result
    
    @classmethod
    async def post_stream(cls, url: str, data: str) -> AsyncGenerator[str, None]:
        client = cls.get_client()

        async with client.post(url, data=data) as response:
            if response.status != 200:
                raise HTTPException(
                    status_code=response.status,
                    detail=StandardResponse(message=(await response.json())["error"]).to_json(),
                )
            
            async for line in response.content:
                yield line.decode("utf-8")



async def create_http_client() -> None:
    logger.info("Initialising async http client")
    HttpClient.get_client()


async def destroy_http_client() -> None:
    logger.info("Destroying async http client")
    await HttpClient.close_client()