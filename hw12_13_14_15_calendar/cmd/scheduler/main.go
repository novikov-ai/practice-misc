package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	sch "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/scheduler"
	sqlstorage "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/storage/sql"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/ampq"
	log "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_scheduler.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := configs.NewConfigScheduler(configFile)

	logger := log.New(&config)
	storage := sqlstorage.New(&config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	err := storage.Connect(ctx)
	if err != nil {
		logger.Error(err.Error())
		return
	}

	messageBroker := ampq.New(config, logger)
	err = messageBroker.Connect()
	if err != nil {
		exitWithError(logger, err)
	}
	defer messageBroker.Close()

	scheduler := sch.New(storage, logger, config, messageBroker)

	errs := make(chan error, 1)

	go func() {
		errs <- scheduler.Scan(ctx)
	}()

	select {
	case <-ctx.Done():
		logger.Info("scheduler stopped worked by cancellation")
		break
	case err := <-errs:
		exitWithError(logger, err)
	}
}

func exitWithError(log *log.Logger, err error) {
	log.Error(err.Error())
	os.Exit(1)
}
