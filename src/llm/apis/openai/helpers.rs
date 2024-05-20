use async_stream::{stream, try_stream};
use axum::response::sse::Event;
use futures_core::stream::BoxStream;
use futures_util::StreamExt;

use crate::{
    http::AppError,
    llm::{
        self,
        types::{FinishReason, InferenceResponseStream},
    },
    utils,
};

use super::{
    types, ChatCompletionMessageToolCall, ChatCompletionResponseMessage,
    ChatCompletionStreamResponseDelta, CreateChatCompletionStreamResponse, ResponseChoice,
    StreamingResponseChoice, Usage,
};

pub fn prepare_chat_completion_message(
    response: &str,
    fn_call: Option<llm::types::llm::FunctionCall>,
    tool_selection: types::ToolSelection,
) -> ResponseChoice {
    // Prepare the function call and tools varibles
    let mut fc = None;
    let mut tools = None;

    match tool_selection {
        types::ToolSelection::Function => {
            if let Some(t) = fn_call {
                fc = Some(types::FunctionCall {
                    name: t.name.clone(),
                    arguments: serde_json::to_string(&t.parameters).unwrap(),
                });
            }
        }

        types::ToolSelection::Tools => {
            if let Some(tc) = fn_call {
                tools = Some(vec![ChatCompletionMessageToolCall {
                    id: tc.name.clone(),
                    tool_type: types::ToolType::Function,
                    function: types::FunctionCall {
                        name: tc.name.clone(),
                        arguments: serde_json::to_string(&tc.parameters).unwrap(),
                    },
                }])
            }
        }

        _ => {}
    }

    let finish_reason = if fc.is_some() {
        FinishReason::FunctionCall
    } else {
        FinishReason::Stop
    };

    return ResponseChoice {
        message: ChatCompletionResponseMessage {
            content: Some(response.to_string()),
            role: "assistant".to_string(),
            function_call: fc,
            tool_calls: tools,
        },
        index: 0,
        finish_reason: finish_reason,
    };
}

pub fn prepare_chat_completion_message_chunk(
    input: BoxStream<Result<InferenceResponseStream, AppError>>,
) -> BoxStream<Result<types::CreateChatCompletionStreamResponse, AppError>> {
    return try_stream! {
        for await chunk in input {
            let chunk = chunk.map_err(|e| {
                return e;
            })?;
            yield CreateChatCompletionStreamResponse {
                id: "inferix".to_string(),
                model: chunk.model.clone(),
                choices: vec![StreamingResponseChoice{
                    index: 0,
                    delta: ChatCompletionStreamResponseDelta {
                        content: Some(chunk.response),
                        role: Some("assistant".to_string()),

                        // TODO: Add support for function call while streaming
                        function_call: None,
                        tool_calls: None,
                    },
                    finish_reason: chunk.finish_reason,
                }],
                created: utils::convert_to_unix_timestamp(&chunk.created_at),
                object: "chat.completion.chunk".to_string(),
                system_fingerprint: None,
                usage: match chunk.stats {
                    Some(stats) => Some(Usage {
                        prompt_tokens: stats.prompt_eval_count.unwrap_or(0),
                        completion_tokens: stats.eval_count.unwrap_or(0),
                        total_tokens: stats.eval_count.unwrap_or(0) + stats.prompt_eval_count.unwrap_or(0),
                    }),
                    None => None,
                },
            }
        }
    }
    .boxed();
}

pub fn prepare_sse_response_stream<'a>(
    input: BoxStream<'a, Result<types::CreateChatCompletionStreamResponse, AppError>>,
) -> BoxStream<'a, Result<Event, AppError>> {
    let stream = stream! {
        for await chunk in input {
            let chunk = match chunk {
                Ok(chunk) => serde_json::to_string(&chunk).unwrap(),
                Err(e) => e.to_string(),
            };
            yield Ok(Event::default().data(&chunk));
        }
    };

    return stream.boxed();
}
