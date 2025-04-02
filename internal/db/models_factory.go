package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
)

type ModelsFactory interface {
	Contracts() models.Contracts
}

type modelsFactory struct {
	contracts models.Contracts
}

func newModelsFactory(db *sqlx.DB) ModelsFactory {
	return &modelsFactory{
		contracts: models.NewContracts(db),
	}
}

func (f *modelsFactory) Contracts() models.Contracts {
	return f.contracts
}
