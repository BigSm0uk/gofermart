package accrual

import (
	"database/sql"

	"github.com/BigSm0uk/gofermart/migration"
	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migration.AccrualMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(db, ".")
}
