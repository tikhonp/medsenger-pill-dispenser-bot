// Package confirmlayout provides the configuration for the layout of pill dispenser contents.
package confirmlayout

import (
	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/apps/confirm_layout/handlers"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureConfirmLayoutPageGroup(g *echo.Group, deps util.Dependencies) {
	clh := handlers.ConfirmLayoutHandler(deps)

	g.File("", "/public/frontend/index.html")
	g.File("/fill-2x2-dispenser", "/public/frontend/index.html")
	g.File("/fill-4x7-dispenser", "/public/frontend/index.html")

	g.GET("/large", clh.Large)
	g.GET("/small", clh.Small)
}
