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

	return c.String(http.StatusOK, fmt.Sprintf(
		"The book title is %s\nThe book is published at %s", bookTitle, bookPublishedAt,
	))

	// http://localhost:8000/books?title=title&published-at=published-at

	// The book title is title
	// The book is published at published-at
}

func main() {
	fmt.Println("Server started:")

	e := echo.New()

	e.GET("/", hello)
	e.GET("/books", listBooks)

	e.Start(":8000")
}
