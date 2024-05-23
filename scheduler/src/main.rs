use std::env;

use redis::Commands;
use tokio::time::{Duration as TokioDuration, sleep_until, Instant};
use chrono::{Duration, DateTime, Local, NaiveDate, Timelike};
use dotenvy::dotenv;
use serde::{Deserialize, Serialize};
use sqlx::postgres::PgPoolOptions;

#[derive(Serialize, Deserialize, Debug)]
struct User {
    id: String,
    name: String,
    phone_number: String,
    email: String
}

async fn get_user_by_birthday(pool: &sqlx::PgPool) -> Result<Vec<User>, sqlx::Error> {
    let today = Local::now().format("%Y-%m-%d %H:%M:%S %z").to_string();
    let today_date = DateTime::parse_from_str(&today, "%Y-%m-%d %H:%M:%S %z").expect("Failed to parsed date");
    let users = sqlx::query_as!(
        User,
        "SELECT id, name, email, phone_number 
        FROM users 
        WHERE date_part('day', birthday) = date_part('day', DATE($1))
        AND date_part('month', birthday) = date_part('month', DATE($1))",
        today_date
    )
    .fetch_all(pool)
    .await?;
    Ok(users)
}

async fn send_users(users: Vec<User>) -> Result<(), Box<dyn std::error::Error>> {
    let client = redis::Client::open("redis://127.0.0.1/")?;
    let mut connection = client.get_connection()?;

    for user in users {
        let user_info = serde_json::to_string(&user)?;
        connection.publish("birthdays_channel", user_info)?;
    }
    Ok(())
} 

async fn scheduler(pool: sqlx::PgPool) {
    let target_hour: u32 = env::var("DAILY_HOUR").expect("DB Url must be set on the .env").parse().expect("hour format on .env wrong");
    let target_minute: u32 = env::var("DAILY_MINUTE").expect("DB Url must be set on the .env").parse().expect("hour format on .env wrong");
    let target_second: u32 = env::var("DAILY_SECOND").expect("DB Url must be set on the .env").parse().expect("hour format on .env wrong");
    loop {
        let now = Local::now();

        let next_activation = if now.hour() < target_hour
            || (now.hour() == target_hour && now.minute() < target_minute)
            || (now.hour() == target_hour
                && now.minute() == target_minute
                && now.second() < target_second)
        {
            // If current time is before target time, calculate time until next target time
            // let mut next_activation = Local::today().and_hms_opt(target_hour, target_minute, target_second).unwrap();
            let mut next_activation = Local::now()
                .with_hour(target_hour)
                .unwrap()
                .with_minute(target_minute)
                .unwrap()
                .with_second(target_second)
                .unwrap();
            if next_activation <= now {
                next_activation = next_activation + Duration::days(1);
            }
            next_activation - now
        } else {
            // If current time is after target time, calculate time until next target time tomorrow
            let mut next_activation = Local::now()
                .with_hour(target_hour)
                .unwrap()
                .with_minute(target_minute)
                .unwrap()
                .with_second(target_second)
                .unwrap() + Duration::days(1);
            if next_activation <= now {
                next_activation = next_activation + Duration::days(1);
            }
            next_activation - now
        };

        let duration = next_activation.to_std().expect("failed to convert to std duration");

        sleep_until(Instant::now() + duration).await;

        match get_user_by_birthday(&pool).await {
            Ok(users) => {
                if !users.is_empty() {
                    println!("{:?}", users);
                    if let Err(e) = send_users(users).await {
                        eprintln!("failed to publish users: {}", e)
                    }
                }
            },
            Err(e) => {
                eprintln!("failed to fetch users: {}", e);
            }
        }
    }
}

#[tokio::main]
async fn main() {
    dotenv().ok();
    let database_url = env::var("DATABASE_URL").expect("DB Url must be set on the .env");
    // let database_url = "postgres://sayakayadev:sahidsudirman@localhost:5432/promo";

    let pool = PgPoolOptions::new().max_connections(5).connect(&database_url).await.expect("failed to create pool");
    scheduler(pool).await;
}
