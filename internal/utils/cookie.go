package utils

import (
	"fmt"
	"log"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"ivanmyagkov/gofermart/internal/config"
)

func CreateToken(login string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"user": login})
	tokenString, _ := token.SignedString(config.TokenKey)
	log.Println(tokenString)
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
		log.Println(err)
		return "", false, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%s", claims["user"]), ok, nil
	}
	return "", false, nil
}

func CreateCookie(c echo.Context, login string) error {

	var err error
	cookie := new(http.Cookie)
	cookie.Path = "/"
	cookie.Value, err = CreateToken(login)

	if err != nil {
		return err
	}
	cookie.Name = "token"
	c.SetCookie(cookie)
	c.Request().AddCookie(cookie)
	return nil
}
