package providesnactionpage

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/provide_sn_action_page/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureProvideSNActionGroup(g *echo.Group, deps util.Dependencies) {
	psnah := handlers.ProvideSNActionHandler(deps)

	g.Use(
		util.ApiKeyGetParam(deps.Cfg.Server),
		util.ContractIdQueryParam(deps.Db),
	)

	g.GET("", psnah.Get)
	g.POST("", psnah.Post)
}
