package gophermart

import (
	"database/sql"

	"github.com/BigSm0uk/gofermart/migration"
	"github.com/pressly/goose/v3"
)

func Migrate(db *sql.DB) error {
	goose.SetBaseFS(migration.LoyaltyMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	// Сначала сбрасываем миграции если они есть, затем применяем заново
	_ = goose.Reset(db, "loyalty") // Игнорируем ошибку если таблица миграций не существует

	// Применяем миграции из папки loyalty
	return goose.Up(db, "loyalty")
}
