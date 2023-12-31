import argparse

import uvicorn


if __name__ == "__main__":
    # Create an ArgumentParser object
    parser = argparse.ArgumentParser(description="Program to host LLMs.")

    # Add all server arguments
    parser.add_argument(
        "-p", "--port", type=int, default=8000, help="The port to start the server on."
    )
    parser.add_argument(
        "-w",
        "--workers",
        type=int,
        default=1,
        help="Set number of worker processes to run.",
    )
    parser.add_argument(
        "--reload", type=bool, default=False, help="Enable auto reload."
    )

    # Parse the arguments
    args = parser.parse_args()

    uvicorn.run(
        "server:app",
        host="0.0.0.0",
        port=args.port,
        reload=args.reload,
        workers=args.workers,
    )