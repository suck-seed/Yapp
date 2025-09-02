CREATE TABLE attachments (
    attachment_id uuid PRIMARY KEY ,
    message_id uuid NOT NULL REFERENCES messages (message_id) ON DELETE CASCADE,
    filename text NOT NULL,
    size_bytes bigint NOT NULL,
    mime_type text,
    url text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);
