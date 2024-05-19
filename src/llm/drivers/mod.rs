pub mod inference;
pub mod embedding;

mod ollama;
mod tei;

use async_trait::async_trait;
use once_cell::sync::OnceCell;
use serde::Deserialize;
use std::{collections::HashMap, fmt::Debug};

use crate::{http::AppError, llm::types::*};

#[async_trait]
pub trait Driver: Send + Sync + Debug {
    async fn run_inference(
        &self,
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> Result<InferenceResponse, AppError>;

    async fn create_embedding(&self, req: &EmbeddingRequest)
        -> Result<EmbeddingResponse, AppError>;
}

#[derive(Deserialize)]
pub struct DriverConfig {
    pub name: String,

    #[serde(rename = "type")]
    pub driver_type: DriverType,
    pub config: serde_json::Value,
}

#[derive(Deserialize)]
pub enum DriverType {
    #[serde(rename = "ollama")]
    Ollama,

    #[serde(rename = "tei")]
    TEI,
}

pub static DRIVERS: OnceCell<HashMap<String, Box<dyn Driver>>> = OnceCell::new();

pub fn init(drivers: Vec<DriverConfig>) {
    let mut d = HashMap::new();
    for driver in drivers {
        match driver.driver_type {
            DriverType::Ollama => {
                d.insert(
                    driver.name,
                    Box::new(ollama::Ollama::new(driver.config)) as Box<dyn Driver>,
                );
            },
            DriverType::TEI => {
                d.insert(
                    driver.name,
                    Box::new(tei::TextEmbeddingsInference::new(driver.config)) as Box<dyn Driver>,
                );
            },
        }
    }
    // DRIVER
    //     .set(Box::new(ollama::Ollama::new(
    //         "localhost".to_string(),
    //         "11434".to_string(),
    //     )))
    //     .unwrap();

    DRIVERS.set(d).unwrap();
}
