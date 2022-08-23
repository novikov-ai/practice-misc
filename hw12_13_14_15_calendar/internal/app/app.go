package app

import (
	"context"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/configs"
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
	Add(event m.Event) (string, error)
	Update(eventId string, updatedEvent m.Event) error
	Delete(eventId string) error
	GetEventsForDay(day time.Time) []m.Event
	GetEventsForWeek(day time.Time) []m.Event
	GetEventsForMonth(day time.Time) []m.Event
}

func New(logger Logger, storage Storage) *App {
	return &App{storage: storage, logger: logger}
}

func (a *App) CreateEvent(title string) error {
	newEvent := m.New()
	newEvent.Title = title

	_, err := a.storage.Add(*newEvent)
	return err
}
