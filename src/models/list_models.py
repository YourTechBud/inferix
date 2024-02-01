from typing import List, Optional
from pydantic import BaseModel
from datetime import datetime

class ModelDetails(BaseModel):
    format: str
    family: str
    families: Optional[List[str]]
    parameter_size: str
    quantization_level: str

class Model(BaseModel):
    name: str
    modified_at: datetime
    size: int
    digest: str
    details: ModelDetails

class ListModelsResponse(BaseModel):
    models: List[Model]