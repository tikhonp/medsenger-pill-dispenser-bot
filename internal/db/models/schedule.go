package models

import (
	"database/sql"
	"errors"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
)

var (
	ErrNoSchedule = errors.New("no schedule set for specified serial number")
)

const DefaultRefreshRateInterval int64 = 60 * 60 // seconds

type Schedule struct {
	ID                            int            `db:"id"`
	IsOfflineNotificationsAllowed bool           `db:"is_offline_notifications_allowed"`
	RefreshRateInterval           sql.NullInt64  `db:"refresh_rate_interval"`
	ContractID                    int            `db:"contract_id"`
	PillDispenserSN               sql.NullString `db:"pill_dispenser_sn"`
	CreatedAt                     time.Time      `db:"created_at"`
}

type Schedules interface {
	GetLastScheduleForContractID(contractID int) (*ScheduleData, error)
	GetSchedules(pillDispenserSN string, contractID int) ([]ScheduleData, error)
	// GetScheduleForSN fetches schedule for pill dispenser with specified serial number
	//  returns ErrNoSchedule
	GetScheduleForSN(serialNumber string) (*ScheduleData, error)
	NewSchedule(schedule ScheduleData) (*ScheduleData, error)
	EditSchedule(schedule ScheduleData) (*ScheduleData, error)
	GetPillNameAndContractID(serialNumber string, cellIndex int) (string, int, error)
}

type schedule struct {
	db *sqlx.DB
}

func NewSchedules(db *sqlx.DB) Schedules {
	return &schedule{
		db: db,
	}
}

func (s *schedule) GetSchedules(pillDispenserSN string, contractID int) ([]ScheduleData, error) {
	var schedules []Schedule
	err := s.db.Select(&schedules, "SELECT * FROM schedule WHERE contract_id = $1 AND pill_dispenser_sn = $2 ORDER BY created_at DESC", contractID, pillDispenserSN)
	if err != nil {
		return nil, err
	}
	var schedulesData []ScheduleData
	for _, schedule := range schedules {
		var cells []ScheduleCell
		err := s.db.Select(&cells, "SELECT * FROM schedule_cell WHERE schedule_id = $1 ORDER BY idx", schedule.ID)
		if err != nil {
			return nil, err
		}
		for i := range cells {
			cells[i].StartTime.Time = cells[i].StartTime.Time.UTC()
			cells[i].EndTime.Time = cells[i].EndTime.Time.UTC()
		}
		schedulesData = append(schedulesData, ScheduleData{
			Schedule: schedule,
			Cells:    cells,
		})
	}
	return schedulesData, nil
}

func (s *schedule) GetScheduleForSN(serialNumber string) (*ScheduleData, error) {
	var schedule Schedule
	query := `
    SELECT s.* FROM schedule s
    JOIN pill_dispenser pd ON s.pill_dispenser_sn = pd.serial_number AND s.contract_id = pd.contract_id
    WHERE pd.serial_number = $1
    ORDER BY created_at DESC
    LIMIT 1
    `
	err := s.db.Get(&schedule, query, serialNumber)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoSchedule
	}
	if err != nil {
		return nil, err
	}
	var cells []ScheduleCell
	query = `
    SELECT * FROM schedule_cell WHERE schedule_id = $1 ORDER BY idx
    `
	err = s.db.Select(&cells, query, schedule.ID)
	if err != nil {
		return nil, err
	}
	for i := range cells {
		cells[i].StartTime.Time = cells[i].StartTime.Time.UTC()
		cells[i].EndTime.Time = cells[i].EndTime.Time.UTC()
	}
	return &ScheduleData{
		Schedule: schedule,
		Cells:    cells,
	}, nil
}

func (s *schedule) NewSchedule(schedule ScheduleData) (*ScheduleData, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return &schedule, err
	}

	query := `
    INSERT INTO schedule (is_offline_notifications_allowed, refresh_rate_interval, contract_id, pill_dispenser_sn) 
    VALUES (:is_offline_notifications_allowed, :refresh_rate_interval, :contract_id, :pill_dispenser_sn) RETURNING *
    `
	q, args, err := sqlx.Named(query, schedule.Schedule)
	if err != nil {
		return &schedule, err
	}
	q = tx.Rebind(q)
	err = tx.QueryRowx(q, args...).StructScan(&schedule.Schedule)
	if err != nil {
		_ = tx.Rollback()
		return &schedule, err
	}

	for i := range schedule.Cells {
		schedule.Cells[i].ScheduleID = schedule.Schedule.ID
	}
	query = `
    INSERT INTO schedule_cell (idx, schedule_id, start_time, end_time, contents_description) 
    VALUES (:idx, :schedule_id, :start_time, :end_time, :contents_description)
    `
	_, err = tx.NamedExec(query, schedule.Cells)
	if err != nil {
		_ = tx.Rollback()
		return &schedule, err
	}

	return &schedule, tx.Commit()
}

