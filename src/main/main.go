package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type BookshelfBook struct {
	Title       string `json:"title"`
	PublishedAt string `json:"publishedAt"`
}

type BookshelfRepository struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

type BookshelfResearchPaper struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func listBookshelfBooks(c echo.Context) error {
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

func createBookshelfBook(c echo.Context) error {
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

func createBookshelfRepository(c echo.Context) error {
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

func createBookshelfResearchPaper(c echo.Context) error {
	researchPaper := BookshelfResearchPaper{}

	err := c.Bind(&researchPaper)

	if err != nil {
		log.Printf("failed to read the request body: %s", err)
		return echo.NewHTTPError(http.StatusInternalServerError, "")
	}

	log.Printf("this is your research paper: %#v", researchPaper)

	return c.String(http.StatusOK, "research paper created")

}

func mainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "admin main page")
}

func main() {
	fmt.Println("Server started:")

	e := echo.New()

	g := e.Group("/admin")

	g.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	g.GET("/main", mainAdmin)

	e.GET("/bookshelfBooks/:data", listBookshelfBooks)

	e.POST("/bookshelfBook", createBookshelfBook)
	e.POST("/bookshelfRepository", createBookshelfRepository)
	e.POST("/bookshelfResearchPaper", createBookshelfResearchPaper)

	e.Start(":8000")
}
