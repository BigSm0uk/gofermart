package migration

import "embed"

//go:embed accrual/*.sql
var AccrualMigrations embed.FS

//go:embed loyalty/*.sql
var LoyaltyMigrations embed.FS
