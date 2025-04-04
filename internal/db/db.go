package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
)

// Connect to the database and return a connection.
func Connect(cfg *config.Database) (ModelsFactory, error) {
    db, err := sqlx.Connect("sqlite3", cfg.DbFilePath)
	if err != nil {
		return nil, err
	}
	return newModelsFactory(db), nil
}
