CREATE TABLE reactions (
    id uuid PRIMARY KEY ,
    message_id uuid NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    emoji text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    UNIQUE (message_id, user_id, emoji)
);

-- each emoji has a unique code, we store that in text form
-- instead of trying to store them in some other ways
