package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func (mah *MedsengerAgentHandler) MainPage(c echo.Context) error {
	return c.String(http.StatusOK, "Купил мужик шляпу, а она ему как раз!")
}
