package diagnosticspage

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/diagnostics/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureDiagnosticsGroup(g *echo.Group, deps util.Dependencies) {
	cd := handlers.ProvideDiagnosticsHandler(deps)

	g.GET("", cd.Get)
}
