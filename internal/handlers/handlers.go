package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/utils"
)

type Handler struct {
	db interfaces.DB
}

func New(db interfaces.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) PostUserRegister(c echo.Context) error {

	var user dto.User

	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if user.Login == "" || user.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = h.db.UserRegister(user); err != nil {
		if errors.Is(err, interfaces.ErrAlreadyExists) {
			return c.NoContent(http.StatusConflict)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	if err = utils.CreateCookie(c, user.Login); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserLogin(c echo.Context) error {
	var user dto.User
	body, err := io.ReadAll(c.Request().Body)
	if err != nil || len(body) == 0 {
		return c.NoContent(http.StatusBadRequest)
	}

	err = json.Unmarshal(body, &user)

	if err != nil {
		return c.NoContent(http.StatusBadRequest)
	}

	if user.Login == "" || user.Password == "" {
		return c.NoContent(http.StatusBadRequest)
	}

	if err = h.db.UserLogin(user); err != nil {
		log.Println(err)
		if errors.Is(err, interfaces.ErrBadPassword) {
			return c.NoContent(http.StatusUnauthorized)
		}
		return c.NoContent(http.StatusInternalServerError)
	}
	if err = utils.CreateCookie(c, user.Login); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)

}

func (h *Handler) PostUserOrders(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserOrders(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserBalance(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) PostUserBalanceWithdraw(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}

func (h *Handler) GetUserBalanceWithdrawals(c echo.Context) error {
	return c.NoContent(http.StatusOK)
}
