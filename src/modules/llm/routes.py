from typing import List
from fastapi import APIRouter, Request
from models.list_models import ListModelsResponse
from modules.llm.handle_list_models import handle_list_models
from sse_starlette import EventSourceResponse

from models.http import StandardResponse
from models.infer import RunInferenceRequest, RunInferenceResponse
from models.openai import CreateChatCompletionRequest, CreateChatCompletionResponse

from .handle_infer import handle_delete_conversation_by_key, handle_delete_conversations_by_context, handle_infer
from .handle_lateral_stream import handle_delete_lateral_stream_state_by_context, handle_delete_lateral_stream_state_by_key, handle_lateral_stream
from .handler_completion import handle_completions

router = APIRouter(tags=["llm"], prefix="/api/llm/v1")


@router.post("/chat/completions", response_model_exclude_none=True)
async def chat_completions(req: CreateChatCompletionRequest) -> CreateChatCompletionResponse:
    return await handle_completions(req)


@router.post("/infer", response_model_exclude_none=True)
async def infer(req: RunInferenceRequest) -> RunInferenceResponse:
    return await handle_infer(req)


@router.delete("/infer/conversations/{ctx_id}", response_model_exclude_none=True)
async def delete_conversations_by_context(ctx_id: str) -> StandardResponse:
    return await handle_delete_conversations_by_context(ctx_id)


@router.delete("/infer/conversations/{ctx_id}/{key}", response_model_exclude_none=True)
async def delete_conversation(ctx_id: str, key: str) -> StandardResponse:
    return await handle_delete_conversation_by_key(ctx_id, key)


@router.get("/infer/streams/{ctx_id}/{ctx_name}", response_model_exclude_none=True)
async def infer_stream(ctx_id: str, ctx_name: str, interval: int = 200) -> EventSourceResponse:
    return await handle_lateral_stream(ctx_id, ctx_name, interval)


@router.delete("/infer/streams/{ctx_id}", response_model_exclude_none=True)
async def delete_infer_stream_state_by_context(ctx_id: str) -> StandardResponse:
    return await handle_delete_lateral_stream_state_by_context(ctx_id)


@router.delete("/infer/streams/{ctx_id}/{key}", response_model_exclude_none=True)
async def delete_infer_stream_state_by_key(ctx_id: str, key: str) -> StandardResponse:
    return await handle_delete_lateral_stream_state_by_key(ctx_id, key)


@router.delete("/infer/contexts/{ctx_id}", response_model_exclude_none=True)
async def delete_infer_context(ctx_id: str) -> StandardResponse:
    await handle_delete_conversations_by_context(ctx_id)
    return await handle_delete_lateral_stream_state_by_context(ctx_id)


@router.get("/models/list", response_model_exclude_none=True)
async def list_models() -> ListModelsResponse:
    return await handle_list_models()
