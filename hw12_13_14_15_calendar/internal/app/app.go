package app

import (
	"context"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"time"

	m "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/models"
)

type App struct {
	storage Storage
	logger  Logger
	config  configs.Config
}

type Logger interface {
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type Storage interface {
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
	Add(ctx context.Context, event m.Event) (string, error)
	Update(ctx context.Context, eventId string, updatedEvent m.Event) error
	Delete(ctx context.Context, eventId string) error
	GetEventsForDay(ctx context.Context, day time.Time) []m.Event
	GetEventsForWeek(ctx context.Context, day time.Time) []m.Event
	GetEventsForMonth(ctx context.Context, day time.Time) []m.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{storage: storage, logger: logger}
}

func (a *App) CreateEvent(ctx context.Context, title string) error {
	newEvent := m.New()
	newEvent.Title = title

	_, err := a.storage.Add(ctx, *newEvent)
	return err
}
