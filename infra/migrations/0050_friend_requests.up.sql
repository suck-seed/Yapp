CREATE TABLE friend_requests (
    id uuid PRIMARY KEY,
    sender_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    receiver_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    -- Prevent duplicate requests between the same two users
    UNIQUE (sender_id, receiver_id),
    -- Prevent users from sending a request to themselves
    CHECK (sender_id != receiver_id)
);

-- Index to quickly find requests sent to a specific user
CREATE INDEX friend_requests_receiver_idx ON friend_requests (receiver_id);
