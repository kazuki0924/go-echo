package controller

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type JwtClaims struct {
	Name string `json:"name"`
	jwt.StandardClaims
}

func Login(c echo.Context) error {
	username := c.QueryParam("username")
	password := c.QueryParam("password")

	if username == "username" && password == "password" {
		cookie := &http.Cookie{} // new(http.)

		cookie.Name = "sessionID"
		cookie.Value = "secret"
		cookie.Expires = time.Now().Add(24 * 2 * time.Hour)

		c.SetCookie(cookie)

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

func MainAdmin(c echo.Context) error {
	return c.String(http.StatusOK, "admin main page")
}

func MainCookie(c echo.Context) error {
	return c.String(http.StatusOK, "cookie")
}

func MainJwt(c echo.Context) error {
	user := c.Get("user")
	token := user.(*jwt.Token)

	claims := token.Claims.(jwt.MapClaims)

	log.Println("User Name: ", claims["name"], "User ID: ", claims["jti"])

	return c.String(http.StatusOK, "jwt")
}

func CheckCookie(next echo.HandlerFunc) echo.HandlerFunc {
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
