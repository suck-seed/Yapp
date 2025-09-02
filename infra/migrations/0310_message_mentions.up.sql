CREATE TABLE message_mentions (
    message_id uuid NOT NULL REFERENCES messages (message_id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    PRIMARY KEY (message_id, user_id)
);
