use async_trait::async_trait;
use serde::Serialize;
use std::fmt::Debug;

use crate::{http::AppError, llm::prompts};

#[async_trait]
pub trait Driver: Send + Sync + Debug {
    async fn call(&self, req: &Request, options: &RequestOptions) -> Result<Response, AppError>;
}

#[derive(Serialize)]

pub struct Request {
    pub model: String,
    pub messages: Vec<RequestMessage>,
    pub tools: Option<Vec<Tool>>,
}

impl Request {
    pub fn new(model: String, messages: Vec<RequestMessage>, tools: Option<Vec<Tool>>) -> Self {
        return Request {
            model: model,
            messages: messages,
            tools: tools,
        };
    }
}

#[derive(Serialize)]
pub struct RequestMessage {
    pub role: String,
    pub content: Option<String>,
}

#[derive(Serialize)]
pub struct Tool {
    pub name: String,
    pub description: Option<String>,
    pub args: serde_json::Value,
    #[serde(rename = "type")]
    pub tool_type: ToolType,
}

#[derive(Serialize)]
pub enum ToolType {
    #[serde(rename = "function")]
    Function,
}

pub struct RequestOptions {
    pub top_p: Option<f64>,
    pub top_k: Option<i32>,
    pub num_ctx: Option<i32>,
    pub temperature: Option<f64>,

    // Internal options only available to the drivers to set
    pub driver_options: serde_json::Value,
    pub prompt_tmpl: String,
}

impl RequestOptions {
    pub fn new(
        top_p: Option<f64>,
        top_k: Option<i32>,
        num_ctx: Option<i32>,
        temperature: Option<f64>,
    ) -> Self {
        return RequestOptions {
            top_p: top_p,
            top_k: top_k,
            num_ctx: num_ctx,
            temperature: temperature,

            driver_options: serde_json::json!({}),
            prompt_tmpl: prompts::CHATML.to_string(),
        };
    }

    pub fn default() -> Self {
        return RequestOptions {
            top_p: Some(0.9),
            top_k: Some(40),
            num_ctx: Some(4096),
            temperature: Some(0.2),

            driver_options: serde_json::json!({}),
            prompt_tmpl: prompts::CHATML.to_string(),
        };
    }
}

pub struct Response {
    pub model: String,
    pub created_at: String,
    pub response: String,
    pub stats: ResponseStats,
    pub fn_call: Option<FunctionCall>,
}

#[derive(serde::Deserialize)]
pub struct FunctionCall {
    pub name: String,
    pub parameters: serde_json::Value,
}

pub struct ResponseStats {
    pub total_duration: Option<u64>,
    pub load_duration: Option<u64>,
    pub prompt_eval_count: Option<u64>,
    pub prompt_eval_duration: Option<u64>,
    pub eval_count: Option<u64>,
    pub eval_duration: Option<u64>,
}
