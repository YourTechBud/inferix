use std::{collections::HashSet, path::PathBuf};

use axum::{
    extract::Request,
    middleware::{self, Next},
    response::Response,
    Router,
};
use axum_server::tls_rustls::RustlsConfig;
use clap::Parser;
use once_cell::sync::OnceCell;
use serde::Serialize;
use tower::ServiceBuilder;

// TODO: Move this to a separate module
#[derive(clap::ValueEnum, Clone, Default, Debug, Serialize)]
#[serde(rename_all = "kebab-case")]
enum LogLevel {
    Debug,
    #[default]
    Info,
    Error,
}

// TODO: Move this to a separate module
#[derive(Debug, Parser)]
#[command(version, about, long_about = None)]
struct Cli {
    /************************************************/
    /********************* Config *******************/
    /************************************************/
    #[arg(short, long, env, value_name = "FILE")]
    config: Option<PathBuf>,

    /************************************************/
    /********************* Loging *******************/
    /************************************************/
    #[arg(long, env, default_value_t, value_enum)]
    log_level: LogLevel,

    /************************************************/
    /****************** HTTP Settings ***************/
    /************************************************/
    #[arg(long, env, default_value = "4386")]
    port: u16,

    /************************************************/
    /****************** TLS Settings ****************/
    /************************************************/
    #[arg(long)]
    tls: bool,

    #[arg(long)]
    tls_port: Option<u16>,

    #[arg(long, value_name = "FILE")]
    tls_cert: Option<PathBuf>,

    #[arg(long, value_name = "FILE")]
    tls_key: Option<PathBuf>,

    /************************************************/
    /*************** JWT Authentication *************/
    /************************************************/
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

    // TODO: Make this cleaner
    let cli = CLI_ARGS.get().unwrap();

    // Initialise the Logger
    tracing_subscriber::fmt()
        .with_max_level(match cli.log_level {
            LogLevel::Debug => tracing::Level::DEBUG,
            LogLevel::Info => tracing::Level::INFO,
            LogLevel::Error => tracing::Level::ERROR,
        })
        .with_target(false)
        .init();

    // Initialize the LLM driver
    // TODO: Load which drivers are required from the config file
    inferix::init(cli.config.clone());

    // Setup the cors middleware
    // TODO: Allow users to modify the cors settings
    let cors = tower_http::cors::CorsLayer::new()
        .allow_methods(tower_http::cors::Any)
        .allow_origin(tower_http::cors::Any)
        .allow_headers(tower_http::cors::Any);

    // Start by creating a router
    let app = Router::new()
        .nest("/api/v1", inferix::llm::routes::new())
        .layer(
            ServiceBuilder::new()
                .layer(cors)
                .layer(middleware::from_fn(jwt_auth_middleware)),
        );

    if cli.tls {
        tokio::spawn(start_tls_sever(cli, app.clone()));
    }

    // Start the server
    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], cli.port));
    tracing::info!("Starting server on {}", addr.port());
    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();
    axum::serve(listener, app).await.unwrap();
}

async fn start_tls_sever(cli: &Cli, app: Router) {
    // Configure certificate and private key used by https
    let config =
        RustlsConfig::from_pem_file(cli.tls_cert.clone().unwrap(), cli.tls_key.clone().unwrap())
            .await
            .unwrap();

    // Configure the port
    let port = cli.tls_port.unwrap_or(4388);
    tracing::info!("Starting server on {}", port);

    // Start the server
    let addr = std::net::SocketAddr::from(([0, 0, 0, 0], port));
    axum_server::bind_rustls(addr, config)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

// TODO: Move this to a middleware module
async fn jwt_auth_middleware(
    req: Request,
    next: Next,
) -> Result<Response, inferix::http::AppError> {
    let cli = CLI_ARGS.get().unwrap();

    if cli.jwt {
        // Check if the request has a valid JWT token
        let token =
            req.headers()
                .get("authorization")
                .ok_or(inferix::http::AppError::Unauthenticated(
                    "No token provided in request".to_string(),
                ))?;
        let token = token.to_str().map_err(|_| {
            inferix::http::AppError::Unauthenticated("Invalid token provided".to_string())
        })?;

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
            return inferix::http::AppError::Unauthenticated(format!(
                "Invalid token provided: {}",
                e
            ));
        })?;
    }

    return Ok(next.run(req).await);
}
