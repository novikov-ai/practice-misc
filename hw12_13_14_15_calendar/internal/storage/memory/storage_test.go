package memorystorage

import (
	"errors"
	"testing"
	"time"

	st "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage"
	m "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/stretchr/testify/require"
)

var (
	dates = []time.Time{
		time.Date(2020, 1, 10, 0, 0, 0, 0, time.UTC),
		dateExistsOnlyOne,
		time.Date(2020, 2, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 13, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 14, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 15, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 16, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 2, 17, 0, 0, 0, 0, time.UTC),

		time.Date(2020, 3, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2020, 4, 13, 0, 0, 0, 0, time.UTC),

		time.Date(2021, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 2, 11, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 3, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2021, 4, 13, 0, 0, 0, 0, time.UTC),

		time.Date(2022, 1, 10, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 3, 2, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 3, 9, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 3, 10, 0, 0, 0, 0, time.UTC),
		dateExistMany, dateExistMany, dateExistMany, dateExistMany, dateExistMany,
		time.Date(2022, 3, 12, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 3, 30, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 3, 31, 0, 0, 0, 0, time.UTC),
		time.Date(2022, 2, 11, 0, 0, 0, 0, time.UTC),
	}

	dateNotExisting   = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)
	dateExistsOnlyOne = time.Date(2020, 2, 11, 0, 0, 0, 0, time.UTC)
	dateExistMany     = time.Date(2022, 3, 11, 0, 0, 0, 0, time.UTC)
)

func TestStorageAdd(t *testing.T) {
	testCases := []struct {
		title  string
		events []m.Event
	}{
		{title: "add 5 items to empty storage", events: generateEvents(5)},
		{title: "add 10 items to empty storage", events: generateEvents(10)},
		{title: "add 15 items to empty storage", events: generateEvents(15)},
	}

	storage := New()

	for _, test := range testCases {
		test := test

		t.Run(test.title, func(t *testing.T) {
			for _, ev := range test.events {
				err := storage.Add(ev)
				require.Nil(t, err)
			}
		})
	}
	require.Equal(t, 30, len(storage.events))

	t.Run("add already existing item", func(t *testing.T) {
		for _, ev := range testCases[0].events {
			err := storage.Add(ev)
			require.True(t, errors.Is(err, st.ErrEventAlreadyExists))
		}

		require.Equal(t, 30, len(storage.events))
	})
}

func TestStorageUpdate(t *testing.T) {
	testCases := []struct {
		title         string
		events        []m.Event
		correctUpdate bool
	}{
		{title: "update empty storage", events: generateEvents(0)},
		{title: "update not existing item", events: generateEvents(10)},
		{title: "update existing item", events: generateEvents(10), correctUpdate: true},
	}

	updateEvent := m.New()
	updateEvent.Title = "updating title"
	updateEvent.Description = "simple description"
	updateEvent.UserID = "5"
	updateEvent.Duration = time.Second
	updateEvent.DateTime = time.Now()
	updateEvent.NotifiedBefore = time.Hour

	for _, test := range testCases {
		test := test

		t.Run(test.title, func(t *testing.T) {
			storage := initStorage(test.events)

			if test.correctUpdate {
				updatingID := test.events[0].ID

				storage.Update(updatingID, *updateEvent)
				require.Equal(t, len(test.events), len(storage.events))

				ev, ok := storage.events[updatingID]
				require.True(t, ok)

				require.NotEqual(t, updateEvent.ID, ev.ID)
				require.Equal(t, updateEvent.Title, ev.Title)
				require.Equal(t, updateEvent.Description, ev.Description)
				require.Equal(t, updateEvent.UserID, ev.UserID)
				require.Equal(t, updateEvent.Duration, ev.Duration)
				require.Equal(t, updateEvent.DateTime, ev.DateTime)
				require.Equal(t, updateEvent.NotifiedBefore, ev.NotifiedBefore)

			} else {
				storage.Update(updateEvent.ID, *updateEvent)
				require.Equal(t, len(test.events), len(storage.events))

				for _, ev := range test.events {
					require.NotEqual(t, updateEvent.ID, ev.ID)
					require.NotEqual(t, updateEvent.Title, ev.Title)
					require.NotEqual(t, updateEvent.Description, ev.Description)
					require.NotEqual(t, updateEvent.UserID, ev.UserID)
					require.NotEqual(t, updateEvent.Duration, ev.Duration)
					require.NotEqual(t, updateEvent.DateTime, ev.DateTime)
					require.NotEqual(t, updateEvent.NotifiedBefore, ev.NotifiedBefore)
				}
			}
		})
	}
}

