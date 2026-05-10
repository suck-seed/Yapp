CREATE TABLE room_members (
    room_id uuid NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    member_id uuid NOT NULL REFERENCES hall_members (id) ON DELETE CASCADE,
    PRIMARY KEY (room_id, member_id)
);

CREATE INDEX IF NOT EXISTS room_members_member_idx ON room_members (member_id);
-- nothing much, just room_id and user_id combo
