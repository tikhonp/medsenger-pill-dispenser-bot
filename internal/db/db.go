// Package db provides functionality to connect to a database and manage models.
package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/config"
)

func DataSourceName(cfg *config.Database) string {
	return fmt.Sprintf("user=%s dbname=%s sslmode=disable password=%s host=%s", cfg.User, cfg.Dbname, cfg.Password, cfg.Host)
}

// Connect to the database and return a connection.
func Connect(cfg *config.Database) (ModelsFactory, error) {
	db, err := sqlx.Connect("postgres", DataSourceName(cfg))
	if err != nil {
		return nil, err
	}
	return newModelsFactory(db), nil
}
