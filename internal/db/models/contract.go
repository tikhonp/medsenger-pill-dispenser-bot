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

	// NewContract creates new contract from Medsenger Core request.
	NewContract(contractId int, agentToken, locale string) error

	// UpdateContractWithPatientData saves contract patient meta data to db
	UpdateContractWithPatientData(contractId int, patientName, patientEmail string) error

	// MarkInactiveContractWithId sets contract with id to inactive.
	// Use it for medsenger remove endpoint.
	// Equivalent to DELETE FROM contracts WHERE id = ?.
	MarkInactiveContractWithId(id int) error
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

func (c *contracts) NewContract(contractId int, agentToken, locale string) error {
	const query = `
		INSERT INTO contracts (id, is_active, agent_token, locale)
        VALUES ($1, TRUE, $2, $3) ON CONFLICT (id)
		DO UPDATE SET is_active = EXCLUDED.is_active, agent_token = EXCLUDED.agent_token, locale = EXCLUDED.locale
	`
	_, err := c.db.Exec(query, contractId, agentToken, locale)
	return err

}

func (c *contracts) UpdateContractWithPatientData(contractId int, patientName, patientEmail string) error {
	const query = `
        UPDATE contracts SET (patient_name = $1, patient_email = $2) WHERE id = $3
    `
	_, err := c.db.Exec(query, patientName, patientEmail, contractId)
	return err
}

func (c *contracts) MarkInactiveContractWithId(id int) error {
	_, err := c.db.Exec(`UPDATE contracts SET is_active = false WHERE id = $1`, id)
	return err
}
