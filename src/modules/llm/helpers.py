import json
from datetime import datetime
from math import ceil
from typing import List, Optional

from models.openai import ChatCompletionFunctions, ChatCompletionMessageFunctionCall, ChatCompletionRequestMessage, ChatCompletionResponseMessage, CreateChatCompletionChoice

from .types import OllamaResponse


def add_message_for_fn_call(available_functions: Optional[List[ChatCompletionFunctions]], messages: List[ChatCompletionRequestMessage]):
    if available_functions is not None:
        system = """You may use the following FUNCTIONS in the response. Only use one function at a time. Give output in following OUTPUT_FORMAT in strict JSON if you want to call a function.
FUNCTIONS:"""
        for f in available_functions:
            system += f"1. Name: {f.name}\n"
            system += f"{f.description}\n"
            system += f"Parameters:\n"
            system += json.dumps(f.parameters) + "\n"
            system += "\n\n"

        system += """OUTPUT_FORMAT:
{
    "type": "FUNC_CALL",
    "name": "<name of function>",
    "parameters": "<parameters to pass to function>"
}
"""
        messages.append(ChatCompletionRequestMessage(content=system, role="system"))


def sanitize_json_text(text):
    # find the index of the first "{"
    start = text.find("{")
    # find the index of the last "}"
    end = text.rfind("}")
    # return the substring between the start and end indices
    return text[start : end + 1]


def to_unix_timestamp(ts: str) -> int:
    return ceil(int(datetime.fromisoformat(ts).timestamp()))


def prepare_chat_completion_message(
    response: str,
) -> CreateChatCompletionChoice:
    """
    Prepare a chat completion message.

    Args:
        response (str): The ollama generation output.

    Returns:
        CreateChatCompletionChoice: The prepared chat completion choice.
    """

    fn_call: Optional[ChatCompletionMessageFunctionCall] = None
    # Check if output suggests a function call
    if "FUNC_CALL" in response:
        f = dict(json.loads(response))
        fn_call = ChatCompletionMessageFunctionCall(name=f.get("name", ""), arguments=json.dumps(f.get("parameters")))

    return CreateChatCompletionChoice(
        index=0,
        finish_reason="stop" if fn_call is None else "function_call",
        message=ChatCompletionResponseMessage(content=response, role="assistant", function_call=fn_call),
    )
