package models

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
)

const DefaultRefreshRateInterval int64 = 60 * 60 // seconds

type Schedule struct {
	ID                            int           `db:"id"`
	IsOfflineNotificationsAllowed bool          `db:"is_offline_notifications_allowed"`
	RefreshRateInterval           sql.NullInt64 `db:"refresh_rate_interval"`
	ContractID                    int           `db:"contract_id"`
	PillDispenserSN               string        `db:"pill_dispenser_sn"`
	CreatedAt                     time.Time     `db:"created_at"`
}

type Schedules interface {
	GetSchedules(pillDispenserSN string, contractID int) ([]ScheduleData, error)
	NewSchedule(schedule ScheduleData) (*ScheduleData, error)
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
	for _, schd := range schedules {
		var cells []ScheduleCell
		err := s.db.Select(&cells, "SELECT * FROM schedule_cell WHERE schedule_id = $1 ORDER BY idx", schd.ID)
		if err != nil {
			return nil, err
		}
		schedulesData = append(schedulesData, ScheduleData{
			Schedule: schd,
			Cells:    cells,
		})
	}
	return schedulesData, nil
}

func (s *schedule) NewSchedule(schedule ScheduleData) (*ScheduleData, error) {
	query := `
    INSERT INTO schedule (is_offline_notifications_allowed, refresh_rate_interval, contract_id, pill_dispenser_sn) 
    VALUES (:is_offline_notifications_allowed, :refresh_rate_interval, :contract_id, :pill_dispenser_sn) RETURNING id, created_at
    `
	rows, err := s.db.NamedQuery(query, schedule.Schedule)
	if err != nil {
		return &schedule, err
	}
	if rows.Next() {
		err = rows.Scan(&schedule.Schedule.ID, &schedule.Schedule.CreatedAt)
		if err != nil {
			return &schedule, err
		}
	}
	err = rows.Close()
	if err != nil {
		return &schedule, err
	}
	for i := range schedule.Cells {
		schedule.Cells[i].ScheduleID = schedule.Schedule.ID
	}
	query = `
    INSERT INTO schedule_cell (idx, schedule_id, time) VALUES (:idx, :schedule_id, :time)
    `
	_, err = s.db.NamedExec(query, schedule.Cells)
	return &schedule, err
}

type ScheduleData struct {
	Schedule Schedule
	Cells    []ScheduleCell
}

func NewSchedule(pillDispenser *PillDispenser) ScheduleData {
	return ScheduleData{
		Schedule: Schedule{
			IsOfflineNotificationsAllowed: true,
			RefreshRateInterval:           sql.NullInt64{Valid: true, Int64: DefaultRefreshRateInterval},
			ContractID:                    int(pillDispenser.ContractID.Int64),
			PillDispenserSN:               pillDispenser.SerialNumber,
		},
		Cells: NewCellsSet(pillDispenser.HWType.GetCellsCount(), 0),
	}
}
