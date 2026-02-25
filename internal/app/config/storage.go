package config

import "time"

type AccrualStorage struct {
	DatabaseURI        string        `env:"DATABASE_URI" env-default:"postgresql://admin:admin@localhost:5437/accrual?sslmode=disable"`
	MinPoolSize        int32         `env:"MIN_POOL_SIZE" env-default:"1"`
	MaxPoolSize        int32         `env:"MAX_POOL_SIZE" env-default:"10"`
	ConnectionLifetime time.Duration `env:"CONNECTION_LIFETIME" env-default:"10s"`
}
type GophermartStorage struct {
	DatabaseURI string `env:"DATABASE_URI" env-default:"postgresql://admin:admin@localhost:5432/gophermart?sslmode=disable"`
}
