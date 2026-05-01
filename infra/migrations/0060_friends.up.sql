CREATE TABLE IF NOT EXISTS friends (
    user_id_1 uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    user_id_2 uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (user_id_1, user_id_2),
    -- Ensure user_id_1 is always the smaller UUID to prevent (A,B) and (B,A) duplicates
    CONSTRAINT user_id_order_check CHECK (user_id_1 < user_id_2)
);

-- Index for scanning friendships where the user is the second ID
CREATE INDEX IF NOT EXISTS idx_friends_user_2 ON friends(user_id_2);
