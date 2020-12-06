package appInit

import (
	"context"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

type CommonApp interface {
	GetLogger() *zap.Logger
	GetErr() chan error
}

type AppConstructor =  func(ctx context.Context) (CommonApp, error)

func Initialize(constructor AppConstructor) {
	defer time.Sleep(1500 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sm, err := constructor(ctx)
	if err != nil {
		panic(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	logger := (*sm).GetLogger()

	select {
	case <-stop:
		logger.Info("Application was interrupted.")
	case err := <-(*sm).GetErr():
		logger.Panic("A fatal error occured", zap.Error(err))
	}
}
