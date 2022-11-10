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
	requestTimeStart := time.Now()

	wr := NewResponseWrapper(w)
	m.handler.ServeHTTP(wr, r)

	methodPathProto := fmt.Sprintf("%s %s %s", r.Method, r.RequestURI, r.Proto)

	latency := time.Now().Sub(requestTimeStart)

	m.logger.Info(fmt.Sprintf("%s [%s] %s %v %v \"%s\"", r.RemoteAddr, time.Now(), methodPathProto,
		wr.statusCode, latency, r.UserAgent()))
}

func NewMiddleware(handler http.Handler, logger app.Logger) *Middleware {
	return &Middleware{handler: handler, logger: logger}
}

type ResponseWrapper struct {
	http.ResponseWriter
	statusCode int
}

func NewResponseWrapper(w http.ResponseWriter) *ResponseWrapper {
	return &ResponseWrapper{w, http.StatusOK}
}

func (wr *ResponseWrapper) WriteHeader(code int) {
	wr.statusCode = code
	wr.ResponseWriter.WriteHeader(code)
}
