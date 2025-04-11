package models

import "database/sql"

type ScheduleCell struct {
	Index      int          `db:"idx"`
	ScheduleID int          `db:"schedule_id"`
	Time       sql.NullTime `db:"time"`
}

func NewCellsSet(n, scheduleID int) []ScheduleCell {
    var cells = make([]ScheduleCell, 0, n)
    for i := range n {
       cells = append(cells, ScheduleCell{Index: i, ScheduleID: scheduleID})
    }
    return cells
}
