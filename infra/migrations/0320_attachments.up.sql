CREATE TABLE attachments (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid (),
    message_id uuid NOT NULL REFERENCES messages (id) ON DELETE CASCADE,
    filename text NOT NULL,
    size_bytes bigint NOT NULL,
    mime_type text,
    url text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);
