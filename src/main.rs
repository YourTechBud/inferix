use std::collections::HashSet;

use axum::{
    extract::Request,
    middleware::{self, Next},
    response::Response,
    Router,
};
use clap::Parser;
use once_cell::sync::OnceCell;
use tower::ServiceBuilder;

#[derive(Debug, Parser)]
#[command(version, about, long_about = None)]
struct Cli {
    #[arg(long)]
    jwt: bool,

    #[arg(long, env)]
    jwt_algo: Option<String>, // TODO: Implement this

    #[arg(long, env)]
    jwt_secret: Option<String>,
}

static CLI_ARGS: OnceCell<Cli> = OnceCell::new();

#[tokio::main]
async fn main() {
    let cli = Cli::parse();

    // TODO: Need to clean this up
    CLI_ARGS.set(cli).unwrap();

    // Initialize the LLM driver
    // TODO: Load which drivers are required from the config file
    inferix::llm::init();

    // Setup the cors middleware
    // TODO: Allow users to modify the cors settings
    let cors = tower_http::cors::CorsLayer::new()
        .allow_methods(tower_http::cors::Any)
        .allow_origin(tower_http::cors::Any)
        .allow_headers(tower_http::cors::Any);

    // Start by creating a router
    let app = Router::new()
        .nest("/api/llm/v1", inferix::llm::routes::new())
        .layer(
            ServiceBuilder::new()
                .layer(cors)
                .layer(middleware::from_fn(jwt_auth_middleware)),
        );

    // Start the server
    let listener = tokio::net::TcpListener::bind("0.0.0.0:4386").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

// TODO: Move this to a middleware module
async fn jwt_auth_middleware(
    req: Request,
    next: Next,
) -> Result<Response, inferix::http::AppError> {
    let cli = CLI_ARGS.get().unwrap();

    if cli.jwt {
        // Check if the request has a valid JWT token
        let token = req
            .headers()
            .get("authorization")
            .ok_or(inferix::http::AppError::Unauthenticated(
                "No token provided in request".to_string(),
            ))?;
        let token = token
            .to_str()
            .map_err(|_| inferix::http::AppError::Unauthenticated("Invalid token provided".to_string()))?;

        // Strip the bearer prefix if it exists
        let token = token.strip_prefix("Bearer ").unwrap_or(token);

        // Create a validation object
        let mut validation = jsonwebtoken::Validation::new(jsonwebtoken::Algorithm::HS256);
        validation.required_spec_claims = HashSet::new();

        // Verify the token using HS256
        let secret = cli.jwt_secret.as_ref().unwrap();
        let _ = jsonwebtoken::decode::<serde_json::Value>(
            &token,
            &jsonwebtoken::DecodingKey::from_secret(secret.as_bytes()),
            &validation,
        )
        .map_err(|e| {
            return inferix::http::AppError::Unauthenticated(format!("Invalid token provided: {}", e));
        })?;
    }

    return Ok(next.run(req).await);
}
