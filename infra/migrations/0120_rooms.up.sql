CREATE TABLE rooms (
    room_id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls (hall_id) ON DELETE CASCADE,
    floor_id uuid REFERENCES floors (id) ON DELETE SET NULL,
    name text NOT NULL,
    room_type room_type NOT NULL,
    is_private boolean NOT NULL DEFAULT false,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);

CREATE INDEX rooms_hall_id_idx ON rooms (hall_id);

-- to retrieve rooms based on hall_id, faster
