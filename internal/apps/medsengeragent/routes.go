// Package medsengeragent provides the routes for the Medsenger Agent service.
package medsengeragent

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/medsengeragent/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureMedsengerAgentGroup(g *echo.Group, deps util.Dependencies) {
	mah := handlers.MedsengerAgentHandler(deps)

	g.Use(util.APIKeyJSON(deps.Cfg.Server))

	// g.Use(middleware.BodyDump(func(c echo.Context, reqBody, resBody []byte) {
	// 	fmt.Printf("Request Body: \n%s\n", reqBody)
	// }))

	g.POST("/init", mah.Init)
	g.POST("/status", mah.Status)
	g.POST("/remove", mah.Remove)
	g.POST("/order", mah.Order)
}
