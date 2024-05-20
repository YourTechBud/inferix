use serde::Serialize;

use crate::llm::types::EmbeddingInput;

#[derive(Serialize)]
pub struct TextEmbeddingsInferenceRequest<'a> {
    pub input: &'a EmbeddingInput,
}
