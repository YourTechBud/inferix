from typing import Any, Optional
from pydantic import BaseModel


class StandardResponse(BaseModel):
  message: str
  error: Optional[Any] = None

  def to_json(self):
    return self.__dict__