CREATE TABLE attachments (
    id uuid PRIMARY KEY ,
    message_id uuid NOT NULL REFERENCES messages (id) ON DELETE CASCADE,

    file_name text NOT NULL,
    url text NOT NULL,
    file_type text,
    file_size INT,

    created_at timestamptz NOT NULL DEFAULT now (),
    updated_at timestamptz NOT NULL DEFAULT now ()
);
