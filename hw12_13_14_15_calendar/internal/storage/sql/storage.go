package sqlstorage

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/configs"
	m "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/models"
	"time"
)

type Storage struct {
	driver, source string
	conn           *sql.Conn
	db             *sqlx.DB
}

func New(config configs.Config) *Storage {
	// TODO: env
	return &Storage{driver: config.Database.Driver, source: config.Database.Source}
}

func (s *Storage) Connect(ctx context.Context) error {
	db, err := sqlx.Connect(s.driver, s.source)
	if err != nil {
		return err
	}
	s.db = db

	conn, err := db.Conn(ctx)
	s.conn = conn

	return err
}

func (s *Storage) Close(ctx context.Context) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	return s.conn.Close()
}

func (s *Storage) Add(ev m.Event) (string, error) {
	query := `INSERT INTO events (title, description, user_id, date, duration, notified_before)
VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING id`

	row := s.db.QueryRow(query, ev.Title, ev.Description, ev.UserID, ev.DateTime, ev.Duration, ev.NotifiedBefore)

	var id string

	err := row.Scan(&id)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Storage) Update(eventID string, updatedEvent m.Event) error {
	query := `UPDATE events SET title = $1, 
                  description = $2, user_id = $3, date=$4, duration=$5, notified_before=$6
	WHERE id = $7`

	_, err := s.db.Exec(query, updatedEvent.Title, updatedEvent.Description, updatedEvent.UserID,
		updatedEvent.DateTime, updatedEvent.Duration, updatedEvent.NotifiedBefore, eventID)

	return err
}

func (s *Storage) Delete(eventID string) error {
	query := `DELETE FROM events
WHERE id = $1`

	_, err := s.db.Exec(query, eventID)
	return err
}

func (s *Storage) GetEventsForDay(day time.Time) []m.Event {
	query := `SELECT * FROM events
WHERE DATE_PART('day', date) = $1`

	return s.getEventsByQueryAndArgs(query, day.YearDay())
}

func (s *Storage) GetEventsForWeek(fromDay time.Time) []m.Event {
	query := `SELECT * FROM events
WHERE DATE_PART('year', date) = $1 AND DATE_PART('month', date) = $2 AND date >= $3 AND DATE_PART('week', date) = $4`

	_, week := fromDay.ISOWeek()
	return s.getEventsByQueryAndArgs(query, fromDay.Year(), fromDay.Month(), fromDay, week)
}

func (s *Storage) GetEventsForMonth(fromDay time.Time) []m.Event {
	query := `SELECT * FROM events
WHERE DATE_PART('year', date) = $1 AND DATE_PART('month', date) = $2 AND date >= $3`

	return s.getEventsByQueryAndArgs(query, fromDay.Year(), fromDay.Month(), fromDay)
}

func (s *Storage) getEventsByQueryAndArgs(query string, args ...interface{}) []m.Event {
	eventsForDay := make([]m.Event, 0)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return eventsForDay
	}
	defer rows.Close()

	var (
		id, title, description, userID, duration, notifiedBefore string
		date                                                     time.Time
	)

	for rows.Next() {
		err = rows.Scan(&id, &title, &description, &userID, &date, &duration, &notifiedBefore)
		if err != nil {
			break
		}

		event := m.Event{ID: id, Title: title, Description: description, UserID: userID, DateTime: date}

		dur, err := time.ParseDuration(duration)
		if err != nil {
			event.Duration = dur
		}

		pushBefore, err := time.ParseDuration(notifiedBefore)
		if err != nil {
			event.NotifiedBefore = pushBefore
		}

		eventsForDay = append(eventsForDay, event)
	}

	return eventsForDay
}
