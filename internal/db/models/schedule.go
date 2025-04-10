package models

import (
	"database/sql"
	"time"
)

type Schedule struct {
	IsOfflineNotificationsAllowed bool          `db:"is_offline_notifications_allowed"`
	RefreshRateInterval           sql.NullInt64 `db:"refresh_rate_interval"`
	ContractID                    int           `db:"contract_id"`
	PillDispenserSN               string        `db:"pill_dispenser_sn"`
	CreatedAt                     time.Time     `db:"created_at"`
}

type Schedules interface{}
