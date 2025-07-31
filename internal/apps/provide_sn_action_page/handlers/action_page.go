// Package handlers provides HTTP handlers for the Medsenger Pill Dispenser Bot service.
package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/provide_sn_action_page/views"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func (psnah *ProvideSNActionHandler) Get(c echo.Context) error {
	return util.TemplRender(c, views.ActionPage(""))
}

func (psnah *ProvideSNActionHandler) Post(c echo.Context) error {
	contract, err := util.GetContract(c)
	if err != nil {
		return err
	}
	serialNumber := c.FormValue("serial-number")
	if serialNumber == "" {
		return util.TemplRender(c, views.ActionPage("Серийный номер не может быть пустым"))
	}
	regContractIDErr := psnah.DB.PillDispensers().RegisterContractID(serialNumber, contract.ID)
	if errors.Is(regContractIDErr, models.ErrPillDispenserNotExists) {
		return util.TemplRender(c, views.ActionPage("Устройство с таким серийным номером не найдено."))
	}
	if errors.Is(regContractIDErr, models.ErrContractIDAlreadySet) {
		return util.TemplRender(c, views.ActionPage("Это устройство уже привязано к другому контракту, сначала отвяжите."))
	}
	if regContractIDErr != nil {
		return err
	}

	schedule, err := psnah.DB.Schedules().GetLastScheduleForContractID(contract.ID)
	if err != nil {
		return err
	}
	schedule.Schedule.PillDispenserSN.String = serialNumber
	schedule.Schedule.PillDispenserSN.Valid = true
	_, err = psnah.DB.Schedules().EditSchedule(*schedule)
	if err != nil {
		return err
	}

	return util.TemplRender(c, views.ActionSuccess())
}
