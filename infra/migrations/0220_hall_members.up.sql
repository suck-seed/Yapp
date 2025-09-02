CREATE TABLE hall_members (
    hall_member_id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls (hall_id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    nickname text,
    joined_at timestamptz NOT NULL DEFAULT now (),
    role_id uuid NOT NULL REFERENCES roles (id) ON DELETE RESTRICT,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    UNIQUE (hall_id, user_id)
);
