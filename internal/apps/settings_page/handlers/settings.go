package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/settings_page/views"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func (mah *SettingsPageHandler) SettingsGet(c echo.Context) error {
	contractID, err := util.GetContractID(c)
	if err != nil {
		return err
	}
	pillDispensers, err := mah.DB.PillDispensers().GetAllByContractID(contractID)
	if err != nil {
		return err
	}
	contract, err := mah.DB.Contracts().Get(contractID)
	if err != nil {
		return err
	}
	agentToken := c.Get("agent_token").(string)
	return util.TemplRender(c, views.Settings(contract, pillDispensers, agentToken))
}

func (mah *SettingsPageHandler) AddContractPillDispenser(c echo.Context) error {
	contractID, err := util.GetContractID(c)
	if err != nil {
		return err
	}
	contract, err := mah.DB.Contracts().Get(contractID)
	if err != nil {
		return err
	}
	agentToken := c.Get("agent_token").(string)
	serialNumber := c.FormValue("serial-number")
	if serialNumber == "" {
		pillDispensers, err := mah.DB.PillDispensers().GetAllByContractID(contractID)
		if err != nil {
			return err
		}
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Введите серийный номер", agentToken))
	}
	regContractIDErr := mah.DB.PillDispensers().RegisterContractID(serialNumber, contract.ID)
	pillDispensers, err := mah.DB.PillDispensers().GetAllByContractID(contract.ID)
	if err != nil {
		return err
	}
	if errors.Is(regContractIDErr, models.ErrPillDispenserNotExists) {
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Устройство с таким серийным номером не найдено.", agentToken))
	}
	if errors.Is(regContractIDErr, models.ErrContractIDAlreadySet) {
		return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "Это устройство уже привязано к другому контракту, сначала отвяжите.", agentToken))
	}
	if regContractIDErr != nil {
		return err
	}
	return util.TemplRender(c, views.PillDispensersList(pillDispensers, contract, "", agentToken))
}

func (mah *SettingsPageHandler) RemoveContractPillDispenser(c echo.Context) error {
	serialNumber := c.FormValue("serial-number")
	if serialNumber == "" {
		return echo.NewHTTPError(http.StatusBadRequest, "provide serial number")
	}
	err := mah.DB.PillDispensers().UnregisterContractID(serialNumber)
	if err != nil {
		return err
	}
	return c.NoContent(http.StatusOK)
}

func (mah *SettingsPageHandler) pillDispenserPagesCommon(c echo.Context) (*models.Contract, *models.PillDispenser, error) {
	contractID, err := util.GetContractID(c)
	if err != nil {
		return nil, nil, err
	}
	serialNumber := c.Param("serial-number")
	pillDispenser, err := mah.DB.PillDispensers().Get(serialNumber)
	if err != nil {
		return nil, nil, echo.NewHTTPError(http.StatusNotFound)
	}
	if pillDispenser.ContractID.Int64 != int64(contractID) {
		return nil, nil, echo.NewHTTPError(http.StatusForbidden)
	}
	contract, err := mah.DB.Contracts().Get(contractID)
	if err != nil {
		return nil, nil, err
	}
	return contract, pillDispenser, nil
}

func (mah *SettingsPageHandler) SetScheduleGet(c echo.Context) error {
	contract, pillDispenser, err := mah.pillDispenserPagesCommon(c)
	if err != nil {
		return err
	}
	schedules, err := mah.DB.Schedules().GetSchedules(pillDispenser.SerialNumber, int(pillDispenser.ContractID.Int64))
	if err != nil {
		return err
	}
	loc, err := time.LoadLocation(contract.Timezone.String)
	if err != nil {
		return err
	}
	agentToken := c.Get("agent_token").(string)
	return util.TemplRender(c, views.ScheduleSettings(pillDispenser, schedules, contract, loc, agentToken))
}

