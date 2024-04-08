use async_trait::async_trait;
use reqwest::Client;
use serde::{Deserialize, Serialize};
use std::option::Option;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::prompts,
    utils,
};

use super::*;

#[derive(Debug)]
pub struct Ollama {
    host: String,
    port: String,
}

impl Ollama {
    pub fn new(host: String, port: String) -> Self {
        return Ollama { host, port };
    }
}

#[derive(Serialize, Deserialize)]
pub struct OllamaRequest<'a> {
    pub model: &'a str,
    pub prompt: &'a str,
    pub raw: bool,
    pub stream: Option<bool>,
    pub options: serde_json::Value,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct OllamaResponse {
    pub model: String,
    pub created_at: String,
    pub response: String,
    pub done: bool,
    pub total_duration: Option<u64>,
    pub load_duration: Option<u64>,
    pub prompt_eval_count: Option<u64>,
    pub prompt_eval_duration: Option<u64>,
    pub eval_count: Option<u64>,
    pub eval_duration: Option<u64>,
}

#[async_trait]
impl Driver for Ollama {
    async fn call(&self, req: &Request, options: &RequestOptions) -> Result<Response, AppError> {
        // Let's first get the prompt
        let prompt = prompts::get_prompt(&options.prompt_tmpl, &req.messages)?;

        // Prepare the request
        let req = OllamaRequest {
            model: &req.model,
            prompt: &prompt,
            raw: true,
            stream: Some(false),
            options: utils::merge_objects(
                &serde_json::json!({
                    "top_p": options.top_p,
                    "top_k": options.top_k,
                    "num_ctx": options.num_ctx,
                    "temperature": options.temperature,
                }),
                &options.driver_options,
            ),
        };

        println!("====================================");
        println!("Request Prompt:");
        println!("{}", prompt);
        println!("====================================");

        // Call the Ollama API
        let client = Client::new();
        let url = format!("http://{}:{}/api/generate", self.host, self.port);
        let response = client.post(&url).json(&req).send().await.map_err(|e| {
            eprintln!("Request error: {}", e);
            return AppError::BadRequest(StandardErrorResponse::new(
                "Unable to make request to Ollama".to_string(),
                e.to_string(),
            ));
        })?;

        let res = response.text().await.map_err(|e| {
            eprintln!("Response error: {}", e);
            return AppError::BadRequest(StandardErrorResponse::new(
                "Unable to get response from Ollama".to_string(),
                e.to_string(),
            ));
        })?;

        let mut res: OllamaResponse = serde_json::from_str(&res).map_err(|e| {
            eprintln!("Response parsing error: {}", e);
            return AppError::BadRequest(StandardErrorResponse::new(
                "Unable to parse Ollama response".to_string(),
                e.to_string(),
            ));
        })?;

        // Remove the weird ':' and '>' prefix that mistral keeps giving us.
        res.response = res.response.trim().to_string();
        if res.response.starts_with(':') || res.response.starts_with('>') {
            res.response = res.response[1..].to_string();
        }

        println!("Response:");
        println!("{}", res.response);
        println!("====================================");

        return Ok(Response {
            model: res.model,
            created_at: res.created_at,
            response: res.response,
            stats: ResponseStats {
                total_duration: res.total_duration,
                load_duration: res.load_duration,
                prompt_eval_count: res.prompt_eval_count,
                prompt_eval_duration: res.prompt_eval_duration,
                eval_count: res.eval_count,
                eval_duration: res.eval_duration,
            },
            fn_call: None,
        });
    }
}

#[cfg(test)]
mod test {

    use super::*;

    #[tokio::test]
    #[ignore]
    async fn test_ollama() {
        prompts::init();

        let driver = Ollama {
            host: "localhost".to_string(),
            port: "11434".to_string(),
        };

        let req = Request {
            model: "mistral-openhermes".to_string(),
            messages: vec![
                RequestMessage {
                    role: "system".to_string(),
                    content: Some("You are a helpful AI assistant".to_string()),
                },
                RequestMessage {
                    role: "user".to_string(),
                    content: Some("What's the capital of British Columbia".to_string()),
                },
            ],
            tools: None,
        };

        let res = driver.call(&req, &RequestOptions::default()).await.unwrap();
        assert_eq!(res.model, "mistral-openhermes");
    }
}
