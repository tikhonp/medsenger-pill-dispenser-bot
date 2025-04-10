package models

import (
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Contract represents Medsenger contract.
// Create on agent /init and persist during agent lifecycle.
type Contract struct {
	ID                  int            `db:"id"`
	IsActive            bool           `db:"is_active"`
	ClinicID            int            `db:"clinic_id"`
	AgentToken          string         `db:"agent_token"`
	PatientAgentToken   string         `db:"patient_agent_token"`
	DoctorAgentToken    string         `db:"doctor_agent_token"`
	Locale              string         `db:"locale"`
	PatientName         sql.NullString `db:"patient_name"`
	PatientEmail        sql.NullString `db:"patient_email"`
	ScheduleID          sql.NullString `db:"schedule_id"`
	RefreshRateInterval sql.NullInt64  `db:"refresh_rate_interval"`
}

type Contracts interface {

	// GetActiveContractIds returns all active contracts ids.
	// Use it for medsenger status endpoint.
	GetActiveContractIds() ([]int, error)

	// NewContract creates new contract from Medsenger Core request.
	NewContract(contractId, clinicId int, agentToken, patientAgentToken, doctorAgentToken, locale string) error

	// UpdateContractWithPatientData saves contract patient meta data to db
	UpdateContractWithPatientData(contractId int, patientName, patientEmail string) error

	// MarkInactiveContractWithId sets contract with id to inactive.
	// Use it for medsenger remove endpoint.
	// Equivalent to DELETE FROM contracts WHERE id = ?.
	MarkInactiveContractWithId(id int) error

	// GetContract get contract data by id
	Get(id int) (*Contract, error)

	// GetByAgentToken get contract by agent token
	GetByAgentToken(agentToken string) (*Contract, error)
}

type contracts struct {
	db *sqlx.DB
}

func NewContracts(db *sqlx.DB) Contracts {
	return &contracts{db: db}
}

func (c *contracts) GetActiveContractIds() ([]int, error) {
	var contractIds []int
	err := c.db.Select(&contractIds, `SELECT id FROM contract WHERE is_active = true`)
	return contractIds, err
}

func (c *contracts) NewContract(contractId, clinicId int, agentToken, patientAgentToken, doctorAgentToken, locale string) error {
	const query = `
		INSERT INTO contract (id, is_active, clinic_id, agent_token, patient_agent_token, doctor_agent_token, locale)
        VALUES ($1, TRUE, $2, $3, $4, $5, $6) ON CONFLICT (id)
		DO UPDATE SET is_active = EXCLUDED.is_active, clinic_id = EXCLUDED.clinic_id, agent_token = EXCLUDED.agent_token, patient_agent_token = EXCLUDED.patient_agent_token,
        doctor_agent_token = EXCLUDED.doctor_agent_token, locale = EXCLUDED.locale
	`
	_, err := c.db.Exec(query, contractId, clinicId, agentToken, patientAgentToken, doctorAgentToken, locale)
	return err
}

func (c *contracts) UpdateContractWithPatientData(contractId int, patientName, patientEmail string) error {
	const query = `
        UPDATE contract SET (patient_name = $1, patient_email = $2) WHERE id = $3
    `
	_, err := c.db.Exec(query, patientName, patientEmail, contractId)
	return err
}

func (c *contracts) MarkInactiveContractWithId(id int) error {
	_, err := c.db.Exec(`UPDATE contract SET is_active = false WHERE id = $1`, id)
	return err
}

func (c *contracts) Get(id int) (*Contract, error) {
	var contract Contract
	err := c.db.Get(&contract, "SELECT * FROM contract WHERE id = $1", id)
	return &contract, err
}

func (c *contracts) GetByAgentToken(agentToken string) (*Contract, error) {
	var contract Contract
	err := c.db.Get(&contract, "SELECT * FROM contract WHERE agent_token = $1", agentToken)
	return &contract, err
}
