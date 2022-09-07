package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/configs"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/internal/sender"
	"github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/ampq"
	log "github.com/novikov-ai/practice-misc/hw12_13_14_15_calendar/pkg/logger"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_scheduler.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config := configs.NewConfigSender(configFile)

	logger := log.New(&config)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	messageBroker := ampq.New(&config, logger)
	err := messageBroker.Connect()
	if err != nil {
		exitWithError(logger, err)
	}
	defer messageBroker.Close()

	sender := sender.New(logger, config, messageBroker)

	errs := make(chan error, 1)

	logger.Info("sender is running...")

	go func() {
		errs <- sender.Start(ctx, os.Stdout)
	}()

	select {
	case <-ctx.Done():
		logger.Info("sender stopped worked by cancellation")
		break
	case err := <-errs:
		exitWithError(logger, err)
	}
}

func exitWithError(log *log.Logger, err error) {
	log.Error(err.Error())
	os.Exit(1)
}
