use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct StandardErrorResponse {
    message: String,

    #[serde(rename = "type")]
    error_type: Option<String>,
}

impl StandardErrorResponse {
    pub fn new(message: String, error_type: String) -> Self {
        Self {
            message,
            error_type: Some(error_type),
        }
    }
}

// The kinds of errors we can hit in our application.
#[derive(Debug, Serialize)]
pub enum AppError {
    StreamEndedd,
    BadRequest(StandardErrorResponse),
    InternalServerError(StandardErrorResponse),
    Unauthenticated(String),
}

// Tell axum how `AppError` should be converted into a response.
//
// This is also a convenient place to log errors.
impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            AppError::StreamEndedd => (
                StatusCode::OK,
                StandardErrorResponse::new("[DONE]".to_string(), "".to_string()),
            ),
            AppError::BadRequest(message) => (StatusCode::BAD_REQUEST, message),
            AppError::InternalServerError(message) => (StatusCode::INTERNAL_SERVER_ERROR, message),
            AppError::Unauthenticated(message) => (
                StatusCode::UNAUTHORIZED,
                StandardErrorResponse::new(message, "unauthenticated".to_string()),
            ),
        };

        if StatusCode::OK == status {
            return (status, message.message).into_response();
        }

        return (status, Json(message)).into_response();
    }
}

impl From<StandardErrorResponse> for AppError {
    fn from(bad_req: StandardErrorResponse) -> Self {
        Self::BadRequest(bad_req)
    }
}

impl std::fmt::Display for AppError {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        match self {
            AppError::BadRequest(message) => write!(f, "Bad Request: {}", message.message),
            AppError::InternalServerError(message) => {
                write!(f, "Internal Server Error: {}", message.message)
            }
            AppError::Unauthenticated(message) => write!(f, "Unauthenticated: {}", message),
            AppError::StreamEndedd => write!(f, "[DONE]"),
        }
    }
}

impl std::error::Error for AppError {}
