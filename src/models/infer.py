from typing import List, Optional
from xml.etree.ElementInclude import include
from pydantic import BaseModel

from .openai import ChatCompletionFunctions, ChatCompletionRequestMessage, CreateChatCompletionResponse


class RunInferenceContext(BaseModel):
    id: str
    key: str


class AddTextInstruction(BaseModel):
    text: str
    include_in_output: bool = True


class ConversationOptions(BaseModel):
    # Store key will add the inference response to the conversation history
    store_key: Optional[str] = None

    # Store entire history will store the entire conversation history overwriting the previous one
    store_entire_history: Optional[bool] = False

    # Load key will load the conversation from the conversation history
    load_key: Optional[str] = None

    # Assistant is used to determine how each message needs to be loaded.
    # Messages with the assitant will be passed with role "assistant" while other will be
    # passed with role "user"
    assistant_name: str = "default"


class RunInferenceInstructions(BaseModel):
    force_json: Optional[bool] = False  # TODO: Implement this
    conversation: Optional[ConversationOptions] = None
    add_prefix: Optional[AddTextInstruction] = None
    add_suffix: Optional[AddTextInstruction] = None
    enable_lateral_stream: Optional[bool] = False
    check_for_dedup: bool = False  # TODO: Implement this


class RunInferenceRequest(BaseModel):
    context: RunInferenceContext = RunInferenceContext(id="default", key="default")
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
