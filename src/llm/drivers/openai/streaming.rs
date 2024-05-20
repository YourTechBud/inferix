use futures_core::stream::BoxStream;
use futures_util::StreamExt;
use reqwest_eventsource::EventSource;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        drivers::StreamingInference,
        types::{InferenceOptions, InferenceRequest, InferenceResponseStream},
    },
};

use super::OpenAI;

impl StreamingInference for OpenAI {
    fn run_streaming_inference(
        &self,
        req: &InferenceRequest,
        options: &InferenceOptions,
    ) -> Result<BoxStream<Result<InferenceResponseStream, AppError>>, AppError> {
        // Prepare openai request
        let mut req = OpenAI::prepare_request_message(req, options);
        req.stream = Some(true);

        // Fire the request
        let res = self.prepare_reqwest_builder(&req);
        let src = EventSource::new(res).map_err(|e| {
            return AppError::InternalServerError(StandardErrorResponse::new(
                format!("Unable to make request to OpenAI server: {}", e),
                "openai_call_error".to_string(),
            ));
        })?;

        let filtered_stream = OpenAI::prepare_event_stream(src);
        let output_stream = OpenAI::convert_to_inference_response_stream(filtered_stream);

        return Ok(output_stream.boxed());
    }
}
