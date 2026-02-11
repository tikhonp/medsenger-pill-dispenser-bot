package util

import (
	"github.com/tikhonp/maigo"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/config"
)

type Dependencies struct {
	Cfg   *config.Config
	Maigo *maigo.Client
	DB    db.ModelsFactory
}

func NewDependencies(cfg *config.Config, maigo *maigo.Client, db db.ModelsFactory) Dependencies {
	return Dependencies{
		Cfg:   cfg,
		Maigo: maigo,
		DB:    db,
	}
}
