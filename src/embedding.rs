pub mod routes;

mod drivers;
mod tei;
mod openai;

pub fn init() {
    drivers::init();
}
