use serde::{Deserialize, Serialize};

#[derive(Serialize, Deserialize)]
pub struct EmbeddingRequest {
    pub model: String,
    pub trucate: Option<bool>,
    pub inputs: EmbeddingInput,
}

#[derive(Serialize, Deserialize)]
#[serde(untagged)]
pub enum EmbeddingInput {
    Input(String),
    Inputs(Vec<String>),
}

#[derive(Serialize, Deserialize)]
pub struct EmbeddingResponse {
    pub model: Option<String>,
    pub data: Vec<Embedding>,
    pub usage: EmbeddingUsage,
}

#[derive(Serialize, Deserialize)]
pub struct Embedding {
    pub index: u32,
    pub embedding: Vec<f64>,
}

#[derive(Serialize, Deserialize)]
pub struct EmbeddingUsage {
    pub prompt_tokens: u64,
    pub total_tokens: u64,
}