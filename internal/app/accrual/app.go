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
func (a *App) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	err := a.preRunActions(ctx)
	if err != nil {
		return err
	}

	go a.Container.Processor.Run(ctx)

	go func() {
		if err := a.Container.Router.Listen(); err != nil {
			zl.Log.Error("failed to start server", zap.Error(err))
			cancel()
		}
	}()
	<-ctx.Done()
	zl.Log.Info("shutting down server")
	a.Container.Router.Shutdown()
	return nil
}
func (a *App) preRunActions(ctx context.Context) error {
	err := a.Container.Usecase.Ping(ctx)
	if err != nil {
		return err
	}
	return a.migrateWithDB()
}
func (a *App) migrateWithDB() error {
	db := a.Container.Repo.StdDB()
	return Migrate(db)
}
