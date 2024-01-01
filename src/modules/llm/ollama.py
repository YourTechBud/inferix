from typing import AsyncGenerator
from fastapi import HTTPException
from pydantic import ValidationError

from models.http import StandardResponse
from clients.http import HttpClient

from .prompts import get_chatml_tmpl
from .types import OllamaRequest, OllamaResponse


async def call_ollama(req: OllamaRequest):
    try:
        res = await HttpClient.post("http://localhost:11434/api/generate", req.to_json())
        return OllamaResponse.model_validate_json(res)
    except ValidationError as e:
        raise HTTPException(
            status_code=500,
            detail=StandardResponse(message=f"Invalid response from Ollama: {e}", error=e.errors()).to_json(),
        )


async def call_ollama_stream(req: OllamaRequest) -> AsyncGenerator[OllamaResponse, None]:
    response: str = ""
    try:
        async for line in HttpClient.post_stream("http://localhost:11434/api/generate", req.to_json()):
            res = OllamaResponse.model_validate_json(line)

            # Keep concating the response
            response += res.response
            res.response = response

            # Now we can yield the response
            yield res
    except ValidationError as e:
        raise HTTPException(
            status_code=500,
            detail=StandardResponse(message=f"Invalid response from Ollama: {e}", error=e.errors()).to_json(),
        )
