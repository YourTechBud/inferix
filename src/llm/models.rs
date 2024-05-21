use std::collections::HashMap;

use once_cell::sync::OnceCell;
use serde::Deserialize;

use crate::http::{AppError, StandardErrorResponse};

#[derive(Debug, Deserialize)]
pub struct ModelConfig {
    name: String,
    pub driver: String,
    pub target_name: Option<String>,
    pub default_options: Option<ModelOptions>,
    pub prompt_tmpl: Option<String>,
}

impl ModelConfig {
    pub fn get_name(&self) -> &str {
        return &self.name;
    }
    pub fn get_target_name(&self) -> &str {
        return match &self.target_name {
            Some(name) => name,
            None => &self.name,
        }
    }
}

#[derive(Debug, Deserialize)]
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
            num_ctx: None,
            temperature: Some(0.2),
            driver_options: None,
        };
    }
}

pub static MODELS: OnceCell<HashMap<String, ModelConfig>> = OnceCell::new();

pub fn init(models: Vec<ModelConfig>) {
    let mut m = HashMap::new();
    for mut model in models {
        // Set the default options if not set
        // TODO: Set only those fields which are not set by the user
        if model.default_options.is_none() {
            model.default_options = Some(ModelOptions::default());
        }

        m.insert(model.name.clone(), model);
    }

    MODELS.set(m).unwrap();
}

pub fn get_model<'a, 'b>(model: &'a str) -> Result<&'b ModelConfig, AppError> {
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

    return Ok(&model);
}
