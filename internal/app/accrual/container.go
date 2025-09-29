package accrual

import (
	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/router"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/BigSm0uk/gofermart/internal/handlers"
	"github.com/BigSm0uk/gofermart/internal/repo"
	"github.com/BigSm0uk/gofermart/internal/usecase"
)

type Container struct {
	Config  *config.Accrual
	Router  *router.AccrualRouter
	Handler *handlers.AccrualHandler
	Usecase *usecase.AccrualUsecase
	Repo    *repo.AccrualRepository
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
	zl.InitLogger(c.Config.Env)
	return c
}
func (c *Container) LoadRouter() *Container {
	c.Router = router.NewAccrualRouter(c.Config, c.Handler)
	return c
}
func (c *Container) LoadHandler() *Container {
	c.Handler = handlers.NewAccrualHandler(c.Usecase)
	return c
}
func (c *Container) LoadUsecase() *Container {
	c.Usecase = usecase.NewAccrualUsecase(c.Repo)
	return c
}
func (c *Container) LoadRepo() *Container {
	c.Repo = repo.NewAccrualRepository(c.Config)
	return c
}
func MustBuild() *Container {
	c := NewContainer().LoadConfig().LoadLogger().LoadRepo().LoadUsecase().LoadHandler().LoadRouter()
	switch true {
	case c.Config == nil:
		panic("config is nil")
	case zl.Log == nil:
		panic("logger is nil")
	case c.Router == nil:
		panic("router is nil")
	case c.Usecase == nil:
		panic("usecase is nil")
	case c.Handler == nil:
		panic("handler is nil")
	case c.Repo == nil:
		panic("repo is nil")
	}
	return c
}
