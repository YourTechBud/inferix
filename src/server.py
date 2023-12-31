import json
from datetime import datetime
from math import ceil
from fastapi import FastAPI
from llm.infer import infer
from llm.types import InferenceOptions, OllamaResponse

from utils.http import on_shutdown, on_start_up
from models.openai import *

app = FastAPI(title="Inferix", description="OpenAI compatible API with extra goodies.", on_startup=[on_start_up], on_shutdown=[on_shutdown])


@app.post(
    "/api/v1/chat/completions",
    response_model_exclude_none=True,
)
async def chat_completion(req: CreateChatCompletionRequest) -> CreateChatCompletionResponse:
    # Add functions as a system prompt
    if req.functions is not None:
        system = """You may use the following FUNCTIONS in the response. Only use one function at a time. Give output in following OUTPUT_FORMAT in strict JSON if you want to call a function.
FUNCTIONS:"""
        for f in req.functions:
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
        req.messages.append(ChatCompletionRequestMessage(content=system, role="system"))

    # Set the inference options
    # print(req.messages)
    # opts = InferenceOptions(messages=[ChatCompletionRequestMessage(content=req.messages[0].content, role=req.messages[0].role)], model=req.model)
    # if req.max_tokens is not None:
    #     opts.num_ctx = req.max_tokens
    # if req.temperature is not None:
    #     opts.temperature = req.temperature
    # if req.top_p is not None:
    #     opts.top_p = req.top_p

    # Make the inference call
    # TODO: Limit the loop to 3 retries
    while True:
        output = await infer(req)

        # Lets start by cleaning up the output
        output.response = output.response.strip()

        # Rerun inference if output was smaller than expected
        if len(output.response) <= 2:
            continue

        # We will need to clean the output a bit if it was a function call request
        try:
            if "FUNC_CALL" in output.response:
                # First try to sanitize the text
                output.response = sanitize_json_text(output.response)

                # Attempt to serialize it
                f = dict(json.loads(output.response))
        except json.JSONDecodeError:
            # We want to retry inference if the response was not JSON serializable
            continue

        # All looks good. We can exit the loop
        break

    # Prepare usage object
    usage: CompletionUsage = CompletionUsage(
        completion_tokens=output.eval_count,
        prompt_tokens=output.prompt_eval_count,
        total_tokens=output.eval_count + output.prompt_eval_count,
    )

    # Prepare the response
    res = CreateChatCompletionResponse(
        id="1",
        created=to_unix_timestamp(output.created_at),
        usage=usage,
        model=output.model,
        object="chat.completion",
        choices=[prepare_chat_completion_message(output)],
    )

    # Return the response
    return res


def sanitize_json_text(text):
    # find the index of the first "{"
    start = text.find("{")
    # find the index of the last "}"
    end = text.rfind("}")
    # return the substring between the start and end indices
    return text[start : end + 1]


def prepare_chat_completion_message(
    output: OllamaResponse,
) -> CreateChatCompletionChoice:
    """
    Prepare a chat completion message.

    Args:
        output (OllamaResponse): The ollama generation.

    Returns:
        CreateChatCompletionChoice: The prepared chat completion choice.
    """

    fn_call: Optional[ChatCompletionMessageFunctionCall] = None
    # Check if output suggests a function call
    if "FUNC_CALL" in output.response:
        f = dict(json.loads(output.response))
        fn_call = ChatCompletionMessageFunctionCall(name=f.get("name", ""), arguments=json.dumps(f.get("parameters")))

    return CreateChatCompletionChoice(
        index=0,
        finish_reason="stop" if fn_call is None else "function_call",
        message=ChatCompletionResponseMessage(content=output.response, role="assistant", function_call=fn_call),
    )


def to_unix_timestamp(ts: str) -> int:
    return ceil(int(datetime.fromisoformat(ts).timestamp()))
