CREATE TABLE roles (
    id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    name text NOT NULL,
    color text,
    icon_url text,
    is_default boolean NOT NULL DEFAULT false, -- "everyone"
    is_admin boolean NOT NULL DEFAULT false, -- full permissions

    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);

CREATE UNIQUE INDEX roles_unique_name_per_hall ON roles (hall_id, lower(name));

-- each hall should have unique hall name
-- thus (hall_id, lower(name))
-- since "sandesh" should be = "Sandesh"
