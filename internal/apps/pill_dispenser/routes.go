package pilldispenser

import (
	"github.com/TikhonP/maigo"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/pill_dispenser/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
)

func ConfigurePillDispenserGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory, maigoClient *maigo.Client) {
	pdh := handlers.PillDispenserHandler{Db: modelsFactory}

	g.GET("/schedule", pdh.GetSchedule)
	g.POST("/submit", pdh.SubmitPills)
}
