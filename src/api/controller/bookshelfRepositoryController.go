package controller

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type BookshelfRepository struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func CreateBookshelfRepository(c echo.Context) error {
	repository := BookshelfRepository{}

	defer c.Request().Body.Close()

	err := json.NewDecoder(c.Request().Body).Decode(&repository)

	if err != nil {
		log.Printf("failed to read the request body: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	log.Printf("this is your repository: %#v", repository)

	return c.String(http.StatusOK, "repository created")
}
