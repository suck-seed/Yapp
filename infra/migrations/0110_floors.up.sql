CREATE TABLE floors (
    floor_id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls (hall_id) ON DELETE CASCADE,
    name text NOT NULL,
    -- position int NOT NULL DEFAULT 0,
    is_private boolean NOT NULL DEFAULT false, -- prompt for user, send userID and floorID, add to floor_members
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);

CREATE INDEX floors_hall_id_idx ON floors (hall_id);
