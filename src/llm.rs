pub mod routes;
pub mod types;

mod drivers;
mod models;
mod openai;
mod prompts;

pub fn init() {
    drivers::init();
    models::init();
    prompts::init();
}
