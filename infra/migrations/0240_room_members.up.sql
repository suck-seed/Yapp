CREATE TABLE room_members (
    room_id uuid NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (room_id, user_id)
);

-- nothing much, just room_id and user_id combo
