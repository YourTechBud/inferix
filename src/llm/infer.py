import json
from fastapi import HTTPException
from pydantic import ValidationError

from llm.prompts import get_chatml_tmpl
from llm.types import OllamaResponse
from models.http import StandardResponse
from models.openai import CreateChatCompletionRequest
from utils.http import HttpClient

chatml_tmpl = get_chatml_tmpl()


async def infer(options: CreateChatCompletionRequest) -> OllamaResponse:
    # First prepare the prompt
    raw_prompt = chatml_tmpl.get_prompt(options.messages)

    # Make the request body for api call
    req_body = {
        "model": options.model,
        "prompt": raw_prompt,
        "raw": True,
        "stream": False,
        "options": {
            "top_p": options.top_p if options.top_p is not None else 0.9,
            "num_ctx": options.max_tokens if options.max_tokens is not None else 4096,
            "temperature": options.temperature if options.temperature is not None else 0.2,
            "top_k": 25,
        },
    }

    # Make the api call
    try:
        res = await HttpClient.post("http://localhost:11434/api/generate", req_body)
        return OllamaResponse.model_validate_json(res)
    except ValidationError as e:
        raise HTTPException(
            status_code=500,
            detail=StandardResponse(message=f"Invalid response from Ollama: {e}", error=e.errors()).to_json(),
        )
