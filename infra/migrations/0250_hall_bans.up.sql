CREATE TABLE hall_bans (
    id uuid,
    reason text,
    banned_by uuid NOT NULL,

    user_id uuid NOT NULL REFERENCES users (id) on DELETE CASCADE,
    hall_id uuid NOT NULL REFERENCES halls (id) on DELETE CASCADE,

    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()

    PRIMARY KEY (user_id, hall_id)
);
