use async_trait::async_trait;
use reqwest::Client;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        drivers::InferenceDriver,
        prompts,
        types::{InferenceOptions, InferenceRequest, InferenceResponseSync, InferenceStats},
    },
    utils,
};

use super::{
    types::{OllamaRequest, OllamaResponse},
    Ollama,
};

#[async_trait]
impl InferenceDriver for Ollama {
    async fn run_inference(
        &self,
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> Result<InferenceResponseSync, AppError> {
        // TODO: We should default to using the chat completion endpoinnt. We can switch
        // to the generation endpoint only if the user explicitly specifies the prompt template/

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

        // Call the Ollama API
        let client = Client::new();
        let url = format!("{}/api/generate", self.base_url);
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

        return Ok(InferenceResponseSync {
            model: res.model,
            created_at: res.created_at,
            response: res.response,
            stats: InferenceStats {
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

    use crate::llm::types::InferenceMessage;

    use super::*;

    #[tokio::test]
    #[ignore]
    async fn test_ollama() {
        prompts::init(vec![
            crate::llm::prompts::PrompTemplate {
                name: "chatml".to_string(),
                tmpl: "{% for msg in messages %}<|im_start|>{{ msg.role }}\n{{ msg.content }}<|im_end|>\n{% endfor %}<|im_start|>assistant".to_string(),
                stop: vec!["<|im_start|>".to_string(), "<|im_end|>".to_string()],
            },
        ]);

        let driver = Ollama {
            base_url: "http://localhost:11434".to_string(),
        };

        let req = InferenceRequest {
            model: "mistral-openhermes".to_string(),
            messages: vec![
                InferenceMessage {
                    role: "system".to_string(),
                    content: Some("You are a helpful AI assistant".to_string()),
                },
                InferenceMessage {
                    role: "user".to_string(),
                    content: Some("What's the capital of British Columbia".to_string()),
                },
            ],
            tools: None,
        };

        let res = driver
            .run_inference(&req, &InferenceOptions::default())
            .await
            .unwrap();
        assert_eq!(res.model, "mistral-openhermes");
    }
}
