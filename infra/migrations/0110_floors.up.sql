CREATE TABLE floors (
    id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    name text NOT NULL,
    position int NOT NULL DEFAULT 0,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);
