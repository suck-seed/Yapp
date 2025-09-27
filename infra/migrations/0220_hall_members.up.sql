CREATE TABLE hall_members (
    id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    username text NOT NULL UNIQUE,
    joined_at timestamptz NOT NULL DEFAULT now (),

    role_id uuid REFERENCES roles (id) ON DELETE RESTRICT,

    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    UNIQUE (hall_id, user_id)
);
