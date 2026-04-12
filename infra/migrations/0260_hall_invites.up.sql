CREATE TABLE hall_invites (
    id          UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    hall_id     UUID        NOT NULL REFERENCES halls(id)  ON DELETE CASCADE,
    created_by  UUID        NOT NULL REFERENCES users(id)  ON DELETE CASCADE,
    code        VARCHAR(12) NOT NULL,
    role_id     UUID        NULL     REFERENCES roles(id)  ON DELETE SET NULL,
    max_uses    INT         NULL,           -- NULL  → unlimited
    used_count  INT         NOT NULL DEFAULT 0,
    expires_at  TIMESTAMPTZ NULL,           -- NULL  → never expires
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    CONSTRAINT hall_invites_used_count_non_negative CHECK (used_count >= 0),
    CONSTRAINT hall_invites_max_uses_positive       CHECK (max_uses IS NULL OR max_uses > 0)
);

CREATE UNIQUE INDEX idx_hall_invites_code    ON hall_invites(code);

-- fast lookup of all invites for a hall (settings page)
CREATE INDEX idx_hall_invites_hall_id ON hall_invites(hall_id);
