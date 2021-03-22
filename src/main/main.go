package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
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

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
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

func mainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "cookie")
}

func login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "username" && password == "password" {
		cookie := &http.Cookie{} // new(http.)

		cookie.Name = "sessionID"
		cookie.Value = "secret"
		cookie.Expires = time.Now().Add(24 * 2 * time.Hour)

		c.SetCookie(cookie)

		// return c.String(http.StatusOK, "you are logged in")
		// create jwt token
		token, err := createJwtToken()
		if err != nil {
			log.Println("Error Creating JWT token", err)
			return c.String(http.StatusInternalServerError, "something went wrong")
		}

		jwtCookie := &http.Cookie{}

		jwtCookie.Name = "JWTCookie"
		jwtCookie.Value = token
		jwtCookie.Expires = time.Now().Add(24 * 2 * time.Hour)

		c.SetCookie(jwtCookie)

		return c.JSON(http.StatusOK, map[string]string{
			"message": "You were logged in!",
			"token":   token,
		})
	}

	return c.String(http.StatusUnauthorized, "your username or password is wrong")
}

func createJwtToken() (string, error) {
	claims := JwtClaims{
		"user",
		jwt.StandardClaims{
			Id:        "main_user_id",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
		},
	}

	rawToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	token, err := rawToken.SignedString([]byte("mySecret"))
	if err != nil {
		return "", err
	}

	return token, nil
}

func mainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	log.Println("User Name: ", claims["name"], "User ID: ", claims["jti"])

	return c.String(http.StatusOK, "jwt")
}

// middlewares
func ServerHeader(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set(echo.HeaderServer, "TestServer/1.0")

		return next(c)
	}
}

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")

		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.String(http.StatusUnauthorized, "you don't have any cookie")
			}
			return err
		}

		if cookie.Value == "secret" {
			return next(c)
		}

		return c.String(http.StatusUnauthorized, "you don't have the right cookie")
	}
}

func main() {
	fmt.Println("Server started:")

	e := echo.New()

	e.Use(ServerHeader)

	cookieGroup := e.Group("/cookie")
	adminGroup := e.Group("/admin")
	jwtGroup := e.Group("/jwt")

	e.Use(middleware.StaticWithConfig(middleware.StaticConfig{
		Root: "./static",
	}))

	adminGroup.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: `[${time_rfc3339}] ${status} ${method} ${host}${path} ${latency_human}` + "\n",
	}))

	adminGroup.Use(middleware.BasicAuth(func(username, password string, c echo.Context) (bool, error) {
		if username == "username" && password == "password" {
			return true, nil
		}

		return false, nil
	}))

	cookieGroup.Use(checkCookie)

	jwtGroup.Use(middleware.JWTWithConfig(middleware.JWTConfig{
		SigningMethod: "HS512",
		SigningKey:    []byte("mySecret"),
	}))

	cookieGroup.GET("/main", mainCookie)
	adminGroup.GET("/main", mainAdmin)
	jwtGroup.GET("/main", mainJwt)

	e.GET("/login", login)
	e.GET("/bookshelfBooks/:data", listBookshelfBooks)

	e.POST("/bookshelfBook", createBookshelfBook)
	e.POST("/bookshelfRepository", createBookshelfRepository)
	e.POST("/bookshelfResearchPaper", createBookshelfResearchPaper)

	e.Start(":8000")
}
