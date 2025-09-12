CREATE TABLE message_mentions (
    message_id uuid NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    PRIMARY KEY (message_id, user_id)
);
