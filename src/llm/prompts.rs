use std::collections::HashMap;

use once_cell::sync::OnceCell;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::drivers::RequestMessage,
};

pub const CHATML: &str = "chatml";
pub const MISTRAL: &str = "mistral";

#[derive(Debug)]
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

pub fn init() {
    let mut m = HashMap::new();

    // Let's add the chatml prompt template by default
    m.insert(
        "chatml".to_string(),
        PrompTemplate::new(
            "chatml".to_string(),
            "{% for msg in messages %}<|im_start|>{{ msg.role }}\n{{ msg.content }}<|im_end|>\n{% endfor %}<|im_start|>assistant".to_string(),
            vec!["<|im_start|>".to_string(), "<|im_end|>".to_string()],
        ),
    );

    m.insert(
        "mistral".to_string(),
        PrompTemplate::new(
            "mistral".to_string(),
            "[INST] {{ messages[0].content }}{% if messages | length == 2 %}\n\n{{messages[1].content}}{% endif %} [/INST]".to_string(),
            vec!["[INST]".to_string(), "[/INST]".to_string()],
        ),
    );

    PROMPT_TEMPLATES.set(m).unwrap();
}

pub fn get_prompt(
    prompt_tmpl_name: &str,
    messages: &Vec<RequestMessage>,
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

    let rendered = tera::Tera::one_off(&prompt_tmpl.tmpl, &context, false).map_err(|e| {
        AppError::InternalServerError(StandardErrorResponse::new(
            format!("Error rendering prompt template: {}", e.to_string()),
            "prompt_template_render_error".to_string(),
        ))
    })?;

    return Ok(rendered);
}
