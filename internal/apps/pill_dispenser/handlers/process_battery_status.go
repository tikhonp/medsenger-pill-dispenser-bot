package handlers

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/tikhonp/maigo"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	pilldispenserprotocol "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/pill_dispenser_protocol"
)

func (pdh *PillDispenserHandler) ProcessBatteryStatus(serialNumber string, batteryVoltage string) error {
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
		return err
	}

	if pilldispenserprotocol.IsBatteryLow(batteryVoltageInt) {
		contractID, err := pdh.DB.PillDispensers().GetContractID(serialNumber)
		if err == nil {
			_, err = pdh.Maigo.SendMessage(contractID, "⚠️ Низкий заряд батареи в вашей таблетнице. Пожалуйста, замените батарейку как можно скорее, чтобы избежать сбоев в работе устройства.", maigo.OnlyPatient())
			if err != nil {
				log.Printf("failed to send low battery message to contract %d: %v", contractID, err)
			}
		}
	}

	return nil
}
