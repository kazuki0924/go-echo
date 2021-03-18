package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

func hello(c echo.Context) error {
	return c.String(http.StatusOK, "hello!")
}

type Book struct {
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
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

func createBook(c echo.Context) error {

	book := Book{}

	defer c.Request().Body.Close()

	body, err := ioutil.ReadAll(c.Request().Body)

	if err != nil {
		log.Printf("failed to read the request body: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	err = json.Unmarshal(body, &book)

	if err != nil {
		log.Printf("failed unmarshalling createBook: %s", err)
		return c.String(http.StatusInternalServerError, "")
	}

	log.Printf("this is your book: ")

	return c.String(http.StatusOK, "book created")
}

func main() {
	fmt.Println("Server started:")

	e := echo.New()

	e.GET("/", hello)
	e.GET("/books/:data", listBooks)

	e.POST("/books", createBook)

	e.Start(":8000")
}
