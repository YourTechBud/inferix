use axum::{routing::post, Router};

use super::{openai::handle_embedding_request, tei::handle_embed_request};

pub fn new() -> Router {
    return axum::Router::new()
        .route("/embed", post(handle_embed_request))
        .route("/embeddings", post(handle_embedding_request));
}
