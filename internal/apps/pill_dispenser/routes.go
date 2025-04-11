package pilldispenser

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/pill_dispenser/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigurePillDispenserGroup(g *echo.Group, deps util.Dependencies) {
	pdh := handlers.PillDispenserHandler(deps)

	g.GET("/schedule", pdh.GetSchedule)
	g.POST("/submit", pdh.SubmitPills)
}
