from typing import List, Optional
from xml.etree.ElementInclude import include
from pydantic import BaseModel

from .openai import ChatCompletionFunctions, ChatCompletionRequestMessage, CreateChatCompletionResponse


class RunInferenceContext(BaseModel):
    id: str
    name: str


class AddTextInstruction(BaseModel):
    text: str
    include_in_output: bool = True


class RunInferenceInstructions(BaseModel):
    force_json: Optional[bool] = False
    conversation_key: Optional[str] = None
    add_prefix: Optional[AddTextInstruction] = None
    add_suffix: Optional[AddTextInstruction] = None
    enable_lateral_stream: Optional[bool] = False


class RunInferenceRequest(BaseModel):
    context: RunInferenceContext = RunInferenceContext(id="default", name="default")
    messages: List[ChatCompletionRequestMessage]
    model: str
    num_ctx: int = 4096
    temperature: float = 0.2
    top_p: float = 0.9
    top_k: int = 40
    functions: Optional[List[ChatCompletionFunctions]] = None
    stream: Optional[bool] = False
    instructions: Optional[RunInferenceInstructions] = RunInferenceInstructions()


class RunInferenceResponse(CreateChatCompletionResponse):
    pass
