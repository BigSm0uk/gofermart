package config

type Storage struct {
	DatabaseURI string `env:"DATABASE_URI" env-default:"User ID=admin;Password=admin;Host=localhost;Port=5437;Database=accrual;Pooling=true;Min Pool Size=1;Max Pool Size=10;Connection Lifetime=10s;"`
}
