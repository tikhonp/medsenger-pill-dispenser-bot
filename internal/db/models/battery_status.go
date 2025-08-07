package models

import (
	"time"

	"github.com/jmoiron/sqlx"
)

type BatteryStatus struct {
	ID           int    `db:"id"`
	SerialNumber string `db:"serial_nu"`
	// Voltage is the battery voltage in millivolts
	// e.g. 3000 for 3.0V
	Voltage   int       `db:"voltage"`
	CreatedAt time.Time `db:"created_at"`
}

type BatteryStatuses interface {
	// InsertBatteryStatus inserts a new battery status record into the database
	InsertBatteryStatus(status BatteryStatus) (*BatteryStatus, error)
}

type batteryStatus struct {
	db *sqlx.DB
}

func NewBatteryStatuses(db *sqlx.DB) BatteryStatuses {
	return &batteryStatus{
		db: db,
	}
}

func (b *batteryStatus) InsertBatteryStatus(status BatteryStatus) (*BatteryStatus, error) {
	query := `INSERT INTO battery_status (serial_nu, voltage, created_at) 
	          VALUES (:serial_nu, :voltage, :created_at) 
	          RETURNING id, serial_nu, voltage, created_at`
	q, args, err := sqlx.Named(query, status)
	if err != nil {
		return &status, err
	}
	q = b.db.Rebind(q)
	if err := b.db.QueryRowx(q, args...).StructScan(&status); err != nil {
		return nil, err
	}
	return &status, nil
}

