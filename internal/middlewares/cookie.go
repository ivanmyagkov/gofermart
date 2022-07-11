package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"ivanmyagkov/gofermart/internal/utils"
)

func SessionWithCookies(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		var ok bool
		cookie, err := c.Cookie("token")
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}

		_, ok, err = utils.CheckToken(cookie.Value)
		if !ok {
			return echo.NewHTTPError(http.StatusUnauthorized)
		}
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
		return next(c)
	}
}
