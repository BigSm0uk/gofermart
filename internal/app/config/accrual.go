package config

import (
	"flag"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Accrual struct {
	Env        string `env:"ENV" env-default:"development"`
	RunAddress string `env:"RUN_ADDRESS" env-default:":3000"`
	Storage    AccrualStorage
	Processor  AccrualProcessor
}

func InitAccrualConfig() (*Accrual, error) {
	var config Accrual
	flag.StringVar(&config.RunAddress, "a", "", "address and port to run the service")
	flag.StringVar(&config.Storage.DatabaseURI, "d", "", "database URI")
	flag.Parse()
	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

type AccrualProcessor struct {
	TTL   time.Duration `env:"TTL" env-default:"10s"`
	Limit uint64        `env:"LIMIT" env-default:"10"`
}
