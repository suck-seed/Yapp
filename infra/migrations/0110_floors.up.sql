CREATE TABLE floors (
    id          uuid             PRIMARY KEY,
    hall_id     uuid             NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    name        text             NOT NULL,
    position    double precision NOT NULL DEFAULT 0,
    is_private  boolean          NOT NULL DEFAULT false,
    created_at  timestamptz      NOT NULL DEFAULT now(),
    updated_at  timestamptz      NOT NULL DEFAULT now()
);

CREATE INDEX floors_hall_id_idx          ON floors (hall_id);
-- composite index speeds up ordered fetches
CREATE INDEX floors_hall_id_position_idx ON floors (hall_id, position);
