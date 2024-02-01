from typing import List

from models.list_models import ListModelsResponse
from modules.llm.ollama import list_models


async def handle_list_models() -> ListModelsResponse:
    return await list_models()