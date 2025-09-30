package config

import "time"

type Storage struct {
	DatabaseURI        string        `env:"DATABASE_URI" env-default:"postgresql://admin:admin@localhost:5437/accrual"`
	MinPoolSize        int32         `env:"MIN_POOL_SIZE" env-default:"1"`
	MaxPoolSize        int32         `env:"MAX_POOL_SIZE" env-default:"10"`
	ConnectionLifetime time.Duration `env:"CONNECTION_LIFETIME" env-default:"10s"`
}
