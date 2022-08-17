package memorystorage

import (
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

func (st *Storage) Add(ev m.Event) error {
	st.mu.Lock()
	defer st.mu.Unlock()

	_, ok := st.events[ev.ID]
	if ok {
		return storage.ErrEventAlreadyExists
	}

	st.events[ev.ID] = &ev
	return nil
}

func (st *Storage) Update(eventId string, updatedEvent m.Event) {
	st.mu.Lock()
	defer st.mu.Unlock()

	event, ok := st.events[eventId]
	if !ok {
		return
	}

	event.Duration = updatedEvent.Duration
	event.UserID = updatedEvent.UserID
	event.Description = updatedEvent.Description
	event.DateTime = updatedEvent.DateTime
	event.NotifiedBefore = updatedEvent.NotifiedBefore
	event.Title = updatedEvent.Title
}

func (st *Storage) Delete(eventId string) {
	st.mu.Lock()
	defer st.mu.Unlock()

	delete(st.events, eventId)
}

func (st *Storage) GetEventsForDay(day time.Time) []m.Event {
	st.mu.RLock()
	defer st.mu.RUnlock()

	eventsForDay := make([]m.Event, 0)
	for _, event := range st.events {
		eventDay := event.DateTime

		if sameYearDay(day, eventDay) {
			eventsForDay = append(eventsForDay, *event)
		}
	}

	return eventsForDay
}

func (st *Storage) GetEventsForWeek(fromDay time.Time) []m.Event {
	st.mu.RLock()
	defer st.mu.RUnlock()

	eventsForWeek := make([]m.Event, 0)
	for _, event := range st.events {
		eventDay := event.DateTime

		if !sameYearMonth(fromDay, eventDay) || !sameWeek(fromDay, eventDay) {
			continue
		}

		if sameDayOrAfter(eventDay, fromDay) {
			eventsForWeek = append(eventsForWeek, *event)
		}

	}

	return eventsForWeek
}

func (st *Storage) GetEventsForMonth(fromDay time.Time) []m.Event {
	st.mu.RLock()
	defer st.mu.RUnlock()

	eventsForMonth := make([]m.Event, 0)
	for _, event := range st.events {
		eventDay := event.DateTime

		if !sameYearMonth(fromDay, eventDay) {
			continue
		}

		if sameDayOrAfter(eventDay, fromDay) {
			eventsForMonth = append(eventsForMonth, *event)
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
