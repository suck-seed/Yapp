CREATE TABLE halls (
    id uuid PRIMARY KEY,
    name text NOT NULL,
    icon_url text,
    banner_color text DEFAULT '#1f1f1f',
    description text,
    owner uuid NOT NULL REFERENCES users (id) ON DELETE SET NULL,
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now()
);
