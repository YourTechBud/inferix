use crate::{
    http::{AppError, StandardErrorResponse},
    llm::types::{EmbeddingRequest, EmbeddingResponse},
};

use super::DRIVERS;

pub async fn create_embeddings(mut req: EmbeddingRequest) -> Result<EmbeddingResponse, AppError> {
    // Get the model
    let model = crate::llm::models::get_model(&req.model)?;
    
    // Override the model name in request
    req.model = model.get_model_name().to_string();

    // Get the driver from the drivers list
    let drivers =
        DRIVERS
            .get()
            .ok_or(AppError::InternalServerError(StandardErrorResponse::new(
                "Drivers not initialized".to_string(),
                "drivers_not_initialized".to_string(),
            )))?;

    // TODO: Check if the model is available in the drivers list during config load
    let driver =
        drivers
            .get(&model.driver)
            .ok_or(AppError::BadRequest(StandardErrorResponse::new(
                "Driver not found".to_string(),
                "driver_not_found".to_string(),
            )))?;

    return driver.create_embedding(&req).await;
}
