package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/logger"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/grpc"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	internalhttp "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_calendar.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := configs.NewConfig(configFile)
	logg := logger.New(&config)

	var storage app.Storage
	if config.Database.InMemory {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New(&config)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	err := storage.Connect(ctx)
	if err != nil {
		logg.Error(err.Error())
		return
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar, storage, logg, &config)

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	go func() {
		if err := grpc.Start(ctx, storage, logg, config); err != nil {
			logg.Error("failed to start protobuf server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
