from typing import List, Optional
import datetime
from pydantic import BaseModel

from models.openai import ChatCompletionRequestMessage


class InferenceOptions(BaseModel):
    messages: List[ChatCompletionRequestMessage]
    num_ctx: int = 4096
    temperature: float = 0.2
    top_p: float = 0.9
    top_k: int = 40
    model: str


class OllamaResponse(BaseModel):
    model: str
    created_at: str
    response: str
    done: bool
    context: Optional[List[int]] = None
    total_duration: int
    load_duration: int
    prompt_eval_count: int
    prompt_eval_duration: int
    eval_count: int
    eval_duration: int
