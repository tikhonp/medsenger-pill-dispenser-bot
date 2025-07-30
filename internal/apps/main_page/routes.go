// Package mainpage provides the main page route for the application.
package mainpage

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tikhonp/medsenger-pill-dispenser-bot/internal/util"
)

func ConfigureMainPageGroup(g *echo.Group, deps util.Dependencies) {
	g.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Купил мужик шляпу, а она ему как раз!")
	})
}
