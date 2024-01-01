from typing import Callable, List
import jinja2

from models.openai import ChatCompletionRequestMessage


class PromptTemplate:
    def __init__(self, tmpl: Callable[[List[ChatCompletionRequestMessage]], str], stop: List[str]) -> None:
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


chatml_tmpl = get_chatml_tmpl()
