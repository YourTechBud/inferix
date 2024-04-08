use std::{collections::HashMap, sync::Arc};

use once_cell::sync::OnceCell;

use crate::http::{AppError, StandardErrorResponse};

#[derive(Debug)]
pub struct Model {
    pub driver: String,
    pub default_options: Option<ModelOptions>,
    pub prompt_tmpl: Option<String>,
}

#[derive(Debug)]
pub struct ModelOptions {
    pub top_p: Option<f64>,
    pub top_k: Option<i32>,
    pub num_ctx: Option<i32>,
    pub temperature: Option<f64>,
    pub driver_options: Option<serde_json::Value>,
}

impl ModelOptions {
    pub fn default() -> Self {
        return ModelOptions {
            top_p: Some(0.9),
            top_k: Some(40),
            num_ctx: Some(4096),
            temperature: Some(0.2),
            driver_options: None,
        };
    }
}

pub static MODELS: OnceCell<HashMap<String, Arc<Model>>> = OnceCell::new();

pub fn init() {
    let mut m = HashMap::new();

    m.insert(
        "mistral-openhermes".to_string(),
        Arc::new(Model {
            driver: "ollama".to_string(),
            default_options: Some(ModelOptions::default()),
            prompt_tmpl: None,
        }),
    );

    MODELS.set(m).unwrap();
}

pub fn get_model(model: &str) -> Result<Arc<Model>, AppError> {
    let models = MODELS
        .get()
        .ok_or(AppError::InternalServerError(StandardErrorResponse::new(
            "Models not initialized".to_string(),
            "models_not_initialized".to_string(),
        )))?;

    let model = models
        .get(model)
        .ok_or(AppError::BadRequest(StandardErrorResponse::new(
            format!("Model {} not found", model),
            "model_not_found".to_string(),
        )))?;

    return Ok(model.clone());
}
