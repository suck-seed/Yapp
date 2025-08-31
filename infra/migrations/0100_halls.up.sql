CREATE TABLE halls (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    icon_url text,
    banner_color text, -- 10 colors in frontend
    description text,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    created_by_id uuid NOT NULL REFERENCES users (id) ON DELETE RESTRICT
);
