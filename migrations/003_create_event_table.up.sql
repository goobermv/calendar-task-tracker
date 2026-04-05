CREATE TABLE IF NOT EXISTS events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,
    end_time TIMESTAMP WITH TIME ZONE NOT NULL,
    event_type VARCHAR(20), 
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_user_id ON events(user_id);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_user_start_time ON events(user_id, start_time);

CREATE TRIGGER update_events_updated_at 
    BEFORE UPDATE ON events 
    FOR EACH ROW 
    EXECUTE FUNCTION update_updated_at_column();

ALTER TABLE events ADD CONSTRAINT chk_events_time_range 
    CHECK (end_time > start_time);

ALTER TABLE events ADD CONSTRAINT chk_event_type
    CHECK (event_type IN ('formal', 'celebration', 'important', 'personal', 'holiday', 'casual', 'business'));

COMMENT ON TABLE events IS 'Stores calendar events and appointments';
COMMENT ON COLUMN events.id IS 'Unique identifier for events';
COMMENT ON COLUMN events.user_id IS 'Identifier of user who created the event';
COMMENT ON COLUMN events.title IS 'Name of the event';
COMMENT ON COLUMN events.description IS 'Description of the event';
COMMENT ON COLUMN events.start_time IS 'Event start timestamp';
COMMENT ON COLUMN events.end_time IS 'Event end timestamp (must be after start_time)';
COMMENT ON COLUMN events.event_type IS 'Category of the event';
COMMENT ON COLUMN events.created_at IS 'Time of event being created';
COMMENT ON COLUMN events.updated_at IS 'Time of event information being updated';