DROP TRIGGER IF EXISTS users_notify_event on users;
CREATE TRIGGER users_notify_event
AFTER INSERT OR UPDATE OR DELETE ON users
    FOR EACH ROW EXECUTE PROCEDURE notify_event();
