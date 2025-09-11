package models

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type PillDispenser struct {
	SerialNumber  string        `db:"serial_number"`
	HWType        HardwareType  `db:"hw_type_id"`
	LastFetchTime sql.NullTime  `db:"last_fetch_time"`
	ContractID    sql.NullInt64 `db:"contract_id"`
	Description   string        `db:"description"`
}

var (
	ErrContractIDAlreadySet   = errors.New("contract id already set for this pill dispenser")
	ErrPillDispenserNotExists = errors.New("pill dispenser not exists for specified serial number")
)

type PillDispensers interface {
	Get(serialNumber string) (*PillDispenser, error)

	// New creates new pill dispenser with specific serial number and hardware type.
	New(serialNumber, description string, hwType HardwareType) error

	// GetContractID fetches contract id for pill dispenser with specific serial number
	GetContractID(serialNumber string) (int, error)

	// GetAllByContractID fetches all pill dispensers related to contract
	GetAllByContractID(contractID int) ([]PillDispenser, error)

	// RegisterContractID sets contract ID for pill-dispenser with specific serial number
	// returns ErrContractIdAlreadySet or ErrPillDispenserNotExists
	RegisterContractID(serialNumber string, contractID int) error

	// UnregisterContractID clears contract id for pill-dispenser with specific serial number
	UnregisterContractID(serialNumber string) error

	// UnregisterByContractID clears contract id for all pill-dispensers with specific set contract id,
	// useful for clearing up links with pill dispensers on contract deactivation.
	UnregisterByContractID(contractID int) error

	// UpdateLastFetchTime updates last fetch time to current time for pill-dispenser with specific serial number
	UpdateLastFetchTime(serialNumber string) error
}

type pillDispensers struct {
	db *sqlx.DB
}

func NewPillDispensers(db *sqlx.DB) PillDispensers {
	return &pillDispensers{db: db}
}

func (pd *pillDispensers) Get(serialNumber string) (*PillDispenser, error) {
	var pillDispenser PillDispenser
	err := pd.db.Get(&pillDispenser, "SELECT * FROM pill_dispenser WHERE serial_number = $1", serialNumber)
	return &pillDispenser, err
}

func (pd *pillDispensers) New(serialNumber, description string, hwType HardwareType) error {
	_, err := pd.db.Exec("INSERT INTO pill_dispenser(serial_number, hw_type_id, description) VALUES ($1, $2, $3)", serialNumber, hwType, description)
	return err
}

func (pd *pillDispensers) GetContractID(serialNumber string) (int, error) {
	var contractID int
	err := pd.db.QueryRow("SELECT contract_id FROM pill_dispenser WHERE serial_number = $1", serialNumber).Scan(&contractID)
	return contractID, err
}

func (pd *pillDispensers) GetAllByContractID(contractID int) ([]PillDispenser, error) {
	var pillDispensers []PillDispenser
	err := pd.db.Select(&pillDispensers, "SELECT * FROM pill_dispenser WHERE contract_id = $1", contractID)
	return pillDispensers, err
}

func (pd *pillDispensers) RegisterContractID(serialNumber string, contractID int) error {
	pillDispenser, err := pd.Get(serialNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPillDispenserNotExists
	}
	if err != nil {
		return err
	}
	if pillDispenser.ContractID.Valid {
		if pillDispenser.ContractID.Int64 == int64(contractID) {
			return nil
		} else {
			return ErrContractIDAlreadySet
		}
	}
	_, err = pd.db.Exec("UPDATE pill_dispenser SET contract_id = $1 WHERE serial_number = $2", contractID, serialNumber)
	return err
}

func (pd *pillDispensers) UnregisterContractID(serialNumber string) error {
	_, err := pd.db.Exec("UPDATE pill_dispenser SET contract_id = NULL WHERE serial_number = $1", serialNumber)
	return err
}

func (pd *pillDispensers) UnregisterByContractID(contractID int) error {
	_, err := pd.db.Exec("UPDATE pill_dispenser SET contract_id = NULL WHERE contract_id = $1", contractID)
	return err
}
func (pd *pillDispensers) UpdateLastFetchTime(serialNumber string) error {
	_, err := pd.db.Exec("UPDATE pill_dispenser SET last_fetch_time = NOW() WHERE serial_number = $1", serialNumber)
	if err != nil {
		return fmt.Errorf("failed to update last fetch time for pill dispenser with SN %s: %w", serialNumber, err)
	}
	return nil
}
