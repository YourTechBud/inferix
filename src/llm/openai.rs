use axum::Json;

use crate::utils;

use super::{
    drivers,
    types::{AppError, StandardErrorResponse},
};
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

    // Select the tools provided from the `tools` or `functions` field in the request
    // We will give priority to the `tools` field since the `functions` field is deprecated.
    // TODO: Make this more efficient and look good
    let mut tool_selection = types::ToolSelection::None;

    let mut tools: Option<Vec<drivers::Tool>> = if let Some(tools) = req.tools {
        tool_selection = types::ToolSelection::Tools;
        Some(
            tools
                .iter()
                .map(|t| crate::llm::drivers::Tool {
                    name: t.function.name.clone(),
                    description: t.function.description.clone(),
                    args: t
                        .function
                        .parameters
                        .clone()
                        .unwrap_or(FunctionParameters {
                            additional_properties: serde_json::json!({}),
                        })
                        .additional_properties,
                    tool_type: crate::llm::drivers::ToolType::Function,
                })
                .collect(),
        )
    } else {
        None
    };

    // Only checks if any functions are provided if tools is still None
    if tools.is_none() {
        if let Some(functions) = req.functions {
            tool_selection = types::ToolSelection::Function;
            tools = Some(
                functions
                    .iter()
                    .map(|f| crate::llm::drivers::Tool {
                        name: f.name.clone(),
                        description: f.description.clone(),
                        args: f.parameters.additional_properties.clone(),
                        tool_type: crate::llm::drivers::ToolType::Function,
                    })
                    .collect(),
            )
        }
    }

    // Run the inference
    let res = crate::llm::drivers::run_inference(
        crate::llm::drivers::Request::new(req.model, messages, tools),
        crate::llm::drivers::RequestOptions::new(req.top_p, None, req.max_tokens, req.temperature),
    )
    .await?;

    // Prepare and return the response
    let res = CreateChatCompletionResponse {
        id: "inferix".to_string(),
        created: utils::convert_to_unix_timestamp(&res.created_at),
        model: res.model.clone(),
        object: "chat.completion".to_string(),
        system_fingerprint: "inferix".to_string(),
        usage: CompletionUsage {
            completion_tokens: res.stats.eval_count.unwrap_or(0),
            prompt_tokens: res.stats.prompt_eval_count.unwrap_or(0),
            total_tokens: res.stats.eval_count.unwrap_or(0)
                + res.stats.prompt_eval_count.unwrap_or(0),
        },
        choices: vec![prepare_chat_completion_message(&res, tool_selection)],
    };

    return Ok(Json(res));
}

fn prepare_chat_completion_message(
    response: &drivers::Response,
    tool_selection: types::ToolSelection,
) -> ResponseChoice {
    // Prepare the function call and tools varibles
    let mut fn_call = None;
    let mut tools = None;

    match tool_selection {
        types::ToolSelection::Function => {
            if let Some(fc) = &response.fn_call {
                fn_call = Some(types::FunctionCall {
                    name: fc.name.clone(),
                    arguments: serde_json::to_string(&fc.parameters).unwrap(),
                });
            }
        }

        types::ToolSelection::Tools => {
            if let Some(tc) = &response.fn_call {
                tools = Some(vec![ChatCompletionMessageToolCall {
                    id: tc.name.clone(),
                    tool_type: types::ToolType::Function,
                    function: FunctionCall {
                        name: tc.name.clone(),
                        arguments: serde_json::to_string(&tc.parameters).unwrap(),
                    },
                }])
            }
        }

        _ => {}
    }

    let finish_reason = if fn_call.is_some() {
        FinishReason::FunctionCall
    } else {
        FinishReason::Stop
    };

    return ResponseChoice {
        message: ChatCompletionResponseMessage {
            content: Some(response.response.to_string()),
            role: "assistant".to_string(),
            function_call: fn_call,
            tool_calls: tools,
        },
        index: 0,
        finish_reason: finish_reason,
    };
}

pub mod types {
    use serde::{Deserialize, Serialize};

    pub enum ToolSelection {
        None,
        Tools,
        Function,
    }

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
        pub id: String,
        #[serde(rename = "type")]
        pub tool_type: ToolType,
        pub function: FunctionCall,
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
        pub tool_calls: Option<Vec<ChatCompletionMessageToolCall>>,
        pub function_call: Option<FunctionCall>,
    }

    #[derive(Debug, Clone, Serialize, Deserialize)]
    pub enum ToolType {
        #[serde(rename = "function")]
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
        #[serde(rename = "type")]
        pub tool_type: ToolType,
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
        pub tool_calls: Option<Vec<ChatCompletionMessageToolCall>>,
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
