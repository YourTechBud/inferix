use axum::Json;

use crate::{
    http::AppError,
    llm::{drivers::embedding::create_embeddings, types::EmbeddingRequest},
};

pub async fn handle_embed_request(
    Json(req): Json<types::EmbeddingRequest>,
) -> Result<Json<Vec<Vec<f64>>>, AppError> {
    let req = EmbeddingRequest {
        model: "default".to_string(),
        trucate: req.trucate,
        inputs: req.inputs,
    };

    let res = create_embeddings(req).await?;

    // Convert the response to a format that can be returned
    let mut embeddings = Vec::with_capacity(res.data.len());
    for embedding in res.data {
        embeddings.push(embedding.embedding);
    }
    return Ok(Json(embeddings));
}

mod types {
    use serde::{Deserialize, Serialize};

    use crate::llm::types::EmbeddingInput;

    #[derive(Serialize, Deserialize)]
    pub struct EmbeddingRequest {
        pub trucate: Option<bool>,
        pub inputs: EmbeddingInput,
    }
}
