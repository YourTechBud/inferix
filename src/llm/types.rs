use axum::{
    http::StatusCode,
    response::{IntoResponse, Response},
    Json,
};
use serde::Serialize;

#[derive(Debug, Serialize)]
pub struct StandardErrorResponse {
    message: String,

    #[serde(rename = "type")]
    error_type: String,
}

impl StandardErrorResponse {
    pub fn new(message: String, error_type: String) -> Self {
        Self {
            message,
            error_type,
        }
    }
}

// The kinds of errors we can hit in our application.
#[derive(Debug)]
pub enum AppError {
    // The request body contained invalid request parameters
    BadRequest(StandardErrorResponse),

    InternalServerError(StandardErrorResponse),
}

// Tell axum how `AppError` should be converted into a response.
//
// This is also a convenient place to log errors.
impl IntoResponse for AppError {
    fn into_response(self) -> Response {
        let (status, message) = match self {
            AppError::BadRequest(message) => (StatusCode::BAD_REQUEST, message),
            AppError::InternalServerError(message) => (StatusCode::INTERNAL_SERVER_ERROR, message),
        };

        (status, Json(message)).into_response()
    }
}

impl From<StandardErrorResponse> for AppError {
    fn from(bad_req: StandardErrorResponse) -> Self {
        Self::BadRequest(bad_req)
    }
}
