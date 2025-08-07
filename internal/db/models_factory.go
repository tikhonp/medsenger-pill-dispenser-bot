package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
)

type ModelsFactory interface {
	Contracts() models.Contracts
	PillDispensers() models.PillDispensers
	Schedules() models.Schedules
	BatteryStatuses() models.BatteryStatuses
}

type modelsFactory struct {
	contracts       models.Contracts
	pillDispensers  models.PillDispensers
	schedules       models.Schedules
	batteryStatuses models.BatteryStatuses
}

func newModelsFactory(db *sqlx.DB) ModelsFactory {
	return &modelsFactory{
		contracts:       models.NewContracts(db),
		pillDispensers:  models.NewPillDispensers(db),
		schedules:       models.NewSchedules(db),
		batteryStatuses: models.NewBatteryStatuses(db),
	}
}

func (f *modelsFactory) Contracts() models.Contracts {
	return f.contracts
}

func (f *modelsFactory) PillDispensers() models.PillDispensers {
	return f.pillDispensers
}

func (f *modelsFactory) Schedules() models.Schedules {
	return f.schedules
}

func (f *modelsFactory) BatteryStatuses() models.BatteryStatuses {
	return f.batteryStatuses
}