func TestStorageDelete(t *testing.T) {
	testCases := []struct {
		title         string
		events        []m.Event
		removeItems   int
		correctDelete bool
	}{
		{title: "delete from empty storage", events: generateEvents(0)},
		{title: "delete not existing item", events: generateEvents(10)},
		{title: "delete existing item", events: generateEvents(15), removeItems: 5, correctDelete: true},
		{title: "delete all items from storage", events: generateEvents(15), removeItems: 15, correctDelete: true},
	}

	for _, test := range testCases {
		test := test

		t.Run(test.title, func(t *testing.T) {
			storage := initStorage(test.events)
			require.Equal(t, len(test.events), len(storage.events))

			if test.correctDelete {
				for i := 0; i < test.removeItems; i++ {
					storage.Delete(test.events[i].ID)
				}
				require.Equal(t, len(test.events)-test.removeItems, len(storage.events))
			} else {
				storage.Delete("not existing UUID")
				require.Equal(t, len(test.events), len(storage.events))
			}
		})
	}
}

func TestStorageGetEventsForDay(t *testing.T) {
	storage := setUpStorage()

	testCase := []struct {
		title  string
		day    time.Time
		events int
	}{
		{title: "get not existing event", day: dateNotExisting, events: 0},
		{title: "get existing event", day: dateExistsOnlyOne, events: 1},
		{title: "get a few existing events", day: dateExistMany, events: 5},
	}

	for _, test := range testCase {
		test := test

		t.Run(test.title, func(t *testing.T) {
			events := storage.GetEventsForDay(test.day)
			require.Equal(t, test.events, len(events))

			for _, ev := range events {
				require.True(t, test.day.Equal(ev.DateTime))
			}
		})
	}
}

func TestStorageGetEventsForWeek(t *testing.T) {
	storage := setUpStorage()

	testCase := []struct {
		title  string
		week   time.Time
		events int
	}{
		{title: "get not existing event", week: dateNotExisting, events: 0},
		{title: "get existing events at the week", week: dateExistsOnlyOne, events: 6},
		{title: "get a few similar events + event at the week", week: dateExistMany, events: 6},
	}

	for _, test := range testCase {
		test := test

		t.Run(test.title, func(t *testing.T) {
			events := storage.GetEventsForWeek(test.week)
			require.Equal(t, test.events, len(events))
		})
	}
}

func TestStorageGetEventsForMonth(t *testing.T) {
	storage := setUpStorage()

	testCase := []struct {
		title  string
		week   time.Time
		events int
	}{
		{title: "get events from not existing date month", week: dateNotExisting, events: 1},
		{title: "get existing events at the week", week: dateExistsOnlyOne, events: 7},
		{title: "get a few similar events + event at the month", week: dateExistMany, events: 8},
	}

	for _, test := range testCase {
		test := test

		t.Run(test.title, func(t *testing.T) {
			events := storage.GetEventsForMonth(test.week)
			require.Equal(t, test.events, len(events))
		})
	}
}

func generateEvents(quantity int) []m.Event {
	events := make([]m.Event, 0, quantity)
	for i := 0; i < quantity; i++ {
		events = append(events, *m.New())
	}
	return events
}

func initStorage(events []m.Event) *Storage {
	storage := New()
	for _, ev := range events {
		err := storage.Add(ev)
		if err != nil {
			continue
		}
	}
	return storage
}

func setUpStorage() *Storage {
	storage := initStorage(generateEvents(len(dates)))
	index := 0
	for _, event := range storage.events {
		event.DateTime = dates[index]
		index++
	}
	return storage
}
