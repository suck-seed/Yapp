CREATE TABLE messages (
    id uuid PRIMARY KEY ,
    room_id uuid NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    author_id uuid NOT NULL REFERENCES users (id) ON DELETE RESTRICT,
    content text NOT NULL,
    sent_at timestamptz NOT NULL DEFAULT now (),
    edited_at timestamptz,
    deleted_at timestamptz,
    mention_everyone boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);

-- pagination-friendly (room_id, sent_at, id)
CREATE INDEX messages_room_time_idx ON messages (room_id, sent_at);

CREATE INDEX messages_room_sent_id_idx ON messages (room_id, sent_at, id);
