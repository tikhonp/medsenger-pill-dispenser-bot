package handlers

import (
	"io"
	"net/http"

	"github.com/labstack/echo/v4"
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

	pillName, contactId, err := pdh.Db.Schedules().GetPillNameAndContractID(serialNumber, cellIndex)
	if err != nil {
		return err
	}

	_, err = pdh.Maigo.AddRecord(contactId, "medicine", pillName, submitTime, nil)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
