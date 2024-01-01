import json
from typing import List, Optional
from pydantic import BaseModel


class OllamaRequestOptions(BaseModel):
    top_p: Optional[float] = 0.9
    top_k: int = 40
    num_ctx: Optional[int] = 4096
    temperature: Optional[float] = 0.2

    def to_json(self):
        return json.dumps(self.__dict__, default=lambda o: o.__dict__)


class OllamaRequest(BaseModel):
    model: str
    prompt: str
    raw: bool = True
    stream: Optional[bool] = False
    options: OllamaRequestOptions

    def to_json(self):
        return json.dumps(self.__dict__, default=lambda o: o.__dict__)


class OllamaResponse(BaseModel):
    model: str
    created_at: str
    response: str
    done: bool
    context: Optional[List[int]] = None
    total_duration: Optional[int] = None
    load_duration: Optional[int] = None
    prompt_eval_count: Optional[int] = None
    prompt_eval_duration: Optional[int] = None
    eval_count: Optional[int] = None
    eval_duration: Optional[int] = None

    def to_json(self):
        return json.dumps(self.__dict__, default=lambda o: o.__dict__)
