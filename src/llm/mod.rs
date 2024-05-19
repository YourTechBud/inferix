pub mod drivers;
pub mod prompts;
pub mod routes;
pub mod models;
pub mod apis;
pub mod types;

pub fn init(c: crate::config::Config) {
    prompts::init(c.prompt_tmpls);
    models::init(c.models);
    drivers::init(c.drivers);
}
