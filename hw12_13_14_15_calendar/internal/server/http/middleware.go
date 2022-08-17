package internalhttp

import (
	"fmt"
	"net/http"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
)

type Middleware struct {
	handler http.Handler
	logger  app.Logger
}

func (m *Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.handler.ServeHTTP(w, r)

	methodPathProto := fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)

	statusCode := 0
	if r.Response != nil { // TODO: handle the status code correctly
		statusCode = r.Response.StatusCode
	}

	latency := 0 // TODO: count request.timeNow - response.timeNow

	m.logger.Info(fmt.Sprintf("%s [%s] %s %v %v \"%s\"", r.RemoteAddr, time.Now(), methodPathProto,
		statusCode, latency, r.UserAgent()))
}

func NewMiddleware(handler http.Handler, logger app.Logger) *Middleware {
	return &Middleware{handler: handler, logger: logger}
}
