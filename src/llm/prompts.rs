use std::collections::HashMap;

use once_cell::sync::OnceCell;
use serde::Deserialize;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::types::InferenceMessage,
};

pub const CHATML: &str = "chatml";
pub const MISTRAL: &str = "mistral";

#[derive(Debug, Deserialize)]
pub struct PrompTemplate {
    pub name: String,
    pub tmpl: String,
    pub stop: Vec<String>,
}

impl PrompTemplate {
    pub fn new(name: String, tmpl: String, stop: Vec<String>) -> Self {
        return PrompTemplate {
            name: name,
            tmpl: tmpl,
            stop: stop,
        };
    }
}

pub static PROMPT_TEMPLATES: OnceCell<HashMap<String, PrompTemplate>> = OnceCell::new();

pub fn init(prompt_tmpls: Vec<PrompTemplate>) {
    let mut m = HashMap::new();
    for pt in prompt_tmpls {
        m.insert(pt.name.clone(), pt);
    }

    PROMPT_TEMPLATES.set(m).unwrap();
}

pub fn get_prompt(
    prompt_tmpl_name: &str,
    messages: &Vec<InferenceMessage>,
) -> Result<String, AppError> {
    let prompts_map =
        PROMPT_TEMPLATES
            .get()
            .ok_or(AppError::InternalServerError(StandardErrorResponse::new(
                "Prompt templates not initialized".to_string(),
                "prompt_templates_not_initialized".to_string(),
            )))?;

    let prompt_tmpl = prompts_map
        .get(prompt_tmpl_name)
        .ok_or(AppError::BadRequest(StandardErrorResponse::new(
            format!("Prompt template {} not found", prompt_tmpl_name),
            "prompt_template_not_found".to_string(),
        )))?;

    let mut context = tera::Context::new();
    context.insert("messages", &messages);

    // TODO: We should probably compile the template once and cache it
    let rendered = tera::Tera::one_off(&prompt_tmpl.tmpl, &context, false).map_err(|e| {
        AppError::InternalServerError(StandardErrorResponse::new(
            format!("Error rendering prompt template: {}", e.to_string()),
            "prompt_template_render_error".to_string(),
        ))
    })?;

    return Ok(rendered);
}
