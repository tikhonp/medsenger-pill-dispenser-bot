package handlers

import (
	"encoding/binary"
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
)

const ContentTypeOctetStream = "application/octet-stream"

func encodeSchedule(s *models.ScheduleData) []byte {
	cellsCount := len(s.Cells)

	// (uint32, uint32, uint8) * cell-count  + uint32
	out := make([]byte, 0, cellsCount*(4+4+1)+4)
	for _, cell := range s.Cells {
		// make array of timestamps from now and every minute after
		timestampStart := uint32(cell.Time.Time.Unix())
		timestampEnd := uint32(cell.Time.Time.Add(time.Minute).Add(time.Second * 30).Unix())
		out = binary.LittleEndian.AppendUint32(out, timestampStart)
		out = binary.LittleEndian.AppendUint32(out, timestampEnd)
		if s.Schedule.IsOfflineNotificationsAllowed {
			out = append(out, 1)
		} else {
			out = append(out, 0)
		}
	}

	var refreshRateSec = uint32(s.Schedule.RefreshRateInterval.Int64)
	out = binary.LittleEndian.AppendUint32(out, refreshRateSec)

	return out
}

func (pdh *PillDispenserHandler) GetSchedule(c echo.Context) error {
	serialNumber := c.QueryParam("serial_number")
	if serialNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provide serial number")
	}

	schedule, err := pdh.Db.Schedules().GetScheduleForSN(serialNumber)
	if errors.Is(err, models.ErrNoSchedule) {
		// TODO: process empty schedule
		return err
	}
	if err != nil {
		return err
	}

	data := encodeSchedule(schedule)

	return c.Blob(http.StatusOK, ContentTypeOctetStream, data)
}
