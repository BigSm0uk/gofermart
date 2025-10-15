package gophermart

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

func (a *App) Run() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	// Запускаем воркер в отдельной горутине
	go a.Container.Worker.Run(ctx)

	// Запускаем HTTP сервер
	go func() {
		a.Container.Log.Info("Starting server", zap.String("address", a.Container.Config.RunAddress))
		if err := a.Container.Server.Listen(a.Container.Config.RunAddress); err != nil {
			a.Container.Log.Error("Failed to start server", zap.Error(err))
			cancel()
		}
	}()

	// Ожидаем сигнал завершения
	<-ctx.Done()

	a.Container.Log.Info("Shutting down server...")

	// Останавливаем сервер
	if err := a.Container.Server.Shutdown(); err != nil {
		a.Container.Log.Error("Server forced to shutdown", zap.Error(err))
	}

	// Закрываем соединения
	a.Container.Close()

	a.Container.Log.Info("Server exited")
}
