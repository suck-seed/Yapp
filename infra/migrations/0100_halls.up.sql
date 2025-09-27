CREATE TABLE halls (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    icon_url text,
    icon_thumbnail_url text,
    banner_color text DEFAULT '#1f1f1f',
    description text,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),
    created_by_id uuid NOT NULL REFERENCES users (id) ON DELETE RESTRICT
);
