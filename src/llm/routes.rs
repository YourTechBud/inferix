use axum::{routing::post, Router};

use super::openai::handle_chat_completion;

pub fn new() -> Router {
    return axum::Router::new().route("/chat/completions", post(handle_chat_completion));
}
