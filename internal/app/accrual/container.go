package accrual

import (
	"fmt"

	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/router"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/handlers"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/usecase"
	"go.uber.org/zap"
)

type Container struct {
	Config    *config.Accrual
	Router    *router.AccrualRouter
	Handler   *handlers.AccrualHandler
	Usecase   *usecase.AccrualUsecase
	Repo      *repo.AccrualRepository
	Processor *AccrualProcessor
	Log       *zap.Logger
}

func NewContainer() *Container {
	return &Container{}
}

func (c *Container) LoadConfig() *Container {
	config, err := config.InitAccrualConfig()
	if err != nil {
		panic(err)
	}
	c.Config = config
	return c
}
func (c *Container) LoadLogger() *Container {
	c.Log = zl.InitLogger(c.Config.Env)
	return c
}
func (c *Container) LoadRouter() *Container {
	c.Router = router.NewAccrualRouter(c.Config, c.Handler, c.Log)
	return c
}
func (c *Container) LoadHandler() *Container {
	c.Handler = handlers.NewAccrualHandler(c.Usecase, c.Log)
	db := c.Repo.StdDB()
	err := Migrate(db)
	if err != nil {
		c.Log.Fatal("failed to migrate database", zap.Error(err))
	}
	return c
}
func (c *Container) LoadUsecase() *Container {
	c.Usecase = usecase.NewAccrualUsecase(c.Repo)
	return c
}
func (c *Container) LoadRepo() *Container {
	c.Repo = repo.NewAccrualRepository(c.Config, c.Log)
	return c
}
func (c *Container) LoadProcessor() *Container {
	c.Processor = NewAccrualProcessor(c.Repo, c.Config.Processor.TTL, c.Config.Processor.Limit, c.Log)
	return c
}
func Build() (*Container, error) {
	c := NewContainer().LoadConfig().LoadLogger().LoadRepo().LoadUsecase().LoadHandler().LoadRouter().LoadProcessor()
	switch true {
	case c.Config == nil:
		return nil, fmt.Errorf("config is nil")
	case c.Log == nil:
		return nil, fmt.Errorf("logger is nil")
	case c.Router == nil:
		return nil, fmt.Errorf("router is nil")
	case c.Usecase == nil:
		return nil, fmt.Errorf("usecase is nil")
	case c.Handler == nil:
		return nil, fmt.Errorf("handler is nil")
	case c.Repo == nil:
		return nil, fmt.Errorf("repo is nil")
	case c.Processor == nil:
		return nil, fmt.Errorf("processor is nil")
	}
	return c, nil
}
