from typing import Callable, List
from fastapi import HTTPException
import jinja2
from models.http import StandardResponse

from models.openai import ChatCompletionRequestMessage


class PromptTemplate:
    def __init__(
        self, tmpl: Callable[[List[ChatCompletionRequestMessage]], str], stop: List[str]
    ) -> None:
        self.template = tmpl
        self.stop = stop if len(stop) > 0 else [""]

    def get_prompt(self, messages: List[ChatCompletionRequestMessage]) -> str:
        return self.template(messages)


# Prompt template for chatml
def get_chatml_tmpl() -> PromptTemplate:
    environment = jinja2.Environment()
    chatml_tmpl: str = """{% for msg in messages %}<|im_start|>{{ msg.role }}
{{ msg.content }}<|im_end|>
{% endfor %}<|im_start|>assistant"""
    template = environment.from_string(chatml_tmpl)

    def fn(messages: List[ChatCompletionRequestMessage]) -> str:
        prompt = template.render(messages=messages)
        print("================")
        print("Generated prompt:", prompt)
        print("================")
        return prompt

    return PromptTemplate(fn, ["<|im_start|>", "<|im_end|>"])


# Prompt template for chatml
def get_vicuna1_1_tmpl() -> PromptTemplate:
    def fn(messages: List[ChatCompletionRequestMessage]) -> str:
        # Select the first message where role is user
        user_message = next((msg for msg in messages if msg.role == "user"), None)
        if user_message is None:
            raise HTTPException(
                status_code=400,
                detail=StandardResponse(message="User message needs to be present").to_json(),
            )
        prompt = f"USER: {user_message.content}\nASSISTANT:"
        print("================")
        print("Generated prompt:", prompt)
        print("================")
        return prompt

    return PromptTemplate(fn, ["</s>"])


# Prompt template for chatml
def get_user_assistant_newlines_tmpl() -> PromptTemplate:
    def fn(messages: List[ChatCompletionRequestMessage]) -> str:
        # Select the first message where role is user
        user_message = next((msg for msg in messages if msg.role == "user"), None)
        if user_message is None:
            raise HTTPException(
                status_code=400,
                detail=StandardResponse(message="User message needs to be present").to_json(),
            )
        prompt = f"### User:\n{user_message.content}\n\n### Assistant:"
        print("================")
        print("Generated prompt:", prompt)
        print("================")
        return prompt

    return PromptTemplate(fn, ["</s>"])


chatml_tmpl = get_chatml_tmpl()
vicuna1_1_tmpl = get_vicuna1_1_tmpl()
user_assistant_newlines_tmpl = get_user_assistant_newlines_tmpl()


def get_prompt(model: str, messages: List[ChatCompletionRequestMessage]) -> str:
    if "solar" in model:
        return user_assistant_newlines_tmpl.get_prompt(messages)
    
    if model.startswith("nous-capybara"):
        return vicuna1_1_tmpl.get_prompt(messages)
    
    return chatml_tmpl.get_prompt(messages)
