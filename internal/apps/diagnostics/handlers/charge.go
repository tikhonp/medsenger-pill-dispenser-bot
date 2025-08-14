// Package handlers provides HTTP handlers for the Medsenger Pill Dispenser Bot service.
package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/diagnostics/views"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func (psnah *ProvideDiagnosticsHandler) Get(c echo.Context) error {
	values := []float64{3.15, 3.28, 3.42, 3.55, 3.67, 3.79, 3.91, 4.05, 4.12, 3.98}
	timeMarkers := []string{
		"08:00", "09:00", "10:00", "11:00", "12:00",
		"13:00", "14:00", "15:00", "16:00", "17:00",
	}
	return util.TemplRender(c, views.ChargePage(values, timeMarkers, ""))
}