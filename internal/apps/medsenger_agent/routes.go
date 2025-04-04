package medsenger_agent

import (
	"github.com/TikhonP/maigo"
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsenger_agent/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureMedsengerAgentGroup(g *echo.Group, cfg *config.Config, modelsFactory db.ModelsFactory, maigoClient *maigo.Client) {
	mah := handlers.MedsengerAgentHandler{
		MaigoClient: maigoClient,
	}

	g.GET("/", mah.MainPage)

	g.POST("/init", mah.Init, util.ApiKeyJSON(cfg.Server))
	g.POST("/status", mah.Status, util.ApiKeyJSON(cfg.Server))
	g.POST("/remove", mah.Remove, util.ApiKeyJSON(cfg.Server))

}
