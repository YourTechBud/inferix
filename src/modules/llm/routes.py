from typing import List
from fastapi import APIRouter, Request
from sse_starlette import EventSourceResponse

from models.http import StandardResponse
from models.infer import RunInferenceRequest, RunInferenceResponse
from models.openai import CreateChatCompletionRequest, CreateChatCompletionResponse

from .handle_infer import handle_delete_conversation, handle_infer
from .handle_lateral_stream import handle_delete_lateral_stream_state, handle_lateral_stream
from .handler_completion import handle_completions

router = APIRouter(tags=["llm"], prefix="/api/llm/v1")


@router.post("/chat/completions", response_model_exclude_none=True)
async def chat_completions(req: CreateChatCompletionRequest) -> CreateChatCompletionResponse:
    return await handle_completions(req)


@router.post("/infer", response_model_exclude_none=True)
async def infer(req: RunInferenceRequest) -> RunInferenceResponse:
    return await handle_infer(req)

@router.delete("/infer/conversations/{ctx_id}/{key}", response_model_exclude_none=True)
async def delete_conversation(ctx_id: str, key: str) -> StandardResponse:
    return await handle_delete_conversation(ctx_id, key)


@router.get("/infer/streams/{ctx_id}/{ctx_name}", response_model_exclude_none=True)
async def infer_stream(ctx_id: str, ctx_name: str, interval: int = 200) -> EventSourceResponse:
    return await handle_lateral_stream(ctx_id, ctx_name, interval)

@router.delete("/infer/streams/{ctx_id}/{ctx_name}", response_model_exclude_none=True)
async def delete_infer_stream_state(ctx_id: str, ctx_name: str) -> StandardResponse:
    return await handle_delete_lateral_stream_state(ctx_id, ctx_name)