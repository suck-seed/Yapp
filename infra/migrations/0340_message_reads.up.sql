CREATE TABLE IF NOT EXISTS message_reads (
    room_id uuid NOT NULL REFERENCES rooms(id) ON DELETE CASCADE,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message_id uuid NOT NULL REFERENCES messages(id) ON DELETE CASCADE,

    read_at timestamptz NOT NULL DEFAULT now(),
    created_at timestamptz NOT NULL DEFAULT now(),
    updated_at timestamptz NOT NULL DEFAULT now(),

    PRIMARY KEY (room_id, user_id)
);

CREATE INDEX IF NOT EXISTS message_reads_room_id_idx
    ON message_reads(room_id);

CREATE INDEX IF NOT EXISTS message_reads_message_id_idx
    ON message_reads(message_id);

CREATE INDEX IF NOT EXISTS message_reads_user_id_idx
    ON message_reads(user_id);

DROP TRIGGER IF EXISTS message_reads_set_updated_at ON message_reads;

CREATE TRIGGER message_reads_set_updated_at
BEFORE UPDATE ON message_reads
FOR EACH ROW EXECUTE FUNCTION set_updated_at();
