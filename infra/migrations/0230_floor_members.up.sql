CREATE TABLE floor_members (
    floor_id uuid NOT NULL REFERENCES floor (floor_id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    PRIMARY KEY (floor_id, user_id)
);

-- nothing much, just room_id and user_id combo
