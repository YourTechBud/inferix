use async_trait::async_trait;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        drivers::EmbeddingDriver,
        types::{EmbeddingRequest, EmbeddingResponse},
    },
};

use super::Ollama;

#[async_trait]
impl EmbeddingDriver for Ollama {
    async fn create_embedding(
        &self,
        _req: &EmbeddingRequest,
    ) -> Result<EmbeddingResponse, AppError> {
        // TODO: Ollama does support embeddings. We should only deny if the user
        // explicitly forbids using that model for embeddings.
        return Err(AppError::BadRequest(StandardErrorResponse::new(
            "Ollama driver does not support streaming inference".to_string(),
            "function_not_supported".to_string(),
        )));
    }
}
