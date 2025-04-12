package medsenger_agent

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsenger_agent/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureMedsengerAgentGroup(g *echo.Group, deps util.Dependencies) {
	mah := handlers.MedsengerAgentHandler(deps)

	g.GET("/", mah.MainPage)

	medsengerAgentGroup := g.Group("", util.ApiKeyJSON(deps.Cfg.Server))
	medsengerAgentGroup.POST("/init", mah.Init)
	medsengerAgentGroup.POST("/status", mah.Status)
	medsengerAgentGroup.POST("/remove", mah.Remove)

	g.GET("/settings", mah.SettingsGet, util.ApiKeyGetParam(deps.Cfg.Server))

	agentApiGetGroup := g.Group("/settings/set-schedule", util.AgentTokenGetParam(mah.Db))
	agentApiGetGroup.GET("/:serial-number", mah.SetScheduleGet)
	agentApiGetGroup.POST("/:serial-number", mah.SetSchedulePost)
	agentApiGetGroup.GET("/:serial-number/new-schedule-form", mah.GetNewScheduleForm)

	agentApiFormGroup := g.Group("/agent", util.AgentTokenForm(mah.Db))
	agentApiFormGroup.POST("/contract-pill-dispenser", mah.AddContractPillDispenser)
	agentApiFormGroup.DELETE("/contract-pill-dispenser", mah.RemoveContractPillDispenser)
}
