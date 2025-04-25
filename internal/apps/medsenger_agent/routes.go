package medsenger_agent

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsenger_agent/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureMedsengerAgentGroup(g *echo.Group, deps util.Dependencies) {
	mah := handlers.MedsengerAgentHandler(deps)

	g.Use(util.ApiKeyJSON(deps.Cfg.Server))

	g.POST("/init", mah.Init)
	g.POST("/status", mah.Status)
	g.POST("/remove", mah.Remove)
	g.POST("/order", mah.Order)
}
