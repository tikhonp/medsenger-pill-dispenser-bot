// Package handlers provides HTTP handlers for the Medsenger Agent service.
package handlers

import (
	"log"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/TikhonP/maigo"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db/models"
)

type initModel struct {
	ContractID        int    `json:"contract_id" validate:"required"`
	ClinicID          int    `json:"clinic_id" validate:"required"`
	AgentToken        string `json:"agent_token" validate:"required"`
	PatientAgentToken string `json:"patient_agent_token" validate:"required"`
	DoctorAgentToken  string `json:"doctor_agent_token" validate:"required"`
	AgentID           int    `json:"agent_id" validate:"required"`
	AgentName         string `json:"agent_name" validate:"required"`
	Locale            string `json:"locale" validate:"required"`
	Preset            string `json:"preset,omitempty"`
	Params            *struct {
		PillCells  string `json:"pill_cells"`
		Medicines  string `json:"medicines"`
		Algorithms string `json:"algorithms"`
		Forms      string `json:"forms"`
	} `json:"params,omitempty"`
}

func (mah *MedsengerAgentHandler) fetchContractDataOnInit(contractID int, ctx echo.Context) {
	ci, err := mah.Maigo.GetContractInfo(contractID)
	if err != nil {
		// sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}

	err = mah.DB.Contracts().UpdateContractWithPatientData(contractID, ci.PatientName, ci.PatientEmail, ci.ClinicTimezone)
	if err != nil {
		// sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	_, err = mah.Maigo.SendMessage(
		contractID,
		"Пожалуйста, введите серийный номер выданного Вам устройства.",
		maigo.WithAction("Ввести", "/provide-sn", maigo.Action),
		maigo.OnlyPatient(),
	)
	if err != nil {
		// sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
}

func (mah *MedsengerAgentHandler) SaveScheduleOnInit(m *initModel, ctx echo.Context) {
	pillsList := strings.Split(m.Params.PillCells, ",")
	var schedule *models.ScheduleData
	if len(pillsList) == 4 {
		schedule = models.New4X4Schedule(m.ContractID)
	} else if len(pillsList) == 28 {
		schedule = models.New4X7Schedule(m.ContractID)
	} else {
		ctx.Logger().Errorf("Cannot create schedule with %d pills", len(pillsList))
		return
	}
	for idx, pillID := range pillsList {
		schedule.Cells[idx].ContentsDescription.Valid = true
		schedule.Cells[idx].ContentsDescription.String = pillID
	}
	_, err := mah.DB.Schedules().NewSchedule(*schedule)
	if err != nil {
		// sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
}

func (mah *MedsengerAgentHandler) Init(c echo.Context) error {
	m := new(initModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	err := mah.DB.Contracts().NewContract(m.ContractID, m.ClinicID, m.AgentToken, m.PatientAgentToken, m.DoctorAgentToken, m.Locale)
	if err != nil {
		return err
	}
	go mah.fetchContractDataOnInit(m.ContractID, c)
	if m.Params != nil {
		go mah.SaveScheduleOnInit(m, c)
	}
	return c.String(http.StatusCreated, "ok")

}

type statusResponseModel struct {
	IsTrackingData     bool     `json:"is_tracking_data"`
	SupportedScenarios []string `json:"supported_scenarios"`
	TrackedContracts   []int    `json:"tracked_contracts"`
}

func (mah *MedsengerAgentHandler) Status(c echo.Context) error {
	trackedContracts, err := mah.DB.Contracts().GetActiveContractIds()
	if err != nil {
		return err
	}
	response := statusResponseModel{
		IsTrackingData:     true,
		SupportedScenarios: []string{},
		TrackedContracts:   trackedContracts,
	}
	return c.JSON(http.StatusOK, response)

}

type contractIDModel struct {
	ContractID int    `json:"contract_id" validate:"required"`
	AgentID    int    `json:"agent_id" validate:"required"`
	AgentName  string `json:"agent_name" validate:"required"`
}

func (mah *MedsengerAgentHandler) Remove(c echo.Context) error {
	m := new(contractIDModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	if err := mah.DB.Contracts().MarkInactiveContractWithID(m.ContractID); err != nil {
		return err
	}
	if err := mah.DB.PillDispensers().UnregisterByContractID(m.ContractID); err != nil {
		return err
	}
	return c.String(http.StatusCreated, "ok")

}

type orderModel struct {
	ContractID int    `json:"contract_id" validate:"required"`
	Order      string `json:"order"`
	Params     struct {
		Schedule []struct {
			Name      string `json:"name" validate:"required"`
			Timestamp int64  `json:"timestamp" validate:"required"`
		} `json:"schedule"`
	} `json:"params" validate:"required"`
	SenderID int `json:"sender_id"`
}

func (mah *MedsengerAgentHandler) Order(c echo.Context) error {
	m := new(orderModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	schedule, err := mah.DB.Schedules().GetLastScheduleForContractID(m.ContractID)
	if err != nil {
		return err
	}

	for _, schedulePoint := range m.Params.Schedule {
		indx := slices.IndexFunc(schedule.Cells, func(c models.ScheduleCell) bool {
			return c.ContentsDescription.Valid &&
				c.ContentsDescription.String == schedulePoint.Name &&
				!c.StartTime.Valid &&
				!c.EndTime.Valid
		})
		if indx != -1 {
			schedule.Cells[indx].StartTime.Time = time.Unix(schedulePoint.Timestamp, 0)
			schedule.Cells[indx].StartTime.Valid = true

			// + 10 minutes to allow some delay in taking medicine
			schedule.Cells[indx].EndTime.Time = time.Unix(schedulePoint.Timestamp+(60*10), 0)
			schedule.Cells[indx].EndTime.Valid = true
		} else {
			log.Printf("Cannot find pill with id %s in schedule for contract %d", schedulePoint.Name, m.ContractID)
		}
	}

	_, err = mah.DB.Schedules().EditSchedule(*schedule)
	if err != nil {
		log.Println("Error saving schedule:", err.Error())
		return err
	}

	return c.NoContent(http.StatusOK)
}
