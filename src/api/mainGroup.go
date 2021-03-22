package api

import (
	"github.com/kazuki0924/go-echo/src/api/controller"
	"github.com/labstack/echo"
)

func MainGroup(e *echo.Echo) {
	e.GET("/login", controller.Login)
	e.GET("/bookshelfBooks/:data", controller.GetBookshelfBook)

	e.POST("/bookshelfBook", controller.CreateBookshelfBook)
	e.POST("/bookshelfRepository", controller.CreateBookshelfRepository)
	e.POST("/bookshelfResearchPaper", controller.CreateBookshelfResearchPaper)
}
