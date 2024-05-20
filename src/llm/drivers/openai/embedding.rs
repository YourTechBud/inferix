use async_trait::async_trait;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        drivers::EmbeddingDriver,
        types::{EmbeddingRequest, EmbeddingResponse},
    },
};

use super::OpenAI;

#[async_trait]
impl EmbeddingDriver for OpenAI {
    async fn create_embedding(
        &self,
        _req: &EmbeddingRequest,
    ) -> Result<EmbeddingResponse, AppError> {
        return Err(AppError::InternalServerError(StandardErrorResponse::new(
            "Inference is not supported".to_string(),
            "function_not_supported".to_string(),
        )));
    }
}
