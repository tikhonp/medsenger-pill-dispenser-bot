package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

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
		return echo.NewHTTPError(http.StatusBadRequest, "invalid contract id")
	}
	contract, err := mah.Db.Contracts().Get(contractID)
	if errors.Is(err, sql.ErrNoRows) {
		return echo.NewHTTPError(http.StatusNotFound, "contract not found")
	} else if err != nil {
		return err
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
	regContractIdErr := mah.Db.PillDispensers().RegisterContractID(serialNumber, contract.ID)
	pillDispensers, err := mah.Db.PillDispensers().GetAllByContractID(contract.ID)
	if err != nil {
		return err
	}
	if errors.Is(regContractIdErr, models.ErrPillDispenserNotExists) {
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Устройство с таким серийным номером не найдено."))
	}
	if errors.Is(regContractIdErr, models.ErrContractIdAlreadySet) {
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Это устройство уже привязано к другому контракту, сначала отвяжите."))
	}
	if regContractIdErr != nil {
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

func (mah *MedsengerAgentHandler) pillDispenserPagesCommon(c echo.Context) (*models.Contract, *models.PillDispenser, error) {
	contract, err := util.GetContract(c)
	if err != nil {
		return nil, nil, err
	}
	serialNumber := c.Param("serial-number")
	pillDispenser, err := mah.Db.PillDispensers().Get(serialNumber)
	if err != nil {
		return nil, nil, echo.NewHTTPError(http.StatusNotFound)
	}
	if pillDispenser.ContractID.Int64 != int64(contract.ID) {
		return nil, nil, echo.NewHTTPError(http.StatusForbidden)
	}
	return contract, pillDispenser, nil
}

func (mah *MedsengerAgentHandler) SetScheduleGet(c echo.Context) error {
	contract, pillDispenser, err := mah.pillDispenserPagesCommon(c)
	if err != nil {
		return err
	}
	schedules, err := mah.Db.Schedules().GetSchedules(pillDispenser.SerialNumber, int(pillDispenser.ContractID.Int64))
	if err != nil {
		return err
	}
	return util.TemplRender(c, views.ScheduleSettings(pillDispenser, schedules, contract))
}

func (mah *MedsengerAgentHandler) SetSchedulePost(c echo.Context) error {
	contract, pillDispenser, err := mah.pillDispenserPagesCommon(c)
	if err != nil {
		return err
	}

	schedule := models.NewSchedule(pillDispenser)

	offlineNotify := c.FormValue("offline-notify") == "on"
	schedule.Schedule.IsOfflineNotificationsAllowed = offlineNotify

	refreshRateIntervalStr := c.FormValue("refresh-rate")
	refreshRateInterval, err := strconv.Atoi(refreshRateIntervalStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid refresh-rate")
	}
	schedule.Schedule.RefreshRateInterval = sql.NullInt64{Valid: true, Int64: int64(refreshRateInterval)}

	schedule.Cells = make([]models.ScheduleCell, pillDispenser.HWType.GetCellsCount())
	for i := range pillDispenser.HWType.GetCellsCount() {
		cellTimeStr := c.FormValue("cell-time-" + strconv.Itoa(i))
		cellTime, err := time.Parse("2006-01-02T15:04", cellTimeStr)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid cell time")
		}
		schedule.Cells[i] = models.ScheduleCell{
			Index: i,
			Time:  sql.NullTime{Valid: true, Time: cellTime},
		}
	}

	newSchedule, err := mah.Db.Schedules().NewSchedule(schedule)
	if err != nil {
		return err
	}

	return util.TemplRender(c, views.Schedule(*newSchedule, pillDispenser, contract, false))
}

func (mah *MedsengerAgentHandler) GetNewScheduleForm(c echo.Context) error {
	contract, pillDispenser, err := mah.pillDispenserPagesCommon(c)
	if err != nil {
		return err
	}
	schedule := models.NewSchedule(pillDispenser)
	return util.TemplRender(c, views.Schedule(schedule, pillDispenser, contract, true))
}
