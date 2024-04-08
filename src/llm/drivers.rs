mod ollama;
mod types;

use once_cell::sync::OnceCell;
pub use types::*;

use crate::utils;

use crate::http::{AppError, StandardErrorResponse};

pub static DRIVER: OnceCell<Box<dyn Driver>> = OnceCell::new();

pub fn init() {
    DRIVER
        .set(Box::new(ollama::Ollama::new(
            "localhost".to_string(),
            "11434".to_string(),
        )))
        .unwrap();
}

pub async fn run_inference(
    mut req: Request,
    mut options: RequestOptions,
) -> Result<Response, AppError> {
    // Get the model
    let model = crate::llm::models::get_model(&req.model)?;

    // Populate the request with default options
    if let Some(model_options) = &model.default_options {
        if options.top_p.is_none() {
            options.top_p = model_options.top_p;
        }
        if options.top_k.is_none() {
            options.top_k = model_options.top_k;
        }
        if options.num_ctx.is_none() {
            options.num_ctx = model_options.num_ctx;
        }
        if options.temperature.is_none() {
            options.temperature = model_options.temperature;
        }

        // Dont forget to load the driver options
        if let Some(driver_options) = &model_options.driver_options {
            options.driver_options = driver_options.clone();
        }
    }

    // Inject the function call message if tools are provided
    if let Some(tools) = &req.tools {
        req.messages = inject_fn_call(req.messages, &tools);
    }

    // Get the driver from the drivers list
    // TODO: Get the right driver based on the model
    let driver = DRIVER
        .get()
        .ok_or(AppError::InternalServerError(StandardErrorResponse::new(
            "Driver not initialized".to_string(),
            "driver_not_initialized".to_string(),
        )))?;

    // TODO: Restrict the loop to a certain number of iterations
    loop {
        // Call the driver and run inference
        let mut res = driver.call(&req, &options).await?;

        // Lets start by cleaning up the output
        res.response = res.response.trim().to_string();

        // Rerun inference if the response is too small
        if res.response.len() <= 2 {
            continue;
        }

        // Check if the response is a function call
        if res.response.contains("FUNC_CALL") {
            // Sanitize the json text
            res.response = utils::sanitize_json_text(&res.response);

            // Check if response is valid JSON
            let json_res = serde_json::from_str::<FunctionCall>(&res.response);
            if let Ok(fn_call) = json_res {
                res.response = "".to_string();
                res.fn_call = Some(fn_call);
            } else {
                // Run inference again if the response is not a valid JSON
                continue;
            }
        }

        return Ok(res);
    }
}

fn inject_fn_call(messages: Vec<RequestMessage>, functions: &Vec<Tool>) -> Vec<RequestMessage> {
    let mut content = "You may use the following FUNCTIONS in the response. Only use one function at a time. Give output in following OUTPUT_FORMAT in strict JSON if you want to call a function.\n\nFUNCTIONS:\n\n".to_string();

    for f in functions {
        content += &format!("- Name: {}\n", f.name);
        if let Some(desc) = &f.description {
            content += &format!("  Description: {}\n", desc);
        }
        content += &format!(
            "  Parameter JSON Schema: {}\n\n",
            serde_json::to_string(&f.args).unwrap()
        );
    }

    content += r#"OUTPUT_FORMAT:
{
    "type": "FUNC_CALL",
    "name": "<name of function>",
    "parameters": "<value to pass to function as parameter>"
}
"#;

    // Return a new messages array with a new system message injected towards the end
    return messages
        .into_iter()
        .chain(std::iter::once(RequestMessage {
            role: "system".to_string(),
            content: Some(content),
        }))
        .collect();
}
