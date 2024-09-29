use async_trait::async_trait;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        apis::openai::types::CreateChatCompletionResponse,
        drivers::InferenceDriver,
        types::{InferenceOptions, InferenceRequest, InferenceResponseSync, InferenceStats},
    },
    utils,
};

use super::OpenAI;

#[async_trait]
impl InferenceDriver for OpenAI {
    async fn run_inference(
        &self,
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> Result<InferenceResponseSync, AppError> {
        // Prepare openai request
        let req = OpenAI::prepare_request_message(req, options);

        // Fire the request
        let res = self
            .prepare_reqwest_builder(&req)
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
            tracing::error!("Response error: {}", e);
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
                tracing::error!("Response parsing error: {}", e);
                return AppError::BadRequest(StandardErrorResponse::new(
                    "Unable to parse server response".to_string(),
                    e.to_string(),
                ));
            })?;

        return Ok(InferenceResponseSync {
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
}
