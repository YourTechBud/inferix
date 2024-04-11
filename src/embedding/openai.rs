use axum::Json;

use crate::{
    embedding::drivers::{run_inference, Request},
    http::AppError,
};

pub async fn handle_embedding_request(
    Json(req): Json<types::OpenAIEmbeddingRequest>,
) -> Result<Json<types::OpenAIEmbeddingResponse>, AppError> {
    let req = Request {
        model: Some(req.model),
        trucate: None, // We will let the driver take care of this
        inputs: req.input,
    };

    // Run the inference
    let res = run_inference(req).await?;

    // Let's prepare the response
    // We gotta do it the hard way cause we don't want to clone stuff around.
    let mut data = Vec::with_capacity(res.data.len());
    for embedding in res.data {
        data.push(types::Embedding {
            index: embedding.index,
            embedding: embedding.embedding,
            object: types::EmbeddingObject::Embedding,
        });
    }
    let res = types::OpenAIEmbeddingResponse {
        model: res.model,
        data,
        usage: types::Usage {
            prompt_tokens: res.usage.prompt_tokens,
            total_tokens: res.usage.total_tokens,
        },
        object: types::OpenAIEmbeddingResponseObject::List,
    };

    // Return the response
    return Ok(Json(res));
}

mod types {
    use serde::{Deserialize, Serialize};

    use crate::embedding::drivers::EmbeddingInput;

    #[derive(Serialize, Deserialize)]
    pub struct OpenAIEmbeddingRequest {
        pub input: EmbeddingInput,
        pub model: String,
    }

    #[derive(Serialize, Deserialize)]
    pub struct OpenAIEmbeddingResponse {
        pub model: Option<String>,
        pub data: Vec<Embedding>,
        pub usage: Usage,
        pub object: OpenAIEmbeddingResponseObject,
    }

    #[derive(Serialize, Deserialize)]
    pub enum OpenAIEmbeddingResponseObject {
        #[serde(rename = "list")]
        List,
    }

    #[derive(Serialize, Deserialize)]
    pub struct Embedding {
        pub index: u32,
        pub embedding: Vec<f64>,
        pub object: EmbeddingObject,
    }

    #[derive(Serialize, Deserialize)]
    pub enum EmbeddingObject {
        #[serde(rename = "embedding")]
        Embedding,
    }

    #[derive(Serialize, Deserialize)]
    pub struct Usage {
        pub prompt_tokens: u64,
        pub total_tokens: u64,
    }
}
