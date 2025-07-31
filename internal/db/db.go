// Package db provides functionality to connect to a database and manage models.
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/config"
)

func DataSourceName(cfg *config.Database) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Dbname)
}

// Connect to the database and return a connection.
func Connect(cfg *config.Database) (ModelsFactory, error) {
	db, err := sqlx.Connect("postgres", DataSourceName(cfg))
	if err != nil {
		return nil, err
	}
	return newModelsFactory(db), nil
}
