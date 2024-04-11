use axum::Json;

use crate::{
    embedding::drivers::{run_inference, Request},
    http::AppError,
};

pub async fn handle_embed_request(
    Json(req): Json<types::EmbeddingRequest>,
) -> Result<Json<Vec<Vec<f64>>>, AppError> {
    let req = Request {
        model: None,
        trucate: req.trucate,
        inputs: req.inputs,
    };

    let res = run_inference(req).await?;

    // Convert the response to a format that can be returned
    let mut embeddings = Vec::with_capacity(res.data.len());
    for embedding in res.data {
        embeddings.push(embedding.embedding);
    }
    return Ok(Json(embeddings));
}

mod types {
    use serde::{Deserialize, Serialize};

    use crate::embedding::drivers::EmbeddingInput;

    #[derive(Serialize, Deserialize)]
    pub struct EmbeddingRequest {
        pub trucate: Option<bool>,
        pub inputs: EmbeddingInput,
    }
}
