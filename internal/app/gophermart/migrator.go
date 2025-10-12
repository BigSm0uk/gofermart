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

	if err := goose.Up(db, "loyalty"); err != nil {
		panic(err)
	}
}

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migration.LoyaltyMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Сначала сбрасываем миграции если они есть, затем применяем заново
	if err := goose.Reset(db, "loyalty"); err != nil {
		// Игнорируем ошибку если таблица миграций не существует
	}
	
	// Применяем миграции из папки loyalty
	return goose.Up(db, "loyalty")
}
