import json
from math import log
from clients.redis import RedisClient
from fastapi import HTTPException

from models.http import StandardResponse
from models.infer import RunInferenceInstructions, RunInferenceRequest, RunInferenceResponse
from models.openai import CompletionUsage
from utils.logger import logger

from .types import OllamaRequest, OllamaRequestOptions
from .helpers import add_message_for_fn_call, prepare_chat_completion_message, sanitize_json_text, to_unix_timestamp
from .ollama import call_ollama_stream
from .prompts import chatml_tmpl


async def handle_infer(req: RunInferenceRequest) -> RunInferenceResponse:
    # We don't support streaming for now
    if req.stream:
        raise HTTPException(status_code=501, detail=StandardResponse(message="Streaming is not supported").to_json())

    if req.instructions == None:
        req.instructions = RunInferenceInstructions()

    # Add functions as a system prompt if needed
    add_message_for_fn_call(req.functions, req.messages)

    # Check if we need to stream the response laterally
    if req.instructions.enable_lateral_stream:
        req.stream = True

    # First prepare the prompt
    raw_prompt = chatml_tmpl.get_prompt(req.messages)

    # Make the request body for api call
    ollama_request_options = OllamaRequestOptions()
    if req.num_ctx is not None:
        ollama_request_options.num_ctx = req.num_ctx
    if req.temperature is not None:
        ollama_request_options.temperature = req.temperature
    if req.top_p is not None:
        ollama_request_options.top_p = req.top_p
    if req.top_k is not None:
        ollama_request_options.top_k = req.top_k

    ollama_request = OllamaRequest(
        model=req.model,
        prompt=raw_prompt,
        stream=True,
        options=ollama_request_options,
    )

    # Populate useful variables
    prefix_text = "" if req.instructions.add_prefix is None else req.instructions.add_prefix.text
    suffix_text = "" if req.instructions.add_suffix is None else req.instructions.add_suffix.text

    # Make the inference call
    # TODO: Limit the loop to 3 retries
    while True:
        # The response text
        response_text = ""

        # Some important stats we'll need
        eval_count = 0
        prompt_eval_count = 0
        created_at = ""

        async for chunk in call_ollama_stream(ollama_request):
            # Let's make sure our response_text is up to date
            # Not we get the entire response again. We don't need to concat the response
            response_text = chunk.response

            # Add the prefix if we need to include it in the response
            if req.instructions.add_prefix is not None and req.instructions.add_prefix.include_in_output:
                response_text = prefix_text + response_text

            # Check if we need to publish a result bilaterally
            if req.instructions.enable_lateral_stream:
                # Push to redis
                # TODO: Add a debouncing logic to avoid spamming redis
                # Update: Deboucing is probably not needed as this doesn't add any noticable latency to the response
                key_name = f"inferix:llm:{req.context.id}:{req.context.name}"
                hash_map = {
                    "done": str(chunk.done),
                    "response": response_text,
                }
                r = RedisClient.get_client()
                await r.hset(key_name, mapping=hash_map)  # type: ignore

            # Check if stream is done
            if chunk.done:
                # Add the suffix if we need to include it in the response
                if req.instructions.add_suffix is not None and req.instructions.add_suffix.include_in_output:
                    response_text = response_text + suffix_text

                # Gather important stats
                eval_count = chunk.eval_count
                prompt_eval_count = chunk.prompt_eval_count
                created_at = chunk.created_at

                # All looks good. We can exit the loop
                break

        # Async loop is over

        # We will need to clean the output a bit if it was a function call request
        try:
            if "FUNC_CALL" in response_text:
                # First try to sanitize the text
                response_text = sanitize_json_text(response_text)

                # Attempt to serialize it
                f = dict(json.loads(response_text))
        except json.JSONDecodeError:
            # We want to retry inference if the response was not JSON serializable
            continue

        # Prepare usage object
        usage: CompletionUsage = CompletionUsage(
            completion_tokens=eval_count or 0,
            prompt_tokens=prompt_eval_count or 0,
            total_tokens=(eval_count or 0) + (prompt_eval_count or 0),
        )

        # Prepare the response
        res = RunInferenceResponse(
            id="1",
            created=to_unix_timestamp(created_at),
            usage=usage,
            model=req.model,
            object="chat.completion",
            choices=[prepare_chat_completion_message(response_text)],
        )

        # Return the response
        return res
