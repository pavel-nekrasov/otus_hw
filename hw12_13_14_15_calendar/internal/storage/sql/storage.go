package sqlstorage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/stdlib" // need import pgx
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/customerrors"
	"github.com/pavel-nekrasov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	goose "github.com/pressly/goose/v3"
)

type Storage struct {
	dsn string
	db  *sql.DB
}

func New(host string, port int, dbname, user, password string) *Storage {
	return &Storage{
		dsn: fmt.Sprintf("postgres://%v:%v@%v:%v/%v", user, password, host, port, dbname),
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	var err error
	s.db, err = sql.Open("pgx", s.dsn)
	if err != nil {
		return fmt.Errorf("cannot open pgx driver: %w", err)
	}

	return s.db.PingContext(ctx)
}

func (s *Storage) Close(_ context.Context) error {
	return s.db.Close()
}

func (s *Storage) Migrate(_ context.Context, migrate string) (err error) {
	//	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("cannot set dialect: %w", err)
	}

	if err := goose.Up(s.db, migrate); err != nil {
		return fmt.Errorf("cannot do up migration: %w", err)
	}

	return nil
}

func (s *Storage) AddEvent(ctx context.Context, event storage.Event) error {
	res, err := s.db.ExecContext(ctx, `INSERT INTO events 
		(id, title, start_time, end_time, description, notify_before, owner_email) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.NotifyBefore,
		event.OwnerEmail,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", event.ID)}
	}

	return nil
}

func (s *Storage) UpdateEvent(ctx context.Context, event storage.Event) error {
	res, err := s.db.ExecContext(ctx, `UPDATE events
		SET title = $2, start_time = $3, end_time = $4, description = $5, notify_before = $6, owner_email = $7 
		WHERE id = $1`,
		event.ID,
		event.Title,
		event.StartTime,
		event.EndTime,
		event.Description,
		event.NotifyBefore,
		event.OwnerEmail,
	)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", event.ID)}
	}

	return nil
}

func (s *Storage) GetEvent(ctx context.Context, eventID string) (storage.Event, error) {
	row := s.db.QueryRowContext(ctx,
		"SELECT id, title, start_time, end_time, description, notify_before, owner_email FROM events WHERE id = $1",
		eventID,
	)
	if errors.Is(row.Err(), sql.ErrNoRows) {
		return storage.Event{}, customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", eventID)}
	}
	if row.Err() != nil {
		return storage.Event{}, row.Err()
	}
	var event storage.Event
	var notify, description sql.NullString
	err := row.Scan(&event.ID,
		&event.Title,
		&event.StartTime,
		&event.EndTime,
		&description,
		&notify,
		&event.OwnerEmail,
	)
	if err != nil {
		return storage.Event{}, err
	}

	if description.Valid {
		event.Description = description.String
	}

	if notify.Valid {
		event.NotifyBefore = notify.String
	}

	return event, nil
}

func (s *Storage) DeleteEvent(ctx context.Context, eventID string) error {
	res, err := s.db.ExecContext(ctx, "delete from events where id = $1", eventID)
	if err != nil {
		return err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if cnt == 0 {
		return customerrors.NotFound{Message: fmt.Sprintf("Event with id = \"%v\" not found", eventID)}
	}

	return nil
}

func (s *Storage) ListEventsForPeriod(
	ctx context.Context,
	ownerEmail string,
	startDate,
	endDate time.Time,
) ([]storage.Event, error) {
	result := make([]storage.Event, 0)
	rows, err := s.db.QueryContext(ctx,
		`SELECT id, title, start_time, end_time, description, notify_before, owner_email 
		FROM events 
		WHERE owner_email = $1 AND start_time >= $2 AND end_time <= $3`,
		ownerEmail,
		startDate,
		endDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return result, nil
	}
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var event storage.Event
		var notify, description sql.NullString
		err := rows.Scan(&event.ID,
			&event.Title,
			&event.StartTime,
			&event.EndTime,
			&description,
			&notify,
			&event.OwnerEmail,
		)
		if err != nil {
			return nil, err
		}
		if description.Valid {
			event.Description = description.String
		}

		if notify.Valid {
			event.NotifyBefore = notify.String
		}
		result = append(result, event)
	}

	return result, rows.Err()
}
