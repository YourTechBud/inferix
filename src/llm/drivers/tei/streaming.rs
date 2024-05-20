use futures_core::stream::BoxStream;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        drivers::StreamingInference,
        types::{InferenceOptions, InferenceRequest, InferenceResponseStream},
    },
};

use super::TextEmbeddingsInference;

impl StreamingInference for TextEmbeddingsInference {
    fn run_streaming_inference(
        &self,
        _req: &InferenceRequest,
        _options: &InferenceOptions,
    ) -> Result<BoxStream<Result<InferenceResponseStream, AppError>>, AppError> {
        // TODO: Ollama does support embeddings. We should only deny if the user
        // explicitly forbids using that model for embeddings.
        return Err(AppError::BadRequest(StandardErrorResponse::new(
            "TEI driver does not support streaming inference".to_string(),
            "function_not_supported".to_string(),
        )));
    }
}
