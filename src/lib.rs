pub mod embedding;
pub mod http;
pub mod llm;
pub mod utils;

pub fn init() {
    embedding::init();
    llm::init();
}
