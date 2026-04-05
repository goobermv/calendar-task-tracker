ALTER TABLE events DROP CONSTRAINT IF EXISTS chk_events_time_range;
ALTER TABLE events DROP CONSTRAINT IF EXISTS chk_event_type;

DROP TRIGGER IF EXISTS update_events_updated_at ON events;

DROP INDEX IF EXISTS idx_events_user_id;
DROP INDEX IF EXISTS idx_events_start_time;
DROP INDEX IF EXISTS idx_events_user_start_time;

DROP TABLE IF EXISTS events;