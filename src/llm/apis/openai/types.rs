use serde::{Deserialize, Serialize};
use serde_with::skip_serializing_none;

use crate::llm::types::{EmbeddingInput, FinishReason};

pub enum ToolSelection {
    None,
    Tools,
    Function,
}

/******************************************************************/
/***********************  Types for Request  **********************/
/******************************************************************/
#[skip_serializing_none]
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreateChatCompletionRequest {
    pub messages: Vec<ChatCompletionRequestMessage>,
    pub model: String,
    pub frequency_penalty: Option<f64>,
    // pub logprobs: Option<bool>,
    // pub top_logprobs: Option<i32>,
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

#[skip_serializing_none]
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

#[skip_serializing_none]
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

#[skip_serializing_none]
#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct CreateChatCompletionResponse {
    pub id: String,
    pub choices: Vec<ResponseChoice>,
    pub created: i64,
    pub model: String,
    pub system_fingerprint: Option<String>,
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

#[skip_serializing_none]
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
pub struct CompletionUsage {
    pub completion_tokens: u64,
    pub prompt_tokens: u64,
    pub total_tokens: u64,
}

/******************************************************************/
/*****************  Types for Streaming response  *****************/
/******************************************************************/

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug)]
pub struct CreateChatCompletionStreamResponse {
    pub id: String,
    pub choices: Vec<StreamingResponseChoice>,
    pub created: i64,
    pub model: String,
    pub system_fingerprint: Option<String>,
    pub object: String,
    pub usage: Option<Usage>,
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug)]
pub struct StreamingResponseChoice {
    pub delta: ChatCompletionStreamResponseDelta,
    pub finish_reason: Option<FinishReason>,
    pub index: i32,
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug)]
pub struct ChatCompletionStreamResponseDelta {
    pub content: Option<String>,
    pub function_call: Option<FunctionCall>,
    pub tool_calls: Option<Vec<ChatCompletionMessageToolCallChunk>>,
    pub role: Option<String>,
}

#[skip_serializing_none]
#[derive(Serialize, Deserialize, Debug)]
pub struct ChatCompletionMessageToolCallChunk {
    pub index: i32,
    pub id: Option<String>,
    #[serde(rename = "type")]
    pub type_field: Option<String>,
    pub function: Option<FunctionCall>,
}

/******************************************************************/
/**********************  Types for Embedding  *********************/
/******************************************************************/

#[derive(Serialize, Deserialize)]
pub struct OpenAIEmbeddingRequest {
    pub input: EmbeddingInput,
    pub model: String,
}

#[derive(Serialize, Deserialize)]
pub struct OpenAIEmbeddingResponse {
    pub model: Option<String>,
    pub data: Vec<Embedding>,
    pub usage: EmbeddingUsage,
    pub object: OpenAIEmbeddingResponseObject,
}

#[derive(Serialize, Deserialize)]
pub enum OpenAIEmbeddingResponseObject {
    #[serde(rename = "list")]
    List,
}

#[derive(Serialize, Deserialize)]
pub struct Embedding {
    pub index: u32,
    pub embedding: Vec<f64>,
    pub object: EmbeddingObject,
}

#[derive(Serialize, Deserialize)]
pub enum EmbeddingObject {
    #[serde(rename = "embedding")]
    Embedding,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct Usage {
    pub prompt_tokens: u64,
    pub completion_tokens: u64,
    pub total_tokens: u64,
}

#[derive(Serialize, Deserialize, Debug)]
pub struct EmbeddingUsage {
    pub prompt_tokens: u64,
    pub total_tokens: u64,
}

