CREATE TABLE halls (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    name text NOT NULL,
    icon_url text,
    banner_color text NOT NULL, -- 10 colors in frontend
    description text,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    created_by_id uuid NOT NULL REFERENCES users (user_id) ON DELETE RESTRICT
);
