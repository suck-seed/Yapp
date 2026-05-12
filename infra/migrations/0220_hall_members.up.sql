CREATE TABLE hall_members (
    id uuid PRIMARY KEY,
    hall_id uuid NOT NULL REFERENCES halls (id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    role_id uuid REFERENCES roles (id) ON DELETE RESTRICT,
    nickname text,

    -- Sidebar pinning / Zen Essentials-like behavior
    is_pinned boolean NOT NULL DEFAULT false,
    pinned_position double precision,


    joined_at timestamptz NOT NULL DEFAULT now (),
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now (),
    UNIQUE (hall_id, user_id),

    CONSTRAINT hall_members_pin_position_check
    CHECK (
        (is_pinned = true AND pinned_position IS NOT NULL)
        OR
        (is_pinned = false AND pinned_position IS NULL)
    )
);


-- Helps GET /halls return:
-- pinned halls first by pinned_position,
-- then unpinned halls sorted by hall name from joined halls query.
CREATE INDEX hall_members_user_sidebar_idx
ON hall_members (user_id, is_pinned DESC, pinned_position ASC, hall_id);
