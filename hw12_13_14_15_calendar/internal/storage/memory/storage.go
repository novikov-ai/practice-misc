package memorystorage

import (
	"context"
	"sync"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage"
	m "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/models"
)

type Storage struct {
	mu     sync.RWMutex
	events map[string]*m.Event
}

func New() *Storage {
	return &Storage{events: map[string]*m.Event{}}
}

func (st *Storage) Connect(ctx context.Context) error {
	// dummy implementation for in-memory database
	return clearStorage(ctx, st)
}

func (st *Storage) Close(ctx context.Context) error {
	// dummy implementation for in-memory database
	return clearStorage(ctx, st)
}

func clearStorage(ctx context.Context, st *Storage) error {
	select {
	default:
	case <-ctx.Done():
		return ctx.Err()
	}

	st.events = map[string]*m.Event{}
	return nil
}

func (st *Storage) Add(ctx context.Context, ev m.Event) (string, error) {
	st.mu.Lock()
	defer st.mu.Unlock()

	select {
	case <-ctx.Done():
		break
	default:
		_, ok := st.events[ev.ID]
		if ok {
			return ev.ID, storage.ErrEventAlreadyExists
		}

		st.events[ev.ID] = &ev
	}

	return ev.ID, nil
}

func (st *Storage) Update(ctx context.Context, eventId string, updatedEvent m.Event) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	select {
	case <-ctx.Done():
		break
	default:
		event, ok := st.events[eventId]
		if !ok {
			return storage.ErrEventNotExists
		}

		event.Duration = updatedEvent.Duration
		event.UserID = updatedEvent.UserID
		event.Description = updatedEvent.Description
		event.DateTime = updatedEvent.DateTime
		event.NotifiedBefore = updatedEvent.NotifiedBefore
		event.Title = updatedEvent.Title
	}

	return nil
}

func (st *Storage) Delete(ctx context.Context, eventId string) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	select {
	case <-ctx.Done():
		break
	default:
		countBefore := len(st.events)
		delete(st.events, eventId)

		if len(st.events) >= countBefore {
			return storage.ErrInMemoryOperationFailed
		}
	}

	return nil
}

func (st *Storage) GetEventsForDay(ctx context.Context, day time.Time) []m.Event {
	st.mu.RLock()
	defer st.mu.RUnlock()

	eventsForDay := make([]m.Event, 0)
	select {
	case <-ctx.Done():
		break
	default:
		for _, event := range st.events {
			eventDay := event.DateTime

			if sameYearDay(day, eventDay) {
				eventsForDay = append(eventsForDay, *event)
			}
		}
	}

	return eventsForDay
}

func (st *Storage) GetEventsForWeek(ctx context.Context, fromDay time.Time) []m.Event {
	st.mu.RLock()
	defer st.mu.RUnlock()

	eventsForWeek := make([]m.Event, 0)
	select {
	case <-ctx.Done():
		break
	default:
		for _, event := range st.events {
			eventDay := event.DateTime

			if !sameYearMonth(fromDay, eventDay) || !sameWeek(fromDay, eventDay) {
				continue
			}

			if sameDayOrAfter(eventDay, fromDay) {
				eventsForWeek = append(eventsForWeek, *event)
			}
		}
	}

	return eventsForWeek
}

func (st *Storage) GetEventsForMonth(ctx context.Context, fromDay time.Time) []m.Event {
	st.mu.RLock()
	defer st.mu.RUnlock()

	eventsForMonth := make([]m.Event, 0)
	select {
	case <-ctx.Done():
		break
	default:
		for _, event := range st.events {
			eventDay := event.DateTime

			if !sameYearMonth(fromDay, eventDay) {
				continue
			}

			if sameDayOrAfter(eventDay, fromDay) {
				eventsForMonth = append(eventsForMonth, *event)
			}
		}
	}

	return eventsForMonth
}

func sameYearDay(d1, d2 time.Time) bool {
	return d1.YearDay() == d2.YearDay() && d1.Year() == d2.Year()
}

func sameYearMonth(d1, d2 time.Time) bool {
	return d1.Year() == d2.Year() && d1.Month() == d2.Month()
}

func sameWeek(d1, d2 time.Time) bool {
	_, w1 := d1.ISOWeek()
	_, w2 := d2.ISOWeek()
	return w1 == w2
}

func sameDayOrAfter(day, pastDay time.Time) bool {
	return day.Equal(pastDay) || day.After(pastDay)
}
