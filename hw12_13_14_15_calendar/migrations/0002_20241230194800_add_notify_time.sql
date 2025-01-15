-- +goose Up
ALTER TABLE events
ADD notify_time timestamptz;

-- +goose Down
ALTER TABLE events
DROP notify_time timestamptz;