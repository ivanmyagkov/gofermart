package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/handlers"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/middlewares"
)

func InitSrv(db interfaces.DB, qu chan dto.Order) *echo.Echo {
	//server
	handler := handlers.New(db, qu)

	//new Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())

	r := e.Group("")
	r.Use(middlewares.SessionWithCookies)
	e.POST("/api/user/register", handler.PostUserRegister)
	e.POST("/api/user/login", handler.PostUserLogin)
	r.POST("/api/user/orders", handler.PostUserOrders)
	r.GET("/api/user/orders", handler.GetUserOrders)
	r.GET("/api/user/balance", handler.GetUserBalance)
	r.POST("/api/user/balance/withdraw", handler.PostUserBalanceWithdraw)
	r.GET("/api/user/withdrawals", handler.GetUserBalanceWithdrawals)

	return e
}
