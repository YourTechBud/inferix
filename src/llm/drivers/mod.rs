pub mod embedding;
pub mod inference;

mod helpers;

mod ollama;
mod openai;
mod tei;

use async_trait::async_trait;
use enum_dispatch::enum_dispatch;
use futures_core::stream::BoxStream;
use once_cell::sync::OnceCell;
use serde::Deserialize;
use std::{collections::HashMap, fmt::Debug};

use self::ollama::Ollama;
use self::openai::OpenAI;
use self::tei::TextEmbeddingsInference;
use crate::{http::AppError, llm::types::*};

#[async_trait]
#[enum_dispatch(Driver)]
pub trait InferenceDriver: Send + Sync + Debug {
    async fn run_inference(
        &self,
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> Result<InferenceResponseSync, AppError>;
}

#[enum_dispatch(Driver)]
pub trait StreamingInference: Send + Sync + Debug {
    fn run_streaming_inference(
        &self,
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> Result<BoxStream<Result<InferenceResponseStream, AppError>>, AppError>;
}

#[async_trait]
#[enum_dispatch(Driver)]
pub trait EmbeddingDriver: Send + Sync + Debug {
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

    #[serde(rename = "openai")]
    OpenAI,

    #[serde(rename = "tei")]
    TEI,
}

#[enum_dispatch]
#[derive(Debug)]
pub enum Driver {
    OpenAI,
    Ollama,
    TextEmbeddingsInference,
}

pub static DRIVERS: OnceCell<HashMap<String, Driver>> = OnceCell::new();

pub fn init(drivers: Vec<DriverConfig>) {
    let mut d = HashMap::new();
    for driver in drivers {
        match driver.driver_type {
            DriverType::Ollama => {
                d.insert(driver.name, ollama::Ollama::new(driver.config).into());
            }

            DriverType::TEI => {
                d.insert(
                    driver.name,
                    tei::TextEmbeddingsInference::new(driver.config).into(),
                );
            }

            DriverType::OpenAI => {
                d.insert(driver.name, openai::OpenAI::new(driver.config).into());
            }
        }
    }

    DRIVERS.set(d).unwrap();
}
