package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsenger_agent/views"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func (mah *MedsengerAgentHandler) SettingsGet(c echo.Context) error {
	return util.TemplRender(c, views.Settings())
}
