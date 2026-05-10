CREATE TABLE floor_members (
    floor_id uuid NOT NULL REFERENCES floors (id) ON DELETE CASCADE,
    member_id uuid NOT NULL REFERENCES hall_members (id) ON DELETE CASCADE,
    PRIMARY KEY (floor_id, member_id)
);

CREATE INDEX IF NOT EXISTS floor_members_member_idx ON floor_members (member_id);

-- nothing much, just room_id and user_id combo
