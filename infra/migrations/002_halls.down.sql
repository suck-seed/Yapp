CREATE TABLE IF NOT EXISTS halls (
    hall_id uuid,
    hall_name VARCHAR(255) NOT NULL UNIQUE,
    avatar TEXT,
);
