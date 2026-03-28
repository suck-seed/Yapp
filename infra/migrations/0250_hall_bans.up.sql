CREATE TABLE hall_bans (
    id uuid,
    reason text NOT NULL,

    user_id uuid NOT NULL REFERENCES users (id) on DELETE CASCADE,
    hall_id uuid NOT NULL REFERENCES halls (id) on DELETE CASCADE,

    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),

    PRIMARY KEY (user_id, hall_id)
);


-- Indexes for performance
CREATE INDEX idx_hall_bans_hall_id ON hall_bans(hall_id);
CREATE INDEX idx_hall_bans_user_id ON hall_bans(user_id);
CREATE INDEX idx_hall_bans_created_at ON hall_bans(hall_id, created_at DESC);
