use dotenvy::dotenv;
use serde::Deserialize;
use std::env;

#[derive(Debug, Deserialize)]
struct Body {
    origin: String,
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    println!("Hello, world!");
    dotenv().expect(".env file not found");
    println!("{:?}", env::var("APP_NAME"));
    let resp = reqwest::get("https://simesaba80.net/").await?;
    let body = resp.json::<Body>().await?;
    println!("{body:#?}");
    Ok(())
}
