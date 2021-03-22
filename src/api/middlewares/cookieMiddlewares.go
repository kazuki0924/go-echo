package middlewares

import (
	"log"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func checkCookie(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("sessionID")
		if err != nil {
			if strings.Contains(err.Error(), "named cookie not present") {
				return c.String(http.StatusUnauthorized, "you don't have any cookie")
			}

			log.Println(err)
			return err
		}

		if cookie.Value == "secret" {
			return next(c)
		}

		return c.String(http.StatusUnauthorized, "you don't have the right cookie, cookie")
	}
}

func SetCookieMiddlewares(g *echo.Group) {
	g.Use(checkCookie)
}
