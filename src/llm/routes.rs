use axum::{routing::post, Router};

use super::apis::{
    openai::{handle_chat_completion, handle_embedding_request},
    tei::handle_embed_request,
};

pub fn new() -> Router {
    return axum::Router::new()
        // Routes for OpenAI
        .route("/chat/completions", post(handle_chat_completion))
        .route("/embeddings", post(handle_embedding_request))

        // Routes for Text Embedding Inference
        .route("/embed", post(handle_embed_request));
}
