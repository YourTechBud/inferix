use async_stream::{stream, try_stream};
use futures_util::stream::Stream;
use reqwest::{Client, RequestBuilder};
use reqwest_eventsource::{Event, EventSource};

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        apis::openai::types::{
            ChatCompletionRequestMessage, CreateChatCompletionRequest,
            CreateChatCompletionStreamResponse,
        },
        types::{InferenceOptions, InferenceRequest, InferenceResponseStream, InferenceStats},
    },
    utils,
};

use super::OpenAI;

impl OpenAI {
    pub fn prepare_request_message(
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> CreateChatCompletionRequest {
        return CreateChatCompletionRequest {
            model: req.model.clone(),
            messages: req
                .messages
                .iter()
                .map(|m| ChatCompletionRequestMessage {
                    role: m.role.clone(),
                    content: m.content.clone(),
                    name: None,
                    tool_call_id: None,
                    function_call: None,
                    tool_calls: None,
                })
                .collect(),
            max_tokens: options.num_ctx,
            temperature: options.temperature,
            top_p: options.top_p,
            frequency_penalty: None,
            presence_penalty: None,
            stop: None,
            function_call: None,
            functions: None, // TODO: Add support for this
            tool_choice: None,
            tools: None,           // TODO: Add support for this
            response_format: None, // TODO: Add support for this
            n: Some(1),
            seed: None,
            stream: Some(false),
        };
    }

    pub fn prepare_reqwest_builder(&self, req: &CreateChatCompletionRequest) -> RequestBuilder {
        let client = Client::new();
        return client
            .post(&format!("{}/chat/completions", self.base_url))
            .json(&req);
    }

    pub fn prepare_event_stream(input: EventSource) -> impl Stream<Item = Result<Event, AppError>> {
        return stream! {
            for await event in input {
                match event {
                    Ok(Event::Open) => {},
                    Ok(Event::Message(message)) => {
                        if message.data == "[DONE]" {
                            yield Err(AppError::StreamEndedd);
                            continue;
                        }
                        yield Ok(Event::Message(message))
                    },
                    Err(e) => {
                        if e.to_string() == "Stream ended" {
                            continue;
                        } else {
                            yield Err(AppError::BadRequest(StandardErrorResponse::new(
                                format!("Unable to get response from driver: {}", e),
                                "openai_stream_response_error".to_string(),
                            )))
                        }
                    },
                }
            }
        };
    }

    pub fn convert_to_inference_response_stream<S: Stream<Item = Result<Event, AppError>>>(
        input: S,
    ) -> impl Stream<Item = Result<InferenceResponseStream, AppError>> {
        let a = try_stream! {
            for await event in input {
                let event = event?;
                match event {
                    Event::Message(message) => {
                        let chunk: CreateChatCompletionStreamResponse = serde_json::from_str(&message.data).map_err(|e| {
                            return AppError::BadRequest(StandardErrorResponse::new(
                                    format!("Unable to parse server response: {}; Data: {}", e, message.data),
                                    "openai_parse_error".to_string(),
                                ));
                            })?;

                        yield InferenceResponseStream {
                            model: chunk.model,
                            response: chunk.choices[0]
                                .delta
                                .content
                                .clone()
                                .unwrap_or("".to_string()),
                            stats: match chunk.usage {
                                Some(usage) => Some(InferenceStats {
                                    eval_count: Some(
                                        usage.total_tokens - usage.prompt_tokens,
                                    ),
                                    eval_duration: None,
                                    prompt_eval_count: Some(usage.prompt_tokens),
                                    prompt_eval_duration: None,
                                    load_duration: None,
                                    total_duration: None, // TODO: Calculate this youourself
                                }),
                                None => None,
                            },
                            finish_reason: chunk.choices[0].finish_reason.clone(),
                            created_at: utils::convert_to_datetime(chunk.created),
                        };
                    },
                    _ => {
                        unimplemented!("This should not happen")
                    }
                }
            }
        };

        return a;
    }
}
