CREATE TABLE role_permissions (
    role_id uuid PRIMARY KEY REFERENCES roles (role_id) ON DELETE CASCADE,
    -- General
    view_channels boolean NOT NULL DEFAULT true,
    manage_channels boolean NOT NULL DEFAULT false,
    manage_roles boolean NOT NULL DEFAULT false,
    manage_servers boolean NOT NULL DEFAULT false,
    change_nickname boolean NOT NULL DEFAULT true,
    manage_nicknames boolean NOT NULL DEFAULT false,
    kick_members boolean NOT NULL DEFAULT false,
    ban_members boolean NOT NULL DEFAULT false,
    -- Text
    text_send_messages boolean NOT NULL DEFAULT true,
    text_attach_files boolean NOT NULL DEFAULT true,
    text_mention_roles boolean NOT NULL DEFAULT true,
    text_manage_messages boolean NOT NULL DEFAULT false,
    text_read_history boolean NOT NULL DEFAULT true,
    text_send_voice boolean NOT NULL DEFAULT true,
    -- Voice
    voice_connect boolean NOT NULL DEFAULT true,
    voice_speak boolean NOT NULL DEFAULT true,
    voice_video boolean NOT NULL DEFAULT false,
    voice_mute_members boolean NOT NULL DEFAULT false
);

-- no need for unique id, each role_id is unique
