use crate::{
    http::AppError,
    llm::types::{EmbeddingRequest, EmbeddingResponse},
};

use super::*;

pub async fn create_embeddings(mut req: EmbeddingRequest) -> Result<EmbeddingResponse, AppError> {
    // Get the model
    let model = crate::llm::models::get_model(&req.model)?;

    // Override the model name in request
    req.model = model.get_target_name().to_string();

    // Get the driver from the drivers list
    let driver = helpers::get_driver(&model)?;

    // Make sure we use the right model name in the response.
    let mut result = driver.create_embedding(&req).await;
    if let Ok(res) = &mut result {
        res.model = Some(model.get_name().to_string());
    }
    return result;
}
