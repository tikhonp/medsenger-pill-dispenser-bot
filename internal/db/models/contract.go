package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Contract represents Medsenger contract.
// Create on agent /init and persist during agent lifecycle.
type Contract struct {
	Id           int            `db:"id"`
	IsActive     bool           `db:"is_active"`
	AgentToken   sql.NullString `db:"agent_token"`
	PatientName  sql.NullString `db:"patient_name"`
	PatientEmail sql.NullString `db:"patient_email"`
	Locale       sql.NullString `db:"locale"`
}

type Contracts interface {

	// GetActiveContractIds returns all active contracts ids.
	// Use it for medsenger status endpoint.
	GetActiveContractIds() ([]int, error)
}

type contracts struct {
	db *sqlx.DB
}

func NewContracts(db *sqlx.DB) Contracts {
	return &contracts{db: db}
}

func (c *contracts) GetActiveContractIds() ([]int, error) {
	var contractIds = make([]int, 0)
	err := c.db.Select(&contractIds, `SELECT id FROM contracts WHERE is_active = true`)
	return contractIds, err
}
