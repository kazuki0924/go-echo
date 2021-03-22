package api

import (
	"github.com/kazuki0924/go-echo/src/api/controller"
	"github.com/labstack/echo"
)

func AdminGroup(g *echo.Group) {
	g.GET("/main", controller.MainAdmin)
}
