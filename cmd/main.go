package main

import (
	"context"
	"go-recorder/app"
	"go-recorder/facility"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

func main() {
	logger := logrus.New()

	cfg, err := facility.NewConfig()
	if err != nil {
		logger.Panicf("failed to create config: %v", err)
	}

	startingCtx, cancel := context.WithCancel(context.Background())
	if err := app.Run(startingCtx, logger, cfg); err != nil {
		logger.Panicf("failed to run app: %v", err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, os.Interrupt)

	<-quit
	cancel()
}
