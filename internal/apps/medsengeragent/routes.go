// Package medsengeragent provides the routes for the Medsenger Agent service.
package medsengeragent

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsengeragent/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureMedsengerAgentGroup(g *echo.Group, deps util.Dependencies) {
	mah := handlers.MedsengerAgentHandler(deps)

	g.Use(util.APIKeyJSON(deps.Cfg.Server))

	g.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
		c.Logger().Infof("Request: %s %s, Body: %s", c.Request().Method, c.Request().URL, string(reqBody))
		c.Logger().Infof("Response: %s %s, Body: %s", c.Request().Method, c.Request().URL, string(resBody))
	}))

	g.POST("/init", mah.Init)
	g.POST("/status", mah.Status)
	g.POST("/remove", mah.Remove)
	g.POST("/order", mah.Order)
}
