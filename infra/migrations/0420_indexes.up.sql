-- other index, goes here
CREATE INDEX IF NOT EXISTS rooms_hall_name_idx ON rooms (hall_id, name);

CREATE INDEX IF NOT EXISTS hall_members_hall_idx ON hall_members (hall_id);

CREATE INDEX IF NOT EXISTS hall_members_user_idx ON hall_members (user_id);

CREATE INDEX IF NOT EXISTS reactions_msg_idx ON reactions (message_id);
