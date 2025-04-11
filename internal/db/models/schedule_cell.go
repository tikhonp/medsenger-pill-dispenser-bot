package models

import "database/sql"

type ScheduleCell struct {
	ID         int          `db:"id"`
	Time       sql.NullTime `db:"time"`
	ScheduleID int          `db:"schedule_id"`
}
