// Package handlers provides HTTP handlers for the pill dispenser service.
package handlers

import (
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	pilldispenserprotocol "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/pill_dispenser_protocol"
)

const ContentTypeOctetStream = "application/octet-stream"

func (pdh *PillDispenserHandler) GetSchedule(c echo.Context) error {
	serialNumber := c.QueryParam("serial_number")
	if serialNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provide serial number")
	}

	schedule, err := pdh.DB.Schedules().GetScheduleForSN(serialNumber)
	if errors.Is(err, models.ErrNoSchedule) {
		pillDispenser, err := pdh.DB.PillDispensers().Get(serialNumber)
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "pill dispenser not found: "+err.Error())
		}

		return c.Blob(http.StatusOK, ContentTypeOctetStream, pilldispenserprotocol.EmptySchedule(pillDispenser.HWType.GetCellsCount()))
	}
	if err != nil {
		return err
	}

	data := pilldispenserprotocol.ScheduleFromScheduleData(schedule)

	return c.Blob(http.StatusOK, ContentTypeOctetStream, data)
}
