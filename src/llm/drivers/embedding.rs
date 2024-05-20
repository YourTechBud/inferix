use crate::{
    http::AppError,
    llm::types::{EmbeddingRequest, EmbeddingResponse},
};

use super::*;

pub async fn create_embeddings(mut req: EmbeddingRequest) -> Result<EmbeddingResponse, AppError> {
    // Get the model
    let model = crate::llm::models::get_model(&req.model)?;

    // Override the model name in request
    req.model = model.get_model_name().to_string();

    // Get the driver from the drivers list
    let driver = helpers::get_driver(&model)?;

    return driver.create_embedding(&req).await;
}
