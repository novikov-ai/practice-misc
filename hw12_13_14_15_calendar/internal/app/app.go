package app

import (
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
	Add(event m.Event) error
	Update(eventId string, updatedEvent m.Event)
	Delete(eventId string)
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

	return a.storage.Add(*newEvent)
}
