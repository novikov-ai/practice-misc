package internalhttp

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/models"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
)

type ServiceAPI struct {
	ctx     context.Context
	logger  app.Logger
	storage app.Storage
}

func NewService(ctx context.Context, st app.Storage, l app.Logger) *ServiceAPI {
	return &ServiceAPI{ctx: ctx, storage: st, logger: l}
}

func (s *ServiceAPI) handlerAdd(w http.ResponseWriter, r *http.Request) {
	event, err := getEventFromBody(r)
	if err != nil {
		s.logger.Warn("provide body for adding new event")
		return
	}

	newID, err := s.storage.Add(s.ctx, event)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "Created event with ID: %s", newID)
}

func (s *ServiceAPI) handlerUpdate(w http.ResponseWriter, r *http.Request) {
	eventID, err := getArgumentFromQuery(r, "id")
	if err != nil {
		s.logger.Warn("provide \"id\" argument")
		return
	}

	event, err := getEventFromBody(r)
	if err != nil {
		s.logger.Warn("provide body for updating the event")
		return
	}

	err = s.storage.Update(s.ctx, eventID, event)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "Event with ID: %s was updated.", eventID)
}

func (s *ServiceAPI) handlerDelete(w http.ResponseWriter, r *http.Request) {
	eventID, err := getArgumentFromQuery(r, "id")
	if err != nil {
		s.logger.Warn("provide \"id\" argument")
		return
	}

	err = s.storage.Delete(s.ctx, eventID)
	if err != nil {
		return
	}

	fmt.Fprintf(w, "Event with ID: %s was deleted.", eventID)
}

func (s *ServiceAPI) handlerGetEventsForDay(w http.ResponseWriter, r *http.Request) {
	parsedTime, err := parseTimeFromQueryArgument(r)
	if err != nil {
		s.logger.Warn(err.Error())
		return
	}

	events := s.storage.GetEventsForDay(s.ctx, parsedTime)

	err = marshalResponse(w, events)
	if err != nil {
		s.logger.Warn(err.Error())
	}
}

func (s *ServiceAPI) handlerGetEventsForWeek(w http.ResponseWriter, r *http.Request) {
	parsedTime, err := parseTimeFromQueryArgument(r)
	if err != nil {
		s.logger.Warn(err.Error())
		return
	}

	events := s.storage.GetEventsForWeek(s.ctx, parsedTime)

	err = marshalResponse(w, events)
	if err != nil {
		s.logger.Warn(err.Error())
	}
}

func (s *ServiceAPI) handlerGetEventsForMonth(w http.ResponseWriter, r *http.Request) {
	parsedTime, err := parseTimeFromQueryArgument(r)
	if err != nil {
		s.logger.Warn(err.Error())
		return
	}

	events := s.storage.GetEventsForMonth(s.ctx, parsedTime)

	err = marshalResponse(w, events)
	if err != nil {
		s.logger.Warn(err.Error())
	}
}

func parseTimeFromQueryArgument(r *http.Request) (time.Time, error) {
	queryLines := r.URL.Query()
	if len(queryLines) == 0 {
		return time.Time{}, errors.New("empty query")
	}

	const Argument = "date"
	value, err := getArgumentFromQuery(r, Argument)
	if err != nil {
		return time.Time{}, fmt.Errorf("argument \"%s\" wasn't found", Argument)
	}

	parsedTime, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, errors.New("error during time string parsing")
	}

	return parsedTime, nil
}

func getArgumentFromQuery(r *http.Request, argument string) (string, error) {
	queryLines := r.URL.Query()
	if len(queryLines) == 0 {
		return "", errors.New("empty query")
	}

	value := queryLines.Get(argument)
	if value == "" {
		return "", fmt.Errorf("argument \"%s\" wasn't found", argument)
	}

	return value, nil
}

func marshalResponse(w http.ResponseWriter, events []models.Event) error {
	encoded, err := json.Marshal(events)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "%s", encoded)
	return nil
}

func getEventFromBody(r *http.Request) (models.Event, error) {
	var event models.Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		return event, err
	}

	return event, nil
}
