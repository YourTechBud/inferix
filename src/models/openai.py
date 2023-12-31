# Define a class for CreateChatCompletionRequest schema
from typing import Any, Dict, Optional, List, Union
from pydantic import BaseModel, validator


class ChatCompletionMessageFunctionCall(BaseModel):
    arguments: str  # A string
    name: str  # A string


# Define a class for ChatCompletionRequestMessage schema
class ChatCompletionRequestMessage(BaseModel):
    content: Optional[str]  # A string or None
    function_call: Optional[ChatCompletionMessageFunctionCall] = None  # A dict or None
    name: Optional[str] = None  # A string or None
    role: str  # A string and one of ['system', 'user', 'assistant', 'function']

    # Validate the role field
    @validator("role")
    def check_role(cls, v):
        # Check if the value is in the enum
        if v not in ["system", "user", "assistant", "function"]:
            raise ValueError('role must be one of ["system", "user", "assistant", "function"]')
        return v


# Define a class for ChatCompletionFunctionCallOption schema
class ChatCompletionFunctionCallOption(BaseModel):
    name: str  # A string


# Define a class for ChatCompletionFunctions schema
class ChatCompletionFunctions(BaseModel):
    description: str  # A string
    name: str  # A string
    parameters: Dict[str, Any]  # A ChatCompletionFunctionParameters instance


class CreateChatCompletionRequest(BaseModel):
    messages: List[ChatCompletionRequestMessage]  # A list of ChatCompletionRequestMessage instances
    model: str  # A string
    frequency_penalty: Optional[float] = 0  # A float or None, default is 0
    function_call: Optional[Union[str, ChatCompletionFunctionCallOption]] = None  # A string, a ChatCompletionFunctionCallOption instance, or None, default is None
    functions: Optional[List[ChatCompletionFunctions]] = None  # A list of ChatCompletionFunctions instances or None, default is None
    logit_bias: Optional[Dict[str, int]] = None  # A dict with string keys and integer values or None, default is None
    max_tokens: Optional[int] = None  # An integer or None, default is None
    n: Optional[int] = 1  # An integer or None, default is 1
    presence_penalty: Optional[float] = 0  # A float or None, default is 0
    stop: Optional[Union[str, List[str]]] = None  # A string, a list of strings, or None, default is None
    stream: Optional[bool] = False  # A boolean or None, default is False
    temperature: Optional[float] = None  # A float or None, default is 1
    top_p: Optional[float] = None  # A float or None, default is 1
    user: Optional[str] = None  # A string or None, default is None

    # Validate the frequency_penalty field
    @validator("frequency_penalty")
    def check_frequency_penalty(cls, v):
        # Check if the value is in the range
        if v is not None and (v < -2 or v > 2):
            raise ValueError("frequency_penalty must be between -2 and 2")
        return v

    # Validate the function_call field
    @validator("function_call")
    def check_function_call(cls, v):
        # Check if the value is in the enum
        if isinstance(v, str) and v not in ["none", "auto"]:
            raise ValueError('function_call must be one of ["none", "auto"]')
        return v

    # Validate the n field
    @validator("n")
    def check_n(cls, v):
        # Check if the value is in the range
        if v is not None and (v < 1 or v > 128):
            raise ValueError("n must be between 1 and 128")
        return v

    # Validate the presence_penalty field
    @validator("presence_penalty")
    def check_presence_penalty(cls, v):
        # Check if the value is in the range
        if v is not None and (v < -2 or v > 2):
            raise ValueError("presence_penalty must be between -2 and 2")
        return v

    # Validate the temperature field
    @validator("temperature")
    def check_temperature(cls, v):
        # Check if the value is in the range
        if v is not None and (v < 0 or v > 2):
            raise ValueError("temperature must be between 0 and 2")
        return v

    # Validate the top_p field
    @validator("top_p")
    def check_top_p(cls, v):
        # Check if the value is in the range
        if v is not None and (v < 0 or v > 1):
            raise ValueError("top_p must be between 0 and 1")
        return v


# Define a class for ChatCompletionResponseMessage schema
class ChatCompletionResponseMessage(BaseModel):
    content: Optional[str]  # A string or None
    function_call: Optional[ChatCompletionMessageFunctionCall] = None  # A dict or None
    role: str  # A string and one of ['system', 'user', 'assistant', 'function']

    # Validate the role field
    @validator("role")
    def check_role(cls, v):
        # Check if the value is in the enum
        if v not in ["system", "user", "assistant", "function"]:
            raise ValueError('role must be one of ["system", "user", "assistant", "function"]')
        return v


# Define a class for CreateChatCompletionChoice schema
class CreateChatCompletionChoice(BaseModel):
    finish_reason: str  # A string and one of ["stop", "length", "function_call", "content_filter"]
    index: int  # An integer
    message: ChatCompletionResponseMessage  # A ChatCompletionResponseMessage instance

    # Validate the finish_reason field
    @validator("finish_reason")
    def check_finish_reason(cls, v):
        # Check if the value is in the enum
        if v not in ["stop", "length", "function_call", "content_filter"]:
            raise ValueError('finish_reason must be one of ["stop", "length", "function_call", "content_filter"]')
        return v


# Define a class for CompletionUsage schema
class CompletionUsage(BaseModel):
    completion_tokens: int  # An integer
    prompt_tokens: int  # An integer
    total_tokens: int  # An integer


# Define a class for CreateChatCompletionResponse schema
class CreateChatCompletionResponse(BaseModel):
    id: str  # A string
    choices: List[CreateChatCompletionChoice]  # A list of CreateChatCompletionChoice instances
    created: int  # An integer
    model: str  # A string
    object: str  # A string
    usage: CompletionUsage  # A CompletionUsage instance
