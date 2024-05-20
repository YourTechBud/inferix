use futures_core::stream::BoxStream;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        drivers::StreamingInference,
        types::{InferenceOptions, InferenceRequest, InferenceResponseStream},
    },
};

use super::Ollama;

impl StreamingInference for Ollama {
    fn run_streaming_inference(
        &self,
        _req: &InferenceRequest,
        _options: &InferenceOptions,
    ) -> Result<BoxStream<Result<InferenceResponseStream, AppError>>, AppError> {
        // TODO: Ollama does support embeddings. We should only deny if the user
        // explicitly forbids using that model for embeddings.
        return Err(AppError::BadRequest(StandardErrorResponse::new(
            "Ollama driver does not support embeddings".to_string(),
            "function_not_supported".to_string(),
        )));
    }
}
