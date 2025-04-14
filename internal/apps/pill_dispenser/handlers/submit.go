package handlers

import (
	"encoding/binary"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (pdh *PillDispenserHandler) SubmitPills(c echo.Context) error {

	body, err := io.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	if len(body) < 5 {
		fmt.Printf("SUBMIT: recieved body less than 5 bytes (%d)", len(body))
		return echo.NewHTTPError(http.StatusBadRequest, "data length must be at least 5 bytes")
	}

	var timestamp uint32
	_, err = binary.Decode(body[:4], binary.BigEndian, &timestamp)
	if err != nil {
		fmt.Printf("SUBMIT: err: %s", err.Error())
		return err
	}
	submitTime := time.Unix(int64(timestamp), 0)
	cellIndex := int(body[4])
	serialNumber := string(body[5:])

	pillName, conractId, err := pdh.Db.Schedules().GetPillNameAndContractID(serialNumber, cellIndex)
	if err != nil {
		return err
	}

	_, err = pdh.Maigo.AddRecord(conractId, "medicine", pillName, submitTime, nil)
	if err != nil {
		return err
	}

	return c.NoContent(http.StatusOK)
}
