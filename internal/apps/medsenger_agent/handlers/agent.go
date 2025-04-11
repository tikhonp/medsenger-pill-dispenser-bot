package handlers

import (
	"net/http"

	"github.com/TikhonP/maigo"
	"github.com/labstack/echo/v4"
)

type initModel struct {
	ContractId        int    `json:"contract_id" validate:"required"`
	ClinicId          int    `json:"clinic_id" validate:"required"`
	AgentToken        string `json:"agent_token" validate:"required"`
	PatientAgentToken string `json:"patient_agent_token" validate:"required"`
	DoctorAgentToken  string `json:"doctor_agent_token" validate:"required"`
	AgentId           int    `json:"agent_id" validate:"required"`
	AgentName         string `json:"agent_name" validate:"required"`
	Locale            string `json:"locale" validate:"required"`
}

func (mah *MedsengerAgentHandler) fetchContractDataOnInit(contractId int, ctx echo.Context) {
	ci, err := mah.Maigo.GetContractInfo(contractId)
	if err != nil {
		// sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}

	err = mah.Db.Contracts().UpdateContractWithPatientData(contractId, ci.PatientName, ci.PatientEmail)
	if err != nil {
		// sentry.CaptureException(err)
		ctx.Logger().Error(err)
		return
	}
	_, err = mah.Maigo.SendMessage(
		contractId,
		"Подключен агент для таблетницы",
		// maigo.WithAction("Настроить", "/setup", maigo.Action),
		maigo.OnlyDoctor(),
	)
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
	err := mah.Db.Contracts().NewContract(m.ContractId, m.ClinicId, m.AgentToken, m.PatientAgentToken, m.DoctorAgentToken, m.Locale)
	if err != nil {
		return err
	}
	go mah.fetchContractDataOnInit(m.ContractId, c)
	return c.String(http.StatusCreated, "ok")

}

type statusResponseModel struct {
	IsTrackingData     bool     `json:"is_tracking_data"`
	SupportedScenarios []string `json:"supported_scenarios"`
	TrackedContracts   []int    `json:"tracked_contracts"`
}

func (mah *MedsengerAgentHandler) Status(c echo.Context) error {
	trackedContracts, err := mah.Db.Contracts().GetActiveContractIds()
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

type contractIdModel struct {
	ContractId int `json:"contract_id" validate:"required"`
}

func (mah *MedsengerAgentHandler) Remove(c echo.Context) error {
	m := new(contractIdModel)
	if err := c.Bind(m); err != nil {
		return err
	}
	if err := c.Validate(m); err != nil {
		return err
	}
	if err := mah.Db.Contracts().MarkInactiveContractWithId(m.ContractId); err != nil {
		return err
	}
	return c.String(http.StatusCreated, "ok")

}
