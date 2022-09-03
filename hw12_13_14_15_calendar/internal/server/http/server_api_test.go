package internalhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/logger"
	memorystorage "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/models"
	"github.com/stretchr/testify/require"
)

const (
	Host                         = "http://localhost:"
	PathHandlerAdd               = "/events/add"
	PathHandlerDelete            = "/events/delete"
	PathHandlerUpdate            = "/events/update"
	PathHandlerGetEventsForDay   = "/events/get-list-for-day"
	PathHandlerGetEventsForWeek  = "/events/get-list-for-week"
	PathHandlerGetEventsForMonth = "/events/get-list-for-month"

	QueryParamID   = "id"
	QueryParamDate = "date"
)

var (
	ctx     = context.Background()
	config  = configs.NewConfig("../../../configs/config_template.toml")
	log     = logger.New(config)
	storage app.Storage

	eventToAdd      = models.Event{ID: eventID, Title: "Mocked event"}
	eventToUpdate   = models.Event{ID: eventID, Title: "Brand new title", Description: "Updated description"}
	eventID         = "123456789"
	eventToDeleteID = "3"

	timeDefaultParam = "0001-01-01"

	eventsInMemory = []models.Event{
		{ID: "1", Title: "first"},
		{ID: "2", Title: "second"},
		{ID: "3", Title: "third"},
	}
)

func init() {
	storage = memorystorage.New()
	for _, ev := range eventsInMemory {
		_, err := storage.Add(ctx, ev)
		if err != nil {
			os.Stderr.WriteString(err.Error())
			os.Exit(1)
		}
	}

	calendar := app.New(log, storage)
	server := NewServer(calendar, storage, log, config)

	errs := make(chan error, 1)
	go func() {
		errs <- server.Start(ctx)
	}()

	go func() {
		select {
		case failed := <-errs:
			if failed != nil {
				os.Stderr.WriteString(failed.Error())
				os.Exit(1)
			}
		}
	}()
}

func TestHandlerAdd(t *testing.T) {
	eventsBefore := getEvents(t, PathHandlerGetEventsForMonth)

	path := Host + config.Server.Port + PathHandlerAdd
	req := createRequestPOST(t, path, eventToAdd)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal()
	}
	defer resp.Body.Close()

	bodyResp := getBodyFromResponse(t, resp.Body)

	eventsAfter := getEvents(t, PathHandlerGetEventsForMonth)

	require.Equal(t, len(eventsBefore)+1, len(eventsAfter))
	require.Equal(t, fmt.Sprintf("Created event with ID: %s", eventToAdd.ID), bodyResp)
}

func createRequestPOST(t *testing.T, path string, body models.Event) *http.Request {
	encodedEvent, err := json.Marshal(body)
	if err != nil {
		t.Fail()
	}

	req, err := http.NewRequest(http.MethodPost, path, bytes.NewBuffer(encodedEvent))
	if err != nil {
		t.Fatal()
	}
	return req
}

func getBodyFromResponse(t *testing.T, resp io.ReadCloser) string {
	b, err := io.ReadAll(resp)
	if err != nil {
		t.Fatal()
	}
	return string(b)
}

func TestHandlerUpdate(t *testing.T) {
	path := Host + config.Server.Port + PathHandlerUpdate
	req := createRequestPOST(t, path, eventToUpdate)

	setUpQueryParam(req, QueryParamID, eventID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal()
	}
	defer resp.Body.Close()

	bodyResp := getBodyFromResponse(t, resp.Body)

	require.Equal(t, fmt.Sprintf("Event with ID: %s was updated.", eventID), bodyResp)
}

func TestHandlerDelete(t *testing.T) {
	eventsBefore := getEvents(t, PathHandlerGetEventsForMonth)

	path := Host + config.Server.Port + PathHandlerDelete
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal()
	}

	setUpQueryParam(req, QueryParamID, eventToDeleteID)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal()
	}
	defer resp.Body.Close()

	eventsAfter := getEvents(t, PathHandlerGetEventsForMonth)

	require.Equal(t, len(eventsBefore)-1, len(eventsAfter))
}

func TestHandlerGetEventsForDay(t *testing.T) {
	events := getEvents(t, PathHandlerGetEventsForDay)
	require.Equal(t, len(eventsInMemory), len(events))
}

func TestHandlerGetEventsForWeek(t *testing.T) {
	events := getEvents(t, PathHandlerGetEventsForWeek)
	require.Equal(t, len(eventsInMemory), len(events))
}

func TestHandlerGetEventsForMonth(t *testing.T) {
	events := getEvents(t, PathHandlerGetEventsForMonth)
	require.Equal(t, len(eventsInMemory), len(events))
}

func getEvents(t *testing.T, handlerPath string) []models.Event {
	path := Host + config.Server.Port + handlerPath
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		t.Fatal()
	}

	setUpQueryParam(req, QueryParamDate, timeDefaultParam)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal()
	}
	defer resp.Body.Close()

	var events []models.Event
	err = json.NewDecoder(resp.Body).Decode(&events)
	if err != nil {
		t.Fatal()
	}

	return events
}

func setUpQueryParam(req *http.Request, key, value string) {
	params := req.URL.Query()
	params.Add(key, value)

	req.URL.RawQuery = params.Encode()
}
