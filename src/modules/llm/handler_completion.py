import json
import stat
from fastapi import HTTPException
from models.http import StandardResponse

from models.openai import CreateChatCompletionRequest, CreateChatCompletionResponse, ChatCompletionRequestMessage, CompletionUsage
from modules.llm.types import OllamaRequest, OllamaRequestOptions
from .helpers import add_message_for_fn_call, prepare_chat_completion_message, sanitize_json_text, to_unix_timestamp
from .ollama import call_ollama
from .prompts import chatml_tmpl


async def handle_completions(req: CreateChatCompletionRequest) -> CreateChatCompletionResponse:
    # We don't support streaming for now
    if req.stream:
        raise HTTPException(status_code=501, detail=StandardResponse(message="Streaming is not supported").to_json())
    
    # Add functions as a system prompt if needed
    add_message_for_fn_call(req.functions, req.messages)

    # Set the inference options
    # print(req.messages)
    # opts = InferenceOptions(messages=[ChatCompletionRequestMessage(content=req.messages[0].content, role=req.messages[0].role)], model=req.model)
    # if req.max_tokens is not None:
    #     opts.num_ctx = req.max_tokens
    # if req.temperature is not None:
    #     opts.temperature = req.temperature
    # if req.top_p is not None:
    #     opts.top_p = req.top_p

    # First prepare the prompt
    raw_prompt = chatml_tmpl.get_prompt(req.messages)

    # Make the request body for api call
    ollama_request_options = OllamaRequestOptions()
    if req.max_tokens is not None:
        ollama_request_options.num_ctx = req.max_tokens
    if req.temperature is not None:
        ollama_request_options.temperature = req.temperature
    if req.top_p is not None:
        ollama_request_options.top_p = req.top_p

    ollama_request = OllamaRequest(
        model=req.model,
        prompt=raw_prompt,
        stream=req.stream,
        options=ollama_request_options,
    )

    # Make the inference call
    # TODO: Limit the loop to 3 retries
    while True:
        output = await call_ollama(ollama_request)

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
        completion_tokens=output.eval_count or 0,
        prompt_tokens=output.prompt_eval_count or 0,
        total_tokens=(output.eval_count or 0) + (output.prompt_eval_count or 0),
    )

    # Prepare the response
    res = CreateChatCompletionResponse(
        id="1",
        created=to_unix_timestamp(output.created_at),
        usage=usage,
        model=output.model,
        object="chat.completion",
        choices=[prepare_chat_completion_message(output.response)],
    )

    # Return the response
    return res
