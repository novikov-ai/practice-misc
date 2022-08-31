package main

import (
	"context"
	"flag"
	pb "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/grpc"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/app"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/configs"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/server/http"
	memorystorage "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "internal/configs/config_template.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := configs.NewConfig(configFile)
	logg := logger.New(config)

	var storage app.Storage
	if config.Database.InMemory {
		storage = memorystorage.New()
	} else {
		storage = sqlstorage.New(config)
	}

	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar, logg, config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	// starts protobuf server
	go func() {
		if err := pb.Start(ctx, storage, logg); err != nil {
			logg.Error("failed to start protobuf server: " + err.Error())
		}
	}()

	if err := server.Start(ctx); err != nil {
		logg.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
