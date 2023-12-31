from socket import AF_INET
from typing import List, Optional, Any

from aiohttp import ClientSession
from fastapi import FastAPI, HTTPException
from fastapi.logger import logger as fastAPI_logger

from models.http import StandardResponse  # convenient name


class HttpClient:
    aiohttp_client: Optional[ClientSession] = None

    @classmethod
    def get_aiohttp_client(cls) -> ClientSession:
        if cls.aiohttp_client is None:
            cls.aiohttp_client = ClientSession()

        return cls.aiohttp_client

    @classmethod
    async def close_aiohttp_client(cls) -> None:
        if cls.aiohttp_client:
            await cls.aiohttp_client.close()
            cls.aiohttp_client = None

    @classmethod
    async def post(cls, url: str, data: Any) -> str:
        client = cls.get_aiohttp_client()

        async with client.post(url, json=data) as response:
            if response.status != 200:
                raise HTTPException(
                    status_code=response.status,
                    detail=StandardResponse(message=(await response.json())["error"]).to_json(),
                )

            result = await response.text()

        return result


async def on_start_up() -> None:
    fastAPI_logger.info("Initialising async http client")
    HttpClient.get_aiohttp_client()


async def on_shutdown() -> None:
    fastAPI_logger.info("Destroying async http client")
    await HttpClient.close_aiohttp_client()