func (s *schedule) EditSchedule(schedule ScheduleData) (*ScheduleData, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return &schedule, err
	}

	query := `
    UPDATE schedule SET is_offline_notifications_allowed = :is_offline_notifications_allowed, 
                        refresh_rate_interval = :refresh_rate_interval, 
                        contract_id = :contract_id, 
                        pill_dispenser_sn = :pill_dispenser_sn
    WHERE id = :id
    RETURNING *
    `
	q, args, err := sqlx.Named(query, schedule.Schedule)
	if err != nil {
		return &schedule, err
	}
	q = tx.Rebind(q)
	err = tx.QueryRowx(q, args...).StructScan(&schedule.Schedule)
	if err != nil {
		_ = tx.Rollback()
		return &schedule, err
	}

	query = `
    UPDATE schedule_cell SET start_time = :start_time, end_time = :end_time, contents_description = :contents_description
    WHERE idx = :idx AND schedule_id = :schedule_id
    `
	for _, cell := range schedule.Cells {
		_, err := tx.NamedExec(query, cell)
		if err != nil {
			log.Printf("Error updating cell: %s\n", err.Error())
			_ = tx.Rollback()
			return &schedule, err
		}
	}

	return &schedule, tx.Commit()
}

func (s *schedule) GetPillNameAndContractID(serialNumber string, cellIndex int) (string, int, error) {
	query := `
    SELECT sc.contents_description, pd.contract_id FROM schedule_cell sc
    JOIN schedule s ON s.id = sc.schedule_id
    JOIN pill_dispenser pd ON s.pill_dispenser_sn = pd.serial_number AND s.contract_id = pd.contract_id
    WHERE pd.serial_number = $1 AND sc.idx = $2
    ORDER BY s.created_at DESC
    LIMIT 1
    `
	var pillName string
	var contractID int
	err := s.db.QueryRow(query, serialNumber, cellIndex).Scan(&pillName, &contractID)
	return pillName, contractID, err
}

func (s *schedule) GetLastScheduleForContractID(contractID int) (*ScheduleData, error) {
	var schedule Schedule
	// JOIN pill_dispenser pd ON s.pill_dispenser_sn = pd.serial_number AND s.contract_id = pd.contract_id
	query := `
    SELECT s.* FROM schedule s
    WHERE s.contract_id = $1
    ORDER BY created_at DESC
    LIMIT 1
    `
	err := s.db.Get(&schedule, query, contractID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNoSchedule
	}
	if err != nil {
		return nil, err
	}
	var cells []ScheduleCell
	query = `
    SELECT * FROM schedule_cell WHERE schedule_id = $1 ORDER BY idx
    `
	err = s.db.Select(&cells, query, schedule.ID)
	if err != nil {
		return nil, err
	}
	return &ScheduleData{
		Schedule: schedule,
		Cells:    cells,
	}, nil
}

type ScheduleData struct {
	Schedule Schedule
	Cells    []ScheduleCell
}

func NewSchedule(pillDispenser *PillDispenser) *ScheduleData {
	return &ScheduleData{
		Schedule: Schedule{
			IsOfflineNotificationsAllowed: true,
			RefreshRateInterval:           sql.NullInt64{Valid: true, Int64: DefaultRefreshRateInterval},
			ContractID:                    int(pillDispenser.ContractID.Int64),
			PillDispenserSN:               sql.NullString{Valid: true, String: pillDispenser.SerialNumber},
		},
		Cells: NewCellsSet(pillDispenser.HWType.GetCellsCount(), 0),
	}
}

func New4X4Schedule(contractID int) *ScheduleData {
	return &ScheduleData{
		Schedule: Schedule{
			IsOfflineNotificationsAllowed: true,
			RefreshRateInterval:           sql.NullInt64{Valid: true, Int64: DefaultRefreshRateInterval},
			ContractID:                    contractID,
			PillDispenserSN:               sql.NullString{Valid: false},
		},
		Cells: NewCellsSet(4, 0),
	}
}
