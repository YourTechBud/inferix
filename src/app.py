import argparse
import uvicorn

from fastapi import FastAPI

from modules.llm.routes import router as llm_router
from clients import create_clients, destroy_clients

app = FastAPI(title="Inferix", description="OpenAI compatible API with extra goodies.", on_startup=[create_clients], on_shutdown=[destroy_clients], logger=True)

app.include_router(llm_router)


def start():
    # Create an ArgumentParser object
    parser = argparse.ArgumentParser(description="Program to host LLMs.")

    # Add all server arguments
    parser.add_argument("-p", "--port", type=int, default=8000, help="The port to start the server on.")
    parser.add_argument(
        "-w",
        "--workers",
        type=int,
        default=1,
        help="Set number of worker processes to run.",
    )
    parser.add_argument("--reload", type=bool, default=False, help="Enable auto reload.")

    # Parse the arguments
    args = parser.parse_args()

    uvicorn.run(
        "app:app",
        host="0.0.0.0",
        port=args.port,
        reload=args.reload,
        workers=args.workers,
        log_level="info",
    )


if __name__ == "__main__":
    start()
