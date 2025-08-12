CREATE TABLE rooms (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    hall_id uuid NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    floor_id uuid REFERENCES floors (id) ON DELETE SET NULL,
    name text NOT NULL,
    room_type room_type NOT NULL,
    is_private boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);

CREATE INDEX rooms_hall_id_idx ON rooms (hall_id);

-- to retrieve rooms based on hall_id, faster