func (mah *SettingsPageHandler) SetSchedulePost(c echo.Context) error {
	contract, pillDispenser, err := mah.pillDispenserPagesCommon(c)
	if err != nil {
		return err
	}

	locationStr := c.FormValue("timezone")
	loc, err := time.LoadLocation(locationStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid timezone: "+locationStr)
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
		cellStartTimeStr := c.FormValue("cell-start-time-" + strconv.Itoa(i))
		cellStartTime, err := time.ParseInLocation(util.HTMLInputTime, cellStartTimeStr, loc)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid cell time")
		}
		cellEndTimeStr := c.FormValue("cell-end-time-" + strconv.Itoa(i))
		cellEndTime, err := time.ParseInLocation(util.HTMLInputTime, cellEndTimeStr, loc)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid cell time")
		}
		contentsDescription := c.FormValue("cell-contents-description-" + strconv.Itoa(i))
		schedule.Cells[i] = models.ScheduleCell{
			Index:               i,
			StartTime:           sql.NullTime{Valid: true, Time: cellStartTime.UTC()},
			EndTime:             sql.NullTime{Valid: true, Time: cellEndTime.UTC()},
			ContentsDescription: sql.NullString{Valid: true, String: contentsDescription},
		}
	}

	newSchedule, err := mah.DB.Schedules().NewSchedule(*schedule)
	if err != nil {
		return err
	}

	agentToken := c.Get("agent_token").(string)
	return util.TemplRender(c, views.Schedule(newSchedule, pillDispenser, contract, false, loc, agentToken))
}

func (mah *SettingsPageHandler) EditSchedulePost(c echo.Context) error {
	contract, pillDispenser, err := mah.pillDispenserPagesCommon(c)
	if err != nil {
		return err
	}

	locationStr := c.FormValue("timezone")
	loc, err := time.LoadLocation(locationStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid timezone: "+locationStr)
	}

	schedule := models.NewSchedule(pillDispenser)

	scheduleIDStr := c.FormValue("schedule-id")
	scheduleID, err := strconv.Atoi(scheduleIDStr)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "invalid schedule id")
	}
	schedule.Schedule.ID = scheduleID

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
		cellStartTimeStr := c.FormValue("cell-start-time-" + strconv.Itoa(i))
		cellStartTime, err := time.ParseInLocation(util.HTMLInputTime, cellStartTimeStr, loc)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid cell time")
		}
		cellEndTimeStr := c.FormValue("cell-end-time-" + strconv.Itoa(i))
		cellEndTime, err := time.ParseInLocation(util.HTMLInputTime, cellEndTimeStr, loc)
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, "invalid cell time")
		}
		contentsDescription := c.FormValue("cell-contents-description-" + strconv.Itoa(i))
		schedule.Cells[i] = models.ScheduleCell{
			Index:               i,
			ScheduleID:          scheduleID,
			StartTime:           sql.NullTime{Valid: true, Time: cellStartTime.UTC()},
			EndTime:             sql.NullTime{Valid: true, Time: cellEndTime.UTC()},
			ContentsDescription: sql.NullString{Valid: true, String: contentsDescription},
		}
	}

	newSchedule, err := mah.DB.Schedules().EditSchedule(*schedule)
	if err != nil {
		return err
	}

	agentToken := c.Get("agent_token").(string)
	return util.TemplRender(c, views.Schedule(newSchedule, pillDispenser, contract, false, loc, agentToken))
}

func (mah *SettingsPageHandler) GetNewScheduleForm(c echo.Context) error {
	contract, pillDispenser, err := mah.pillDispenserPagesCommon(c)
	if err != nil {
		return err
	}
	schedule := models.NewSchedule(pillDispenser)
	loc, err := time.LoadLocation(contract.Timezone.String)
	if err != nil {
		return err
	}

	agentToken := c.Get("agent_token").(string)
	return util.TemplRender(c, views.Schedule(schedule, pillDispenser, contract, true, loc, agentToken))
}
