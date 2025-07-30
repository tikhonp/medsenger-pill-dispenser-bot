package settingspage

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/settings_page/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureSettingsPageGroup(g *echo.Group, deps util.Dependencies) {
	sph := handlers.SettingsPageHandler(deps)

	g.GET("", sph.SettingsGet, util.APIKeyGetParam(deps.Cfg.Server), util.ContractIDQueryParam(deps.Db))

	agentApiGetGroup := g.Group("/set-schedule", util.AgentTokenGetParam(deps.Db))
	agentApiGetGroup.GET("/:serial-number", sph.SetScheduleGet)
	agentApiGetGroup.POST("/:serial-number", sph.SetSchedulePost)
	agentApiGetGroup.GET("/:serial-number/new-schedule-form", sph.GetNewScheduleForm)

	g.POST("/edit-schedule/:serial-number", sph.EditSchedulePost, util.AgentTokenGetParam(deps.Db))

	agentApiFormGroup := g.Group("/pill-dispenser", util.AgentTokenForm(sph.Db))
	agentApiFormGroup.POST("", sph.AddContractPillDispenser)
	agentApiFormGroup.DELETE("", sph.RemoveContractPillDispenser)
}
