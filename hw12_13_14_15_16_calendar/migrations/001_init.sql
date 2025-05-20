-- +goose Up
CREATE TABLE IF NOT EXISTS events (
                                      id VARCHAR PRIMARY KEY,
                                      title TEXT NOT NULL,
                                      start_time TIMESTAMPTZ NOT NULL,
                                      end_time TIMESTAMPTZ NOT NULL,
                                      description TEXT,
                                      user_id VARCHAR NOT NULL,
                                      notify_before INTEGER,
                                      CONSTRAINT valid_time CHECK (end_time > start_time)
);

CREATE INDEX IF NOT EXISTS idx_user_id ON events(user_id);
CREATE INDEX IF NOT EXISTS idx_user_start ON events(user_id, start_time);

-- +goose Down
DROP TABLE IF EXISTS events;
DROP INDEX IF EXISTS idx_user_id;
DROP INDEX IF EXISTS idx_user_start;