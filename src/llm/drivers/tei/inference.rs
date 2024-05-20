use async_trait::async_trait;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        drivers::InferenceDriver,
        types::{InferenceOptions, InferenceRequest, InferenceResponseSync},
    },
};

use super::TextEmbeddingsInference;

#[async_trait]
impl InferenceDriver for TextEmbeddingsInference {
    async fn run_inference(
        &self,
        _: &InferenceRequest,
        _: &InferenceOptions,
    ) -> Result<InferenceResponseSync, AppError> {
        return Err(AppError::InternalServerError(StandardErrorResponse::new(
            "Inference is not supported".to_string(),
            "function_not_supported".to_string(),
        )));
    }
}
