use axum::Router;

#[tokio::main]
async fn main() {
    // Initialize the LLM driver
    // TODO: Load which drivers are required from the config file
    inferix::llm::init();

    // Start by creating a router
    let app = Router::new().nest("/api/llm/v1", inferix::llm::routes::new());

    // Start the server
    let listener = tokio::net::TcpListener::bind("0.0.0.0:4386").await.unwrap();
    axum::serve(listener, app).await.unwrap();
}
