package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"ivanmyagkov/gofermart/internal/utils"
)

func SessionWithCookies(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Path() == "/api/user/login" || c.Path() == "/api/user/register" {
			return next(c)
		}
		var ok bool

		cookie, err := c.Cookie("token")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		} else {
			_, ok, err = utils.CheckToken(cookie.Value)
			if err != nil {
				return echo.NewHTTPError(http.StatusInternalServerError)
			} else {
				if !ok {
					return echo.NewHTTPError(http.StatusUnauthorized)
				}
			}
		}

		return next(c)
	}
}
