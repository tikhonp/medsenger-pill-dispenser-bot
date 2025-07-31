// Package providesnactionpage contains the routes for the Provide SN Action page.
package providesnactionpage

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/provide_sn_action_page/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureProvideSNActionGroup(g *echo.Group, deps util.Dependencies) {
	psnah := handlers.ProvideSNActionHandler(deps)

	g.Use(
		util.APIKeyGetParam(deps.Cfg.Server),
		util.ContractIDQueryParam(deps.DB),
	)

	g.GET("", psnah.Get)
	g.POST("", psnah.Post)
}
