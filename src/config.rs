use std::{
    io::{Error, ErrorKind},
    path::PathBuf,
};

use serde::Deserialize;

use super::llm::{drivers, models, prompts};

#[derive(Deserialize)]
pub struct Config {
    pub drivers: Vec<drivers::DriverConfig>,
    pub prompt_tmpls: Vec<prompts::PrompTemplate>,
    pub models: Vec<models::ModelConfig>,
}

impl Config {
    pub fn load_config(path: Option<PathBuf>) -> Result<Self, Error> {
        match path {
            Some(p) => {
                // Load the config file from the path
                let data = std::fs::read_to_string(p)?;

                // Parse the config file
                let c = serde_yaml::from_str(&data)
                    .map_err(|e| Error::new(ErrorKind::InvalidData, e))?;

                // Simply return the config object
                return Ok(c);
            }
            None => {
                print!("No config file provided. Using default config.\n");
                return Ok(Config {
                    drivers: vec![],
                    prompt_tmpls: vec![],
                    models: vec![],
                });
            }
        }
    }
}
