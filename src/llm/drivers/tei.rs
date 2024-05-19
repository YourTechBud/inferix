use async_trait::async_trait;
use reqwest::Client;
use serde::Serialize;

use crate::http::{AppError, StandardErrorResponse};

use super::*;

#[derive(Debug, Deserialize)]
pub struct TextEmbeddingsInference {
    base_url: String,
}

impl TextEmbeddingsInference {
    pub fn new(config: serde_json::Value) -> Self {
        return serde_json::from_value(config).unwrap();
    }
}

#[derive(Serialize)]
struct TextEmbeddingsInferenceRequest<'a> {
    input: &'a EmbeddingInput,
}

#[async_trait]
impl Driver for TextEmbeddingsInference {
    async fn run_inference(
        &self,
        _: &InferenceRequest,
        _: &InferenceOptions,
    ) -> Result<InferenceResponse, AppError> {
        return Err(AppError::InternalServerError(StandardErrorResponse::new(
            "Inference is not supported".to_string(),
            "function_not_supported".to_string(),
        )));
    }

    async fn create_embedding(
        &self,
        req: &EmbeddingRequest,
    ) -> Result<EmbeddingResponse, AppError> {
        // Prepare the request
        let req = TextEmbeddingsInferenceRequest { input: &req.inputs };

        // Send the request
        let client = Client::new();
        let res = client
            .post(&format!("{}/embeddings", self.base_url))
            .json(&req)
            .send()
            .await
            .map_err(|e| {
                return AppError::InternalServerError(StandardErrorResponse::new(
                    format!(
                        "Unable to make request to Text Embeddings Inference server: {}",
                        e
                    ),
                    "tei_call_error".to_string(),
                ));
            })?;

        // Parse the response
        let res = res.text().await.map_err(|e| {
            return AppError::InternalServerError(StandardErrorResponse::new(
                format!(
                    "Unable to get response from Text Embeddings Inference server: {}",
                    e
                ),
                "tei_response_error".to_string(),
            ));
        })?;

        // Convert the response to a list of byte array
        let res = serde_json::from_str(&res).map_err(|e| {
            return AppError::InternalServerError(StandardErrorResponse::new(
                format!(
                    "Unable to parse response from Text Embeddings Inference server: {}",
                    e
                ),
                "tei_parse_error".to_string(),
            ));
        })?;

        return Ok(res);
    }
}
