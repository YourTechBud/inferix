use futures_util::StreamExt;

use crate::{
    http::{AppError, StandardErrorResponse},
    llm::{
        models::ModelConfig,
        types::{
            FunctionCall, InferenceMessage, InferenceOptions, InferenceRequest,
            InferenceResponseSync, Tool,
        },
    },
    utils,
};

use super::*;

pub async fn run_inference(
    mut req: InferenceRequest,
    mut options: InferenceOptions,
) -> Result<InferenceResponseSync, AppError> {
    // Get the model
    let model_config = crate::llm::models::get_model(&req.model)?;

    // Override the model name in request
    req.model = model_config.get_target_name().to_string();

    // Populate the request with default options
    populate_options(&mut options, model_config);

    // Inject the function call message if tools are provided
    if let Some(tools) = &req.tools {
        req.messages = inject_fn_call(req.messages, &tools);
    }

    // TODO: Check if the model is available in the drivers list during config load
    let driver = helpers::get_driver(&model_config)?;

    // Print the request and options
    tracing::debug!("---------------");
    tracing::debug!("Request: {:?}", req);
    tracing::debug!("Options: {:?}", options);
    tracing::debug!("---------------");

    // TODO: Restrict the loop to a certain number of iterations
    let mut i = 0;
    loop {
        // Call the driver and run inference
        let mut res = driver.run_inference(&req, &options).await?;

        // Lets start by cleaning up the output
        res.response = res.response.trim().to_string();

        tracing::debug!("---------------");
        tracing::debug!("Response: {}", res.response);
        tracing::debug!("---------------");

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
                res.response = format!(
                    "Execute function {} with arguments: {}",
                    fn_call.name,
                    serde_json::to_string(&fn_call.parameters).unwrap()
                )
                .to_string();
                res.fn_call = Some(fn_call);
            } else if i < 5 {
                // Run inference again if the response is not a valid JSON
                i = i + 1;
                continue;
            }
        }

        // Make sure we use the right model name in the response.
        res.model = model_config.get_name().to_string();
        return Ok(res);
    }
}

pub fn run_streaming_inference<'a>(
    mut req: InferenceRequest,
    mut options: InferenceOptions,
) -> Result<BoxStream<'a, Result<InferenceResponseStream, AppError>>, AppError> {
    // Function calling isn't allowed in streaming mode
    if req.tools.is_some() {
        return Err(AppError::BadRequest(StandardErrorResponse::new(
            "Function calling isn't allowed in streaming mode".to_string(),
            "function_calling_not_allowed".to_string(),
        )));
    }

    // Get the model
    let model_config = crate::llm::models::get_model(&req.model)?;

    // Override the model name in request
    req.model = model_config.get_target_name().to_string();

    // Populate the request with default options
    populate_options(&mut options, model_config);

    // Inject the function call message if tools are provided
    if let Some(tools) = &req.tools {
        req.messages = inject_fn_call(req.messages, &tools);
    }

    // TODO: Check if the model is available in the drivers list during config load
    let driver = helpers::get_driver(&model_config)?;

    // Call the driver and run inference
    let stream = driver.run_streaming_inference(&req, &options)?;

    // Make sure we inject the proper model name in the response
    let stream = stream
        .map(move |chunk| match chunk {
            Ok(mut c) => {
                c.model = model_config.get_name().to_string();
                return Ok(c);
            }
            _ => return chunk,
        })
        .boxed();
    return Ok(stream);
}

fn populate_options(options: &mut InferenceOptions, model_config: &ModelConfig) {
    if let Some(model_options) = &model_config.default_options {
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

        // Dont forget to load the driver options and the prompt template
        if let Some(driver_options) = &model_options.driver_options {
            options.driver_options = driver_options.clone();
        }
        if let Some(prompt_tmpl) = &model_config.prompt_tmpl {
            options.prompt_tmpl = prompt_tmpl.clone();
        }
    }
}

fn inject_fn_call(messages: Vec<InferenceMessage>, functions: &Vec<Tool>) -> Vec<InferenceMessage> {
    let mut content = "You may use the following FUNCTIONS in the response. Only use one function at a time. Give output in following OUTPUT_FORMAT if you want to call a function.\n\nFUNCTIONS:\n\n".to_string();

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
Parameter Selection:
<Provide the step by step thought process to select the parameters. Go through the entire conversation>

Function Call:
{
    "type": "FUNC_CALL",
    "reasoning": "<reasoning for choosing the parameters>",
    "name": "<name of function>",
    "parameters": "<value to pass to function as parameter>"
}
"#;

    // Return a new messages array with a new system message injected towards the end
    return messages
        .into_iter()
        .chain(std::iter::once(InferenceMessage {
            role: "system".to_string(),
            content: Some(content),
        }))
        .collect();
}
