DROP TRIGGER IF EXISTS update_user_updated_at ON users;

DROP FUNCTION IF EXISTS update_updated_at_column();

DROP INDEX IF EXISTS idx_user_email;
DROP INDEX IF EXISTS idx_user_username;

DROP TABLE IF EXISTS users;