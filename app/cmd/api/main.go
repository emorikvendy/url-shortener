package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/emorikvendy/url-shortener/internal/api"
	"github.com/emorikvendy/url-shortener/internal/diagnostics"
	"github.com/emorikvendy/url-shortener/internal/resources"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()

	//nolint:errcheck Не может быть ошибки, т.к. работаем с stdout
	defer logger.Sync()

	slogger := logger.Sugar()
	slogger.Info("Starting the application...")
	slogger.Info("Reading configuration and initializing resources...")

	rsc, err := resources.New(slogger)
	if err != nil {
		slogger.Panic("Can't initialize resources.", "err", err)
	}
	defer func() {
		err = rsc.Release()
		if err != nil {
			slogger.Errorw("Got an error during resources release.", "err", err)
		}
	}()

	slogger.Info("Configuring the application units...")

	slogger.Info("Starting the servers...")

	restAPI := api.New(slogger, rsc.Config.RESTAPIPort, *rsc)
	restAPI.Start()

	diag := diagnostics.New(slogger, rsc.Config.DiagPort, rsc.Healthz)
	diag.Start()
	slogger.Info("The application is ready to serve requests.")

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	select {
	case x := <-interrupt:
		slogger.Infow("Received a signal.", "signal", x.String())
	case errRest := <-restAPI.Notify():
		slogger.Errorw("Received an error from the restAPI server.", "err", errRest)
	case errDiag := <-diag.Notify():
		slogger.Errorw("Received an error from the diagnostics server.", "err", errDiag)
	}

	slogger.Info("Stopping the servers...")

	err = restAPI.Stop()
	if err != nil {
		slogger.Error("Got an error while stopping the restAPI logic server.", "err", err)
	}

	err = diag.Stop()
	if err != nil {
		slogger.Error("Got an error while stopping the diag logic server.", "err", err)
	}

	slogger.Info("The app is calling the last defers and will be stopped.")
}
