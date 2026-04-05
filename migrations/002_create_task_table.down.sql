ALTER TABLE tasks DROP CONSTRAINT IF EXISTS chk_tasks_status;
ALTER TABLE tasks DROP CONSTRAINT IF EXISTS chk_tasks_priority;

DROP TRIGGER IF EXISTS update_tasks_updated_at ON tasks;

DROP INDEX IF EXISTS idx_tasks_user_id;
DROP INDEX IF EXISTS idx_tasks_status;
DROP INDEX IF EXISTS idx_tasks_due_date;
DROP INDEX IF EXISTS idx_tasks_user_status;

DROP TABLE IF EXISTS tasks;