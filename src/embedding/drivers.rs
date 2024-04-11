mod tei;

use async_trait::async_trait;
use once_cell::sync::OnceCell;
use serde::{Deserialize, Serialize};
use std::fmt::Debug;

use crate::http::AppError;

#[async_trait]
pub trait Driver: Send + Sync + Debug {
    async fn call(&self, req: &Request) -> Result<Response, AppError>;
}

#[derive(Serialize, Deserialize)]
pub struct Request {
    pub model: Option<String>,
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
pub struct Response {
    pub model: Option<String>,
    pub data: Vec<Embedding>,
    pub usage: Usage,
}

#[derive(Serialize, Deserialize)]
pub struct Embedding {
    pub index: u32,
    pub embedding: Vec<f64>,
}

#[derive(Serialize, Deserialize)]
pub struct Usage {
    pub prompt_tokens: u64,
    pub total_tokens: u64,
}

pub static DRIVER: OnceCell<Box<dyn Driver>> = OnceCell::new();

pub fn init() {
    DRIVER
        .set(Box::new(tei::TextEmbeddingsInference::new(
            "localhost".to_string(),
            "8090".to_string(),
        )))
        .unwrap();
}

pub async fn run_inference(req: Request) -> Result<Response, AppError> {
    // TODO: Get the right driver based on the model. Use the default model if none is provided
    let driver = DRIVER.get().unwrap();
    return driver.call(&req).await;
}
