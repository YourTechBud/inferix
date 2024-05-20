mod embedding;
mod inference;
mod streaming;
mod types;

use serde::Deserialize;

#[derive(Debug, Deserialize)]
pub struct Ollama {
    base_url: String,
}

impl Ollama {
    pub fn new(config: serde_json::Value) -> Self {
        return serde_json::from_value(config).unwrap();
    }
}
