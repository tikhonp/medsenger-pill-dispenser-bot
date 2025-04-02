package router

import (
	"fmt"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	mainpage "github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/main_page"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/config"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/db"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func New(cfg *config.Config) *echo.Echo {
	e := echo.New()

	e.HideBanner = true
	e.Debug = cfg.Server.Debug

	e.Logger.SetLevel(log.DEBUG)

	e.Pre(middleware.RemoveTrailingSlash())

	if e.Debug {
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: "[${time_rfc3339}] ${status} ${method} ${path} (${remote_ip}) ${latency_human}\n",
			Output: e.Logger.Output(),
		}))
	} else {
		e.Use(middleware.Logger())
	}
	// e.Use(middleware.Recover())

	e.Use(middleware.CORSWithConfig(
		middleware.CORSConfig{
			// TODO: Set proper CORS configuration
			AllowOrigins: []string{"*"},
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
			AllowMethods: []string{echo.GET},
		},
	))

	e.Validator = util.NewDefaultValidator()

	return e
}

func RegisterRoutes(e *echo.Echo, cfg *config.Config, modelsFactory db.ModelsFactory) {
	mainpage.ConfigureMainPageGroup(e.Group(""), cfg, modelsFactory)
}

func Start(e *echo.Echo, cfg *config.Config) error {
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	return e.Start(addr)
}
