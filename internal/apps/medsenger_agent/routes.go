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
		Db:          modelsFactory,
		MaigoClient: maigoClient,
	}

	g.GET("/", mah.MainPage)

	g.POST("/init", mah.Init, util.ApiKeyJSON(cfg.Server))
	g.POST("/status", mah.Status, util.ApiKeyJSON(cfg.Server))
	g.POST("/remove", mah.Remove, util.ApiKeyJSON(cfg.Server))

	g.GET("/settings", mah.SettingsGet, util.ApiKeyGetParam(cfg.Server))

	agentGroup := g.Group("/agent", util.AgentTokenForm(mah.Db))
	agentGroup.POST("/contract-pill-dispenser", mah.AddContractPillDispenser)
	agentGroup.DELETE("/contract-pill-dispenser", mah.RemoveContractPillDispenser)
}
