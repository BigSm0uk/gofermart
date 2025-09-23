package config

import (
	"flag"

	"github.com/ilyakaznacheev/cleanenv"
)

type Gophermart struct {
	Env                  string `env:"ENV" env-default:"development"`
	RunAddress           string `env:"RUN_ADDRESS" env-default:":8090"`
	DatabaseURI          string `env:"DATABASE_URI" env-default:"postgresql://admin:admin@localhost:5438/gophermart"`
	AccrualSystemAddress string `env:"ACCRUAL_SYSTEM_ADDRESS" env-default:"http://localhost:8080"`
}

func InitGophermartConfig() (*Gophermart, error) {
	var config Gophermart

	flag.String("a", config.RunAddress, "address and port to run the service")
	flag.String("d", config.DatabaseURI, "database URI")
	flag.String("r", config.AccrualSystemAddress, "accrual system address")
	flag.Parse()

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
