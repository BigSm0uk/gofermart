package accrual

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"go.uber.org/zap"
)

type App struct {
	Container *Container
}

func InitApp() (*App, error) {
	c := MustBuild()
	return &App{Container: c}, nil
}
func (a *App) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()
	go func() {
		if err := a.Container.Router.Listen(); err != nil {
			zl.Log.Error("failed to start server", zap.Error(err))
			cancel()
		}
	}()
	<-ctx.Done()
	zl.Log.Info("shutting down server")
	a.Container.Router.Shutdown()
}
