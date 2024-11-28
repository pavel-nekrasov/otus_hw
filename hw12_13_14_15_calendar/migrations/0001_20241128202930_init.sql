-- +goose Up
CREATE table events (
    id              varchar(255) primary key,
    title           varchar(255) not null,
    start_time      timestamptz  not null,
    end_time        timestamptz  not null,
    description     text,
    notify_before   varchar(255),
    owner_email     varchar(255) not null
);

INSERT INTO events (id, title, start_time, end_time, description, notify_before, owner_email)
VALUES
    ('event-1', 'event-1', '2024-11-28 12:00:00+00', '2024-11-28 12:30:00+00', 'description 1', null, 'user1@example.com'),
    ('event-2', 'event-2', '2024-11-28 12:00:00+00', '2024-11-28 12:30:00+00', 'description 2', null, 'user2@example.com');

-- +goose Down
drop table events;