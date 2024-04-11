use async_trait::async_trait;
use reqwest::Client;
use serde::Serialize;

use crate::http::{AppError, StandardErrorResponse};

use super::{Driver, Request, EmbeddingInput, Response};

#[derive(Debug)]
pub struct TextEmbeddingsInference {
    host: String,
    port: String,
}

impl TextEmbeddingsInference {
    pub fn new(host: String, port: String) -> Self {
        return TextEmbeddingsInference { host, port };
    }
}

#[derive(Serialize)]
struct TextEmbeddingsInferenceRequest<'a> {
    input: &'a EmbeddingInput,
}

#[async_trait]
impl Driver for TextEmbeddingsInference {
    async fn call(&self, req: &Request) -> Result<Response, AppError> {
        // Prepare the request
        let req = TextEmbeddingsInferenceRequest {
            input: &req.inputs,
        };

        // Send the request
        let client = Client::new();
        let res = client
            .post(&format!("http://{}:{}/embeddings", self.host, self.port))
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
