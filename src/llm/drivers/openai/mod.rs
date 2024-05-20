mod embedding;
mod helpers;
mod inference;
mod streaming;

use serde::Deserialize;


#[derive(Debug, Deserialize)]
pub struct OpenAI {
    base_url: String,
}

impl OpenAI {
    pub fn new(config: serde_json::Value) -> Self {
        return serde_json::from_value(config).unwrap();
    }
}
