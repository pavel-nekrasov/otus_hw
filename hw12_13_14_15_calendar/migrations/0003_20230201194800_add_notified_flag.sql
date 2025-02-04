-- +goose Up
ALTER TABLE events
ADD notified_flag boolean default false;

-- +goose Down
ALTER TABLE events
DROP notified_flag boolean;