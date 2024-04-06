mod ollama;
mod types;

use once_cell::sync::OnceCell;
pub use types::*;

use super::types::{AppError, StandardErrorResponse};

pub static DRIVER: OnceCell<Box<dyn Driver>> = OnceCell::new();

pub fn init() {
    DRIVER
        .set(Box::new(ollama::Ollama::new(
            "localhost".to_string(),
            "11434".to_string(),
        )))
        .unwrap();
}

pub async fn run_inference(req: Request, options: RequestOptions) -> Result<Response, AppError> {
    // Get the model
    let model = crate::llm::models::get_model(&req.model)?;

     // Populate the request with default options
     let mut options = options;
     if let Some(default_options) = &model.default_options {
         if options.top_p.is_none() {
             options.top_p = default_options.top_p;
         }
         if options.top_k.is_none() {
             options.top_k = default_options.top_k;
         }
         if options.num_ctx.is_none() {
             options.num_ctx = default_options.num_ctx;
         }
         if options.temperature.is_none() {
             options.temperature = default_options.temperature;
         }
 
         // Dont forget to load the driver options
         if let Some(driver_options) = &default_options.driver_options {
             options.driver_options = driver_options.clone();
         }
     }

    let driver = DRIVER
        .get()
        .ok_or(AppError::InternalServerError(StandardErrorResponse::new(
            "Driver not initialized".to_string(),
            "driver_not_initialized".to_string(),
        )))?;
    return driver.call(req, options).await;
}
