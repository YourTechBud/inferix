# Inferix

## Description

Inferix is a wrapper on top of [Ollama](https://ollama.ai/). It aims to expose an OpenAI compatible API with some extra goodies. It currently supports:

1. OpenAI compatible RestAPI on top of Ollama.
2. LLM powered function calling.
3. Ability to stream response laterally.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)

## Installation

To install and run Inferix, follow these steps:

1. Clone the repository: `git clone https://github.com/YourTechBud/inferix.git`
2. Navigate to the project directory: `cd inferix`
3. Install the dependencies: `poetry install`

## Usage

To start Inferix, run the following command:

```bash
poetry run start
```

- Open [http://localhost:8000/docs](http://localhost:8000/docs) to open swagger ui
- Open [http://localhost:8000/openapi.json](http://localhost:8000/openapi.json) to access the raw openapi spec