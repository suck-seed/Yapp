CREATE TABLE hall_requests (
    id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    CONSTRAINT hall_requests_unique_hall_user UNIQUE (hall_id, user_id)
);
