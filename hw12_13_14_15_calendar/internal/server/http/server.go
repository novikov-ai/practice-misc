package internalhttp

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
)

type Server struct {
	application *app.App
	logger      app.Logger
	host, port  string
}

func NewServer(app *app.App, logger app.Logger, config configs.Config) *Server {
	return &Server{
		application: app, logger: logger,
		host: config.Server.Host, port: config.Server.Port,
	}
}

func (s *Server) Start(ctx context.Context) error {
	serverError := make(chan error)

	mux := http.NewServeMux()
	mux.HandleFunc("/welcome", handlerWelcome)
	wrappedMux := NewMiddleware(mux, s.logger)

	go func() {
		address := net.JoinHostPort(s.host, s.port)
		serverError <- http.ListenAndServe(address, wrappedMux)
	}()

	select {
	case err := <-serverError:
		return err
	case <-ctx.Done():
		return nil
	}
}

func (s *Server) Stop(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

func handlerWelcome(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	body := "Hello there!\nThe current time is"
	fmt.Fprintf(w, "<h1 align=\"center\">%s</h1>", body)
	fmt.Fprintf(w, "<h2 align=\"center\">%s</h2>", t)
}
