package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsenger_agent/views"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func (mah *MedsengerAgentHandler) SettingsGet(c echo.Context) error {
	contractIdStr := c.QueryParam("contract_id")
	if contractIdStr == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provide contract id")
	}
	contractID, err := strconv.Atoi(contractIdStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contract id", err)
	}
	contract, err := mah.Db.Contracts().Get(contractID)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "contract not found", err)
	}
	pillDispensers, err := mah.Db.PillDispensers().GetAllByContractID(contractID)
	if err != nil {
		return err
	}

	return util.TemplRender(c, views.Settings(contract, pillDispensers))
}

func (mah *MedsengerAgentHandler) AddContractPillDispenser(c echo.Context) error {
	contract, err := util.GetContract(c)
	if err != nil {
		return err
	}
	serialNumber := c.FormValue("serial-number")
	if serialNumber == "" {
		pillDispensers, err := mah.Db.PillDispensers().GetAllByContractID(contract.ID)
		if err != nil {
			return err
		}
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Введите серийный номер"))
	}
	regContrIDerr := mah.Db.PillDispensers().RegisterContractID(serialNumber, contract.ID)
	pillDispensers, err := mah.Db.PillDispensers().GetAllByContractID(contract.ID)
	if err != nil {
		return err
	}
	if regContrIDerr == models.ErrPillDispenserNotExists {
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Устройство с таким серийным номером не найдено."))
	}
	if regContrIDerr == models.ErrContractIdAlreadySet {
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Это устройство уже привязано к другому контракту, сначала отвяжите."))
	}
	if regContrIDerr != nil {
		return err
	}
	return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, ""))
}

func (mah *MedsengerAgentHandler) RemoveContractPillDispenser(c echo.Context) error {
	serialNumber := c.FormValue("serial-number")
	if serialNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provide serial number")
	}
	err := mah.Db.PillDispensers().UnregisterContractID(serialNumber)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (mah *MedsengerAgentHandler) SetScheduleGet(c echo.Context) error {
	serialNumber := c.Param("serial-number")

	pillDispenser, err := mah.Db.PillDispensers().Get(serialNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	return util.TemplRender(c, views.ScheduleSettings(pillDispenser))
}

func (mah *MedsengerAgentHandler) SetSchedulePost(c echo.Context) error {
	contract, err := util.GetContract(c)
	if err != nil {
		return err
	}

	serialNumber := c.Param("serial-number")

	pillDispenser, err := mah.Db.PillDispensers().Get(serialNumber)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound)
	}

	if pillDispenser.ContractID.Int64 != int64(contract.ID) {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	offlineNotify := c.FormValue("offline-notify")
	fmt.Println("offline-notify", offlineNotify)

    for i := range pillDispenser.HWType.GetCellsCount() {
        cellTime := c.FormValue("cell-time-"+strconv.Itoa(i))
        fmt.Println(cellTime, i)
    }

	return c.NoContent(http.StatusOK)
}
