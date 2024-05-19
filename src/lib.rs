pub mod http;
pub mod llm;
pub mod utils;

mod config;

pub fn init(config_path: Option<std::path::PathBuf>) {
    // TODO: Handle the error
    let c = config::Config::load_config(config_path).unwrap();
    llm::init(c);
}
