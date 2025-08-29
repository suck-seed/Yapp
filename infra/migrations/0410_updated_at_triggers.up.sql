CREATE TRIGGER users_set_updated_at BEFORE
UPDATE ON users FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER user_devices_set_updated_at BEFORE
UPDATE ON user_devices FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER user_app_links_set_updated_at BEFORE
UPDATE ON user_app_links FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER halls_set_updated_at BEFORE
UPDATE ON halls FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER floors_set_updated_at BEFORE
UPDATE ON floors FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER rooms_set_updated_at BEFORE
UPDATE ON rooms FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER roles_set_updated_at BEFORE
UPDATE ON roles FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER hall_members_set_updated_at BEFORE
UPDATE ON hall_members FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER messages_set_updated_at BEFORE
UPDATE ON messages FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER attachments_set_updated_at BEFORE
UPDATE ON attachments FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

CREATE TRIGGER reactions_set_updated_at BEFORE
UPDATE ON reactions FOR EACH ROW EXECUTE FUNCTION set_updated_at ();

-- automatically calls the trigger function for any update task on the
-- corresponding tables
