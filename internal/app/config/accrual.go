package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Accrual struct {
	Env         string `env:"ENV" env-default:"development"`
	RunAddress  string `env:"RUN_ADDRESS" env-default:":3000"`
	DatabaseURI string `env:"DATABASE_URI" env-default:"postgresql://admin:admin@localhost:5437/accrual"`
}

func InitAccrualConfig() (*Accrual, error) {
	var config Accrual
	flag.StringVar(&config.RunAddress, "a", config.RunAddress, "address and port to run the service")
	flag.StringVar(&config.DatabaseURI, "d", config.DatabaseURI, "database URI")
	flag.Parse()
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
