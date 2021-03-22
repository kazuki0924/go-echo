package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type BookshelfBook struct {
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
}

func GetBookshelfBook(c echo.Context) error {
	title := c.QueryParam("title")
	publishedAt := c.QueryParam("published-at")

	dataType := c.Param("data")

	if dataType == "string" {
		return c.String(http.StatusOK, fmt.Sprintf(
			"The book title is %s\nThe book is published at %s", title, publishedAt,
		))
	}

	if dataType == "json" {
		return c.JSON(http.StatusOK, map[string]string{
			"title":       title,
			"publishedAt": publishedAt,
		})
	}

	return c.JSON(http.StatusBadRequest, map[string]string{
		"error": "you need to let us know if you want json or string data",
	})
}

func CreateBookshelfBook(c echo.Context) error {
	book := BookshelfBook{}

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

	log.Printf("this is your book: %#v", book)

	return c.String(http.StatusOK, "book created")
}
