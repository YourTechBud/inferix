pub mod types;

mod helpers;

use axum::{
    response::{sse::KeepAlive, IntoResponse, Response, Sse},
    Json,
};

use crate::{
    http::AppError,
    llm::{
        self,
        drivers::{self, embedding::create_embeddings},
    },
    utils,
};

use types::*;

pub async fn handle_embedding_request(
    Json(req): Json<types::OpenAIEmbeddingRequest>,
) -> Result<Json<OpenAIEmbeddingResponse>, AppError> {
    let req = llm::types::EmbeddingRequest {
        model: req.model,
        trucate: None, // We will let the driver take care of this
        inputs: req.input,
    };

    // Run the inference
    let res = create_embeddings(req).await?;

    // Let's prepare the response
    // We gotta do it the hard way cause I don't know any better.
    let mut data = Vec::with_capacity(res.data.len());
    for embedding in res.data {
        data.push(types::Embedding {
            index: embedding.index,
            embedding: embedding.embedding,
            object: types::EmbeddingObject::Embedding,
        });
    }
    let res = types::OpenAIEmbeddingResponse {
        model: res.model,
        data,
        usage: types::EmbeddingUsage {
            prompt_tokens: res.usage.prompt_tokens,
            total_tokens: res.usage.total_tokens,
        },
        object: types::OpenAIEmbeddingResponseObject::List,
    };

    // Return the response
    return Ok(Json(res));
}

pub async fn handle_chat_completion(Json(req): Json<CreateChatCompletionRequest>) -> Result<Response, AppError> {
    // Convert request messages to the format expected by the driver
    let messages = req
        .messages
        .iter()
        .map(|m| llm::types::InferenceMessage {
            role: m.role.clone(),
            content: m.content.clone(),
        })
        .collect();

    // Select the tools provided from the `tools` or `functions` field in the request
    // We will give priority to the `tools` field since the `functions` field is deprecated.
    // TODO: Make this more efficient and look good
    let mut tool_selection = types::ToolSelection::None;

    let mut tools: Option<Vec<llm::types::Tool>> = if let Some(tools) = req.tools {
        tool_selection = types::ToolSelection::Tools;
        Some(
            tools
                .iter()
                .map(|t| llm::types::Tool {
                    name: t.function.name.clone(),
                    description: t.function.description.clone(),
                    args: t
                        .function
                        .parameters
                        .clone()
                        .unwrap_or(FunctionParameters {
                            additional_properties: serde_json::json!({}),
                        })
                        .additional_properties,
                    tool_type: llm::types::ToolType::Function,
                })
                .collect(),
        )
    } else {
        None
    };

    // Only checks if any functions are provided if tools is still None
    if tools.is_none() {
        if let Some(functions) = req.functions {
            tool_selection = types::ToolSelection::Function;
            tools = Some(
                functions
                    .iter()
                    .map(|f| llm::types::Tool {
                        name: f.name.clone(),
                        description: f.description.clone(),
                        args: f.parameters.additional_properties.clone(),
                        tool_type: llm::types::ToolType::Function,
                    })
                    .collect(),
            )
        }
    }

    let inference_request = llm::types::InferenceRequest::new(req.model, messages, tools);
    let inference_options =
        llm::types::InferenceOptions::new(req.top_p, None, req.max_tokens, req.temperature);

    match req.stream {
        Some(true) => {
            let res =
                drivers::inference::run_streaming_inference(inference_request, inference_options)?;

            // Prepare and return the response
            let converted_stream = helpers::prepare_chat_completion_message_chunk(res);
            let stream = helpers::prepare_sse_response_stream(converted_stream);

            return Ok(Sse::new(stream)
                .keep_alive(KeepAlive::default())
                .into_response());
        }

        _ => {
            // Run the inference
            let res = drivers::inference::run_inference(inference_request, inference_options)
                .await?;

            // Prepare and return the response
            let res = CreateChatCompletionResponse {
                id: "inferix".to_string(),
                created: utils::convert_to_unix_timestamp(&res.created_at),
                model: res.model.clone(),
                object: "chat.completion".to_string(),
                system_fingerprint: None,
                usage: CompletionUsage {
                    completion_tokens: res.stats.eval_count.unwrap_or(0),
                    prompt_tokens: res.stats.prompt_eval_count.unwrap_or(0),
                    total_tokens: res.stats.eval_count.unwrap_or(0)
                        + res.stats.prompt_eval_count.unwrap_or(0),
                },
                choices: vec![helpers::prepare_chat_completion_message(
                    &res.response,
                    res.fn_call,
                    tool_selection,
                )],
            };

            return Ok(Json(res).into_response());
        }
    }
}
