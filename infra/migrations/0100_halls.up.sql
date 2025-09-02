CREATE TABLE halls (
    hall_id uuid PRIMARY KEY,
    name text NOT NULL,
    icon_url text,
    banner_color text DEFAULT "#1f1f1f", -- 10 colors in frontend
    description text,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    created_by_id uuid NOT NULL REFERENCES users (user_id) ON DELETE RESTRICT
);
