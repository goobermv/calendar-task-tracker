CREATE TABLE IF NOT EXISTS tasks (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(15) DEFAULT 'pending',
    due_date TIMESTAMP WITH TIME ZONE, 
    priority VARCHAR(15) DEFAULT 'low',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_tasks_user_id ON tasks(user_id);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_due_date ON tasks(due_date);
CREATE INDEX idx_tasks_user_status ON tasks(user_id, status);

CREATE TRIGGER update_tasks_updated_at
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE tasks ADD CONSTRAINT chk_tasks_status
    CHECK (status IN ('pending', 'in_progress', 'completed', 'cancelled'));

ALTER TABLE tasks ADD CONSTRAINT chk_tasks_priority
    CHECK (priority IN ('low', 'medium', 'high', 'urgent'));

COMMENT ON TABLE tasks IS 'Stores user tasks and to-dos';
COMMENT ON COLUMN tasks.id IS 'Unique identifier for tasks';
COMMENT ON COLUMN tasks.user_id IS 'Identifier of user who created the task';
COMMENT ON COLUMN tasks.title IS 'Name of the task';
COMMENT ON COLUMN tasks.description IS 'Description of the task';
COMMENT ON COLUMN tasks.status IS 'Task status: pending, in progress, completed, cancelled';
COMMENT ON COLUMN tasks.due_date IS 'Due date of the given task';
COMMENT ON COLUMN tasks.priority IS 'Task priority: low, medium, high, urgent';
COMMENT ON COLUMN tasks.created_at IS 'Time of task being created';
COMMENT ON COLUMN tasks.updated_at IS 'Time of task information being updated';
