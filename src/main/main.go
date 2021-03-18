package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello!")
}

func main() {
	fmt.Println("Server started:")

	e := echo.New()

	e.GET("/", hello)

	e.Start(":8000")
}
