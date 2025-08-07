package handlers

import (
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	pilldispenserprotocol "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/pill_dispenser_protocol"
)

func (pdh *PillDispenserHandler) SubmitPills(c echo.Context) error {
	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}

	submitTime, cellIndex, serialNumber, err := pilldispenserprotocol.PillDispenserSubmitBody(body).Decode()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	pillName, contactID, err := pdh.DB.Schedules().GetPillNameAndContractID(serialNumber, cellIndex)
	if err != nil {
		return err
	}

	_, err = pdh.Maigo.AddRecord(contactID, "medicine", pillName, submitTime, nil)
	if err != nil {
		return err
	}

	batteryVoltage := c.QueryParam("battery_voltage")
	if batteryVoltage != "" {
		batteryVoltageInt, err := strconv.Atoi(batteryVoltage)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid battery voltage: "+err.Error())
		}
		batteryStatus := models.BatteryStatus{
			SerialNumber: serialNumber,
			Voltage:      batteryVoltageInt,
			CreatedAt:    time.Now(),
		}
		_, err = pdh.DB.BatteryStatuses().InsertBatteryStatus(batteryStatus)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "failed to insert battery status: "+err.Error())
		}
	}

	return c.NoContent(http.StatusOK)
}
