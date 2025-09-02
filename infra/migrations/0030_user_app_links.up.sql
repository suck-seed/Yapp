CREATE TABLE user_app_links (
    user_app_link_id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    provider app_provider NOT NULL,
    url text NOT NULL,
    show_on_profile boolean NOT NULL DEFAULT true,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    UNIQUE (user_id, provider)
);

-- user can have only 1 external app provider
