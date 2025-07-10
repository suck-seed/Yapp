CREATE TABLE IF NOT EXISTS users (
    user_id uuid,
    -- Infos
    username VARCHAR(255) NOT NULL UNIQUE,
    avatar TEXT,
    email VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL,
    locale VARCHAR(255),
    number VARCHAR(50) UNIQUE,
    warning_flag int NOT NULL DEFAULT 10,
    -- misc
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    -- PRIMARY KEY
    PRIMARY KEY (user_id)
);
