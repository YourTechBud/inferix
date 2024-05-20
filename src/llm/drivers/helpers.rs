use crate::{
    http::{AppError, StandardErrorResponse},
    llm::models::ModelConfig,
};

use super::{Driver, DRIVERS};


pub fn get_driver(model: &ModelConfig) -> Result<&Driver, AppError> {
    let drivers =
        DRIVERS
            .get()
            .ok_or(AppError::InternalServerError(StandardErrorResponse::new(
                "Drivers not initialized".to_string(),
                "drivers_not_initialized".to_string(),
            )))?;

    // TODO: Check if the model is available in the drivers list during config load
    return Ok(drivers.get(&model.driver).ok_or(AppError::BadRequest(
        StandardErrorResponse::new(
            "Driver not found".to_string(),
            "driver_not_found".to_string(),
        ),
    ))?);
}
