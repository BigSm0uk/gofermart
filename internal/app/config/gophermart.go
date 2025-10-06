package config

import (
	"flag"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Gophermart struct {
	Env                  string        `env:"ENV" env-default:"development"`
	RunAddress           string        `env:"RUN_ADDRESS" env-default:":8080"`
	DatabaseURI          string        `env:"DATABASE_URI" env-default:"postgresql://postgres:postgres@localhost:5432/gophermart?sslmode=disable"`
	AccrualSystemAddress string        `env:"ACCRUAL_SYSTEM_ADDRESS" env-default:"http://localhost:8081"`
	JWTSecret            string        `env:"JWT_SECRET" env-default:"your-secret-key"`
	WorkerInterval       time.Duration `env:"WORKER_INTERVAL" env-default:"5s"`
}

func InitGophermartConfig() (*Gophermart, error) {
	var config Gophermart

	flag.String("a", config.RunAddress, "address and port to run the service")
	flag.String("d", config.DatabaseURI, "database URI")
	flag.String("r", config.AccrualSystemAddress, "accrual system address")
	flag.String("j", config.JWTSecret, "JWT secret key")
	flag.Duration("w", config.WorkerInterval, "worker interval")
	flag.Parse()

	if err := cleanenv.ReadEnv(&config); err != nil {
		return nil, err
	}
	return &config, nil
}