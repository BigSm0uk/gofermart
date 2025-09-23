package accrual

import (
	"github.com/BigSm0uk/gofermart/internal/app/config"
	"github.com/BigSm0uk/gofermart/internal/app/zl"
	"github.com/gofiber/fiber/v3"
)

type Container struct {
	Config *config.Accrual
	Server *fiber.App
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
func (c *Container) LoadServer() *Container {
	c.Server = fiber.New()
	return c
}
func MustBuild() *Container {
	return NewContainer().LoadConfig().LoadLogger().LoadServer()
}
