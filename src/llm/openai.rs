use axum::Json;

use crate::utils;

use super::types::{AppError, StandardErrorResponse};
use types::*;

pub async fn handle_chat_completion(
    Json(req): Json<CreateChatCompletionRequest>,
) -> Result<Json<CreateChatCompletionResponse>, AppError> {
    // Stream is not supported for now
    if req.stream.unwrap_or(false) {
        return Err(AppError::BadRequest(StandardErrorResponse::new(
            "Streaming responses is not supported".to_string(),
            "stream_not_supported".to_string(),
        )));
    }

    // Convert request messages to the format expected by the driver
    let messages = req
        .messages
        .iter()
        .map(|m| crate::llm::drivers::RequestMessage {
            role: m.role.clone(),
            content: m.content.clone(),
        })
        .collect();

    // Run the inference
    let res = crate::llm::drivers::run_inference(
        crate::llm::drivers::Request::new(req.model, messages),
        crate::llm::drivers::RequestOptions::new(req.top_p, None, req.max_tokens, req.temperature),
    )
    .await?;

    // Prepare and return the response
    let res = CreateChatCompletionResponse {
        id: "inferix".to_string(),
        created: utils::convert_to_unix_timestamp(&res.created_at),
        model: res.model,
        object: "chat.completion".to_string(),
        system_fingerprint: "inferix".to_string(),
        usage: CompletionUsage {
            completion_tokens: res.stats.eval_count.unwrap_or(0),
            prompt_tokens: res.stats.prompt_eval_count.unwrap_or(0),
            total_tokens: res.stats.eval_count.unwrap_or(0)
                + res.stats.prompt_eval_count.unwrap_or(0),
        },
        choices: vec![prepare_chat_completion_message(res.response.as_str())],
    };

    return Ok(Json(res));
}

fn prepare_chat_completion_message(response: &str) -> ResponseChoice {
    return ResponseChoice {
        message: ChatCompletionResponseMessage {
            content: Some(response.to_string()),
            role: "assistant".to_string(),
            function_call: None,
            tool_calls: None,
        },
        index: 0,
        finish_reason: FinishReason::Stop,
    };
}

pub mod types {
    use serde::{Deserialize, Serialize};

    /******************************************************************/
    /***********************  Types for Request  **********************/
    /******************************************************************/

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct CreateChatCompletionRequest {
        pub messages: Vec<ChatCompletionRequestMessage>,
        pub model: String,
        pub frequency_penalty: Option<f64>,
        pub logprobs: Option<bool>,
        pub top_logprobs: Option<i32>,
        pub max_tokens: Option<i32>,
        pub n: Option<i32>,
        pub presence_penalty: Option<f64>,
        pub response_format: Option<ResponseFormat>,
        pub seed: Option<i64>,
        pub stop: Option<Stop>,
        pub stream: Option<bool>,
        pub temperature: Option<f64>,
        pub top_p: Option<f64>,
        pub tools: Option<Vec<ChatCompletionTool>>,
        pub tool_choice: Option<ChatCompletionToolChoiceOption>,
        pub function_call: Option<FunctionCallRequest>,
        pub functions: Option<Vec<ChatCompletionFunctions>>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct FunctionParameters {
        #[serde(flatten)]
        pub additional_properties: serde_json::Value,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionRequestSystemMessage {
        pub content: String,
        pub role: String,
        pub name: Option<String>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionRequestUserMessage {
        pub content: String,
        pub role: String,
        pub name: Option<String>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct FunctionCall {
        pub name: String,
        pub arguments: String,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionMessageToolCall {
        id: String,
        #[serde(rename = "type")]
        type_: String,
        function: FunctionCall,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionMessageToolCalls {
        pub items: Vec<ChatCompletionMessageToolCall>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionRequestAssistantMessage {
        pub content: Option<String>,
        pub role: String,
        pub name: Option<String>,
        pub tool_calls: Option<ChatCompletionMessageToolCalls>,
        pub function_call: Option<FunctionCall>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionRequestToolMessage {
        pub role: String,
        pub content: String,
        pub tool_call_id: String,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionRequestFunctionMessage {
        pub role: String,
        pub content: Option<String>,
        pub name: String,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionRequestMessage {
        pub content: Option<String>,
        pub role: String,
        pub name: Option<String>,
        pub tool_call_id: Option<String>,
        pub tool_calls: Option<ChatCompletionMessageToolCalls>,
        pub function_call: Option<FunctionCall>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub enum ToolType {
        Function,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionFunctionCallOption {
        pub name: String,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionFunctions {
        pub description: Option<String>,
        pub name: String,
        pub parameters: FunctionParameters,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct FunctionObject {
        pub description: Option<String>,
        pub name: String,
        pub parameters: Option<FunctionParameters>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionTool {
        pub r#type: ToolType,
        pub function: FunctionObject,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionNamedToolChoice {
        pub r#type: ToolType,
        pub function: FunctionObject,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    #[serde(untagged)]
    pub enum ChatCompletionToolChoiceOption {
        None,
        Auto,
        NamedToolChoice(ChatCompletionNamedToolChoice),
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub enum ResponseFormat {
        #[serde(rename = "text")]
        Text,
        #[serde(rename = "json_object")]
        JsonObject,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    #[serde(untagged)]
    pub enum Stop {
        String(String),
        Array(Vec<String>),
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    #[serde(untagged)]
    pub enum FunctionCallRequest {
        String(String),
        ChatCompletionFunctionCallOption(ChatCompletionFunctionCallOption),
    }

    /******************************************************************/
    /**********************  Types for Response  **********************/
    /******************************************************************/

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct CreateChatCompletionResponse {
        pub id: String,
        pub choices: Vec<ResponseChoice>,
        pub created: i64,
        pub model: String,
        pub system_fingerprint: String,
        pub object: String,
        // Placeholder for the unimplemented type
        pub usage: CompletionUsage,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ResponseChoice {
        pub finish_reason: FinishReason,
        pub index: i32,
        pub message: ChatCompletionResponseMessage,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct ChatCompletionResponseMessage {
        #[serde(skip_serializing_if = "Option::is_none")]
        pub content: Option<String>,
        #[serde(skip_serializing_if = "Option::is_none")]
        pub tool_calls: Option<ChatCompletionMessageToolCalls>,
        pub role: String,
        #[serde(skip_serializing_if = "Option::is_none")]
        pub function_call: Option<FunctionCall>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub enum FinishReason {
        #[serde(rename = "stop")]
        Stop,
        #[serde(rename = "length")]
        Length,
        #[serde(rename = "tool_calls")]
        ToolCalls,
        #[serde(rename = "content_filter")]
        ContentFilter,
        #[serde(rename = "function_call")]
        FunctionCall,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub struct CompletionUsage {
        pub completion_tokens: u64,
        pub prompt_tokens: u64,
        pub total_tokens: u64,
    }
}
