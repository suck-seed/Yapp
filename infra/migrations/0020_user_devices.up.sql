CREATE TABLE user_devices (
    user_device_id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    device_id text NOT NULL,
    device_name text,
    last_ip inet,
    last_seen timestamptz NOT NULL DEFAULT now (),
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    UNIQUE (user_id, device_id)
);

-- id -> unique row
-- device_id -> fingerprint, token, hardware_id
