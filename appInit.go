package appInit

import (
	"context"
	"github.com/skeris/appInit/version"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"time"
)

type CommonApp interface {
	GetLogger() *zap.Logger
	GetErr() chan error
}

type AppConstructor = func(context.Context, interface{}, Version) (CommonApp, error)

type Version struct {
	Release, Commit, BuildTime string
}

func Initialize(constructor AppConstructor, opts interface{}) {
	defer time.Sleep(1500 * time.Millisecond)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sm, err := constructor(ctx, getEnv(opts), Version{
		Release:   version.Release,
		Commit:    version.Commit,
		BuildTime: version.BuildTime,
	})
	if err != nil {
		panic(err)
	}
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	logger := sm.GetLogger()

	select {
	case <-stop:
		logger.Info("Application was interrupted.")
	case err := <-sm.GetErr():
		logger.Panic("A fatal error occured", zap.Error(err))
	}
}
