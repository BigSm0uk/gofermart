package accrual

import (
	"context"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

type App struct {
	Container *Container
}

func InitApp() (*App, error) {
	c, err := Build()
	if err != nil {
		return nil, err
	}
	return &App{Container: c}, nil
}
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	go a.Container.Processor.Run(ctx)

	go func() {
		if err := a.Container.Router.Listen(); err != nil {
			a.Container.Log.Error("failed to start server", zap.Error(err))
			cancel()
		}
	}()
	<-ctx.Done()
	a.Container.Log.Info("shutting down server")
	a.Container.Router.Shutdown()
	return nil
}
