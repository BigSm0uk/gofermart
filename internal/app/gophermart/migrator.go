package gophermart

import (
	"database/sql"

	"github.com/BigSm0uk/gofermart/migration"
	"github.com/pressly/goose/v3"
)

func MustMigrate(db *sql.DB) {
	goose.SetBaseFS(migration.LoyaltyMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "."); err != nil {
		panic(err)
	}
}
