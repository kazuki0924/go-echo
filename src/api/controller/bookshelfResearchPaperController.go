package controller

import (
	"log"
	"net/http"

	"github.com/labstack/echo"
)

type BookshelfResearchPaper struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func CreateBookshelfResearchPaper(c echo.Context) error {
	researchPaper := BookshelfResearchPaper{}

	err := c.Bind(&researchPaper)

	if err != nil {
		log.Printf("failed to read the request body: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	log.Printf("this is your research paper: %#v", researchPaper)

	return c.String(http.StatusOK, "research paper created")

}
