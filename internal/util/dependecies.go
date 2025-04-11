package util

import (
	"github.com/TikhonP/maigo"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
)

type Dependencies struct {
	Cfg   *config.Config
	Maigo *maigo.Client
	Db    db.ModelsFactory
}

func NewDependencies(cfg *config.Config, maigo *maigo.Client, db db.ModelsFactory) Dependencies {
	return Dependencies{
		Cfg:   cfg,
		Maigo: maigo,
		Db:    db,
	}
}
