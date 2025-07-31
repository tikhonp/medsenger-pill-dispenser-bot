// Package settingspage provides the routes for the settings page of the application.
package settingspage

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/settings_page/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureSettingsPageGroup(g *echo.Group, deps util.Dependencies) {
	sph := handlers.SettingsPageHandler(deps)

	g.GET("", sph.SettingsGet, util.APIKeyGetParam(deps.Cfg.Server), util.ContractIDQueryParam(deps.DB))

	agentAPIGetGroup := g.Group("/set-schedule", util.AgentTokenGetParam(deps.DB))
	agentAPIGetGroup.GET("/:serial-number", sph.SetScheduleGet)
	agentAPIGetGroup.POST("/:serial-number", sph.SetSchedulePost)
	agentAPIGetGroup.GET("/:serial-number/new-schedule-form", sph.GetNewScheduleForm)

	g.POST("/edit-schedule/:serial-number", sph.EditSchedulePost, util.AgentTokenGetParam(deps.DB))

	agentAPIFormGroup := g.Group("/pill-dispenser", util.AgentTokenForm(sph.DB))
	agentAPIFormGroup.POST("", sph.AddContractPillDispenser)
	agentAPIFormGroup.DELETE("", sph.RemoveContractPillDispenser)
}
