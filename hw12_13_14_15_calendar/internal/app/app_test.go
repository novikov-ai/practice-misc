package app

import (
	"context"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app/mocks"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestApp(t *testing.T) {
	logger := mocks.NewLogger(t)
	storage := mocks.NewStorage(t)

	app := New(logger, storage)

	const (
		eventTitle = "Mocking"
		eventID    = "12345678987654321"
	)

	newEvent := models.Event{ID: eventID, Title: eventTitle}

	ctx := context.Background()

	call := storage.On("Add", ctx, newEvent).Return(eventID, nil)
	call.Once()

	dummyGenID := func() (string, error) {
		return eventID, nil
	}

	err := app.CreateEvent(ctx, eventTitle, dummyGenID)
	require.Nil(t, err)

	require.True(t, storage.AssertExpectations(t))
}
