package utils

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"ivanmyagkov/gofermart/internal/config"
)

func CreateToken(login string, userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": login, "userID": userID})
	tokenString, _ := token.SignedString(config.TokenKey)
	return tokenString, nil
}

func CheckToken(tokenString string) (string, bool, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexected signing method: %v", token.Header["alg"])
		}
		return config.TokenKey, nil
	})
	if err != nil {
		return "", false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%d", claims["userID"]), ok, nil
	}
	return "", false, nil
}

func CreateCookie(c echo.Context, login string, userID int) error {

	var err error
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Value, err = CreateToken(login, userID)
	if err != nil {
		return err
	}
	cookie.Name = "token"
	c.SetCookie(cookie)
	c.Request().AddCookie(cookie)
	return nil
}

func GetUserID(c echo.Context) int {
	cookie, _ := c.Request().Cookie("token")
	user := cookie.Value
	token, _ := jwt.Parse(user, func(token *jwt.Token) (interface{}, error) {
		return config.TokenKey, nil
	})
	var claims = token.Claims.(jwt.MapClaims)
	return int(claims["userID"].(float64))

}
