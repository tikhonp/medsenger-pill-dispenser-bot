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
	regContractIdErr := psnah.Db.PillDispensers().RegisterContractID(serialNumber, contract.ID)
	if errors.Is(regContractIdErr, models.ErrPillDispenserNotExists) {
		return util.TemplRender(c, views.ActionPage("Устройство с таким серийным номером не найдено."))
	}
	if errors.Is(regContractIdErr, models.ErrContractIdAlreadySet) {
		return util.TemplRender(c, views.ActionPage("Это устройство уже привязано к другому контракту, сначала отвяжите."))
	}
	if regContractIdErr != nil {
		return err
	}
	return util.TemplRender(c, views.ActionSuccess())
}
