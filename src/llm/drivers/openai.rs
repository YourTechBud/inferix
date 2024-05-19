use async_trait::async_trait;
use reqwest::Client;
use serde::Deserialize;

use super::Driver;
use crate::llm::apis::openai::types::{
    ChatCompletionRequestMessage, CreateChatCompletionRequest, CreateChatCompletionResponse,
};
use crate::{
    http::{AppError, StandardErrorResponse},
    llm::types::*,
    utils,
};

#[derive(Debug, Deserialize)]
pub struct OpenAI {
    base_url: String,
}

impl OpenAI {
    pub fn new(config: serde_json::Value) -> Self {
        return serde_json::from_value(config).unwrap();
    }
}

#[async_trait]
impl Driver for OpenAI {
    async fn run_inference(
        &self,
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> Result<InferenceResponse, AppError> {
        // Prepare openai request
        let req = CreateChatCompletionRequest {
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
            functions: None,        // TODO: Add support for this
            tool_choice: None,
            tools: None,            // TODO: Add support for this
            response_format: None,  // TODO: Add support for this
            n: Some(1),
            seed: None,
            stream: Some(false),
        };

        // Fire the request
        let client = Client::new();
        let res = client
            .post(&format!("{}/chat/completions", self.base_url))
            .json(&req)
            .send()
            .await
            .map_err(|e| {
                return AppError::InternalServerError(StandardErrorResponse::new(
                    format!("Unable to make request to OpenAI server: {}", e),
                    "openai_call_error".to_string(),
                ));
            })?;

        // Read the response
        let status = res.status();
        let is_success = status.is_success();
        let response_text = res.text().await.map_err(|e| {
            eprintln!("Response error: {}", e);
            return AppError::BadRequest(StandardErrorResponse::new(
                "Unable to get response from driver".to_string(),
                e.to_string(),
            ));
        })?;

        if !is_success {
            let message = serde_json::from_str(&response_text).unwrap_or_else(|e| {
                StandardErrorResponse::new("Unable to service request".to_string(), e.to_string())
            });
            return Err(AppError::BadRequest(message));
        }

        let res: CreateChatCompletionResponse =
            serde_json::from_str(&response_text).map_err(|e| {
                eprintln!("Response parsing error: {}", e);
                return AppError::BadRequest(StandardErrorResponse::new(
                    "Unable to parse server response".to_string(),
                    e.to_string(),
                ));
            })?;

        return Ok(InferenceResponse {
            model: res.model,
            response: res.choices[0]
                .message
                .content
                .clone()
                .unwrap_or("".to_string()),
            fn_call: None, // This won't necessarily be None. You need to implement this/
            stats: InferenceStats {
                eval_count: Some(res.usage.completion_tokens),
                eval_duration: None,
                prompt_eval_count: Some(res.usage.prompt_tokens),
                prompt_eval_duration: None,
                load_duration: None,
                total_duration: None, // TODO: Calculate this youourself
            },
            created_at: utils::convert_to_datetime(res.created),
        });
    }

    async fn create_embedding(
        &self,
        _req: &EmbeddingRequest,
    ) -> Result<EmbeddingResponse, AppError> {
        // TODO: OpenAI does support embeddings. We should only deny if the user
        // explicitly forbids using that model for embeddings.
        return Err(AppError::BadRequest(StandardErrorResponse::new(
            "OpenAI driver does not support embeddings".to_string(),
            "function_not_supported".to_string(),
        )));
    }
}
