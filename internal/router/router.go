// Package router provides the HTTP server and routes for the application.
package router

import (
	"fmt"
	"time"

	sentryecho "github.com/getsentry/sentry-go/echo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	confirmlayout "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/confirm_layout"
	mainpage "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/main_page"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsengeragent"
	pilldispenser "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/pill_dispenser"
	providesnactionpage "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/provide_sn_action_page"
	settingspage "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/settings_page"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util/config"
)

func New(cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Debug = cfg.Server.Debug
	e.Validator = util.NewDefaultValidator()

	e.Pre(middleware.RemoveTrailingSlash())

	if e.Debug {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: e.Logger.Output(),
		}))
	} else {
		e.Use(sentryecho.New(sentryecho.Options{
			Repanic:         true,
			WaitForDelivery: false,
			Timeout:         5 * time.Second,
		}))
		e.Use(middleware.Logger())
	}

	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	return e
}

func RegisterRoutes(e *echo.Echo, deps util.Dependencies) {
	mainpage.ConfigureMainPageGroup(e.Group(""), deps)
	pilldispenser.ConfigurePillDispenserGroup(e.Group("/pill-dispenser"), deps)
	medsengeragent.ConfigureMedsengerAgentGroup(e.Group("/medsenger"), deps)
	settingspage.ConfigureSettingsPageGroup(e.Group("/medsenger/settings"), deps)
	providesnactionpage.ConfigureProvideSNActionGroup(e.Group("/provide-sn"), deps)
	confirmlayout.ConfigureConfirmLayoutPageGroup(e.Group("/confirm-layout"), deps)
}

func Start(e *echo.Echo, cfg *config.Config) error {
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	return e.Start(addr)
}
