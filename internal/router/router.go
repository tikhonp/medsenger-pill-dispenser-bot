package router

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	mainpage "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/main_page"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsenger_agent"
	pilldispenser "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/pill_dispenser"
	providesnactionpage "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/provide_sn_action_page"
	settingspage "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/settings_page"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func New(cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Debug = cfg.Server.Debug

	e.Pre(middleware.RemoveTrailingSlash())

	if e.Debug {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: e.Logger.Output(),
		}))
	} else {
		e.Use(middleware.Logger())
	}

	e.Use(middleware.CORS())

	e.Validator = util.NewDefaultValidator()

	return e
}

func RegisterRoutes(e *echo.Echo, deps util.Dependencies) {
	mainpage.ConfigureMainPageGroup(e.Group(""), deps)
	pilldispenser.ConfigurePillDispenserGroup(e.Group("/pill-dispenser"), deps)
	medsenger_agent.ConfigureMedsengerAgentGroup(e.Group("/medsenger"), deps)
	settingspage.ConfigureSettingsPageGroup(e.Group("/medsenger/settings"), deps)
	providesnactionpage.ConfigureProvideSNActionGroup(e.Group("/medsenger/provide-sn"), deps)
}

func Start(e *echo.Echo, cfg *config.Config) error {
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	return e.Start(addr)
}
