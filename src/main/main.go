package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello!")
}

func listBooks(c echo.Context) error {
	bookTitle := c.QueryParam("title")
	bookPublishedAt := c.QueryParam("published-at")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf(
			"The book title is %s\nThe book is published at %s", bookTitle, bookPublishedAt,
		))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"title":       bookTitle,
			"publishedAt": bookPublishedAt,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to let us know if you want json or string data",
	})
}

func main() {
	fmt.Println("Server started:")

	e := echo.New()

	e.GET("/", hello)
	e.GET("/books/:data", listBooks)

	e.Start(":8000")
}
