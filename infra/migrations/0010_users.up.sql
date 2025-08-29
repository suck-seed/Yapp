CREATE TABLE users (
    id uuid PRIMARY KEY,
    username text NOT NULL UNIQUE,
    display_name text NOT NULL,
    email text NOT NULL UNIQUE,
    password_hash text NOT NULL,
    phone_number text UNIQUE,
    avatar_url text,
    friend_policy friend_policy DEFAULT 'everyone',
    active boolean NOT NULL DEFAULT false, -- consider derived via Redis; keep for compatibility
    last_seen timestamptz,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);

CREATE UNIQUE INDEX users_username_lower_uq ON users (lower(username));

CREATE UNIQUE INDEX users_email_lower_uq ON users (lower(email));

-- "Sandesh" and "sandesh" should be treated as same username
-- so, we create a lower(username) copy of each username in the index
-- thus, we can compare from this index for valid unique username
--
-- Same goes for email
--
-- SELECT * FROM users WHERE lower(username) = lower('Alice');
