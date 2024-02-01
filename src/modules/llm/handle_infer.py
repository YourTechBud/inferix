import json
from clients.redis import RedisClient
from fastapi import HTTPException

from models.http import StandardResponse
from models.infer import RunInferenceInstructions, RunInferenceRequest, RunInferenceResponse
from models.openai import ChatCompletionRequestMessage, CompletionUsage
from utils.logger import logger

from .types import OllamaRequest, OllamaRequestOptions
from .helpers import add_message_for_fn_call, prepare_chat_completion_message, sanitize_json_text, to_unix_timestamp
from .ollama import call_ollama_stream
from .prompts import get_prompt


async def handle_infer(req: RunInferenceRequest) -> RunInferenceResponse:
    # We don't support streaming for now
    if req.stream:
        raise HTTPException(status_code=501, detail=StandardResponse(message="Streaming is not supported. Consider using lateral streaming instead.").to_json())

    if req.instructions == None:
        req.instructions = RunInferenceInstructions()

    # Number of messages must be 2 when using conversation instruction. The role of the first message must be system
    # while the role of the second message must be user
    if req.instructions.conversation is not None:
        if len(req.messages) != 2:
            raise HTTPException(
                status_code=400,
                detail=StandardResponse(message="Conversation instruction requires exactly 2 messages").to_json(),
            )

        if req.messages[0].role != "system":
            raise HTTPException(
                status_code=400,
                detail=StandardResponse(message="First message must be system").to_json(),
            )

        if req.messages[1].role != "user":
            raise HTTPException(
                status_code=400,
                detail=StandardResponse(message="Second message must be user").to_json(),
            )

    # Check if we need to load conversation history
    if req.instructions.conversation is not None and req.instructions.conversation.load_key is not None:
        # Load the conversation from redis
        key_name = f"inferix:llm:conversation:{req.context.id}:{req.instructions.conversation.load_key}"
        r = RedisClient.get_client()
        conversation = await r.zrange(key_name, 0, -1)  # type: ignore

        # Print warning if conversation is empty
        if len(conversation) == 0:
            logger.warning(f"Conversation at {key_name} is empty")

        # Time to add this conversation to the messages between index 0 and 1
        for i, c in enumerate(conversation):
            # Split the message into role and content
            role, content = c.split(":::")

            # Check if value of role is euqal to current assitant. We will store the message with role "user" if it isn't
            role = "assistant" if role == req.instructions.conversation.assistant_name else "user"

            # Add the message
            req.messages.insert(
                1 + i,
                ChatCompletionRequestMessage(
                    content=content,
                    role=role,
                ),
            )

    # Add functions as a system prompt if needed
    add_message_for_fn_call(req.functions, req.messages)

    # Check if we need to stream the response laterally
    if req.instructions.enable_lateral_stream:
        req.stream = True

    # First prepare the prompt
    raw_prompt = get_prompt(req.model, req.messages)

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
                key_name = f"inferix:llm:result:{req.context.id}:{req.context.key}"
                hash_map = {
                    "done": str(chunk.done),
                    "response": response_text,
                }
                r = RedisClient.get_client()
                await r.hset(key_name, mapping=hash_map)  # type: ignore
                await r.expire(key_name, 60 * 10) # Expire in 10 minutes

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

        # Check if we need to store the conversation
        if req.instructions.conversation is not None and req.instructions.conversation.store_key is not None:
            conversation_text = response_text

            # Check if we need to add a prefix
            if req.instructions.add_prefix is not None and not req.instructions.add_prefix.include_in_output:
                conversation_text = prefix_text + conversation_text

            # Check if we need to add a suffix
            if req.instructions.add_suffix is not None and not req.instructions.add_suffix.include_in_output:
                conversation_text = conversation_text + suffix_text

            # Push the user prompt and assistant response to redis
            # For the user message, we'll get the content of the last message since that is the message that was sent by the user
            key_name = f"inferix:llm:conversation:{req.context.id}:{req.instructions.conversation.store_key}"
            logger.info(f"Storing conversation at {key_name}")
            r = RedisClient.get_client()

            if req.instructions.conversation.store_entire_history:
                load_key_name = f"inferix:llm:conversation:{req.context.id}:{req.instructions.conversation.load_key}"
                conversation = await r.zrange(load_key_name, 0, -1)  # type: ignore
                conversation.append(f"user:::{req.messages[len(req.messages) -1].content}")
                conversation.append(f"{req.instructions.conversation.assistant_name}:::{conversation_text}")
                await RedisClient.store_as_sorted_set(key_name, conversation)
            else:
                await RedisClient.append_to_sorted_set(key_name, f"user:::{req.messages[len(req.messages) -1].content}")
                await RedisClient.append_to_sorted_set(key_name, f"{req.instructions.conversation.assistant_name}:::{conversation_text}")

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


async def handle_delete_conversations_by_context(ctx_id: str) -> StandardResponse:
    # This is the key we'll use to retrieve the result
    key_name = f"inferix:llm:conversation:{ctx_id}:*"

    # Get the redis client
    redis_client = await RedisClient.get_client()

    # Scan all the keys which match the pattern
    keys = redis_client.scan_iter(match=key_name)

    # Delete the keys
    async for key in keys:
        await redis_client.delete(key)

    return StandardResponse(message=f"Deleted conversations for {ctx_id}")


async def handle_delete_conversation_by_key(ctx_id: str, key: str) -> StandardResponse:
    # This is the key we'll use to retrieve the result
    key_name = f"inferix:llm:conversation:{ctx_id}:{key}"

    # Get the redis client
    redis_client = await RedisClient.get_client()

    # Delete the key
    await redis_client.delete(key_name)

    return StandardResponse(message=f"Deleted conversation for {ctx_id}:{key}")
