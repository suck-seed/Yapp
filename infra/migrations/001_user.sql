CREATE TABLE IF NOT EXISTS users (
    user_id uuid,
    username VARCHAR(255) NOT NULL UNIQUE,
    avatar TEXT,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    locale VARCHAR(255),
    number VARCHAR(50) UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,


    PRIMARY KEY (user_id)
);
