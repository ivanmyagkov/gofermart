package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"ivanmyagkov/gofermart/internal/handlers"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/middlewares"
)

func InitSrv(db interfaces.DB) *echo.Echo {
	//server
	handler := handlers.New(db)

	//new Echo instance
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.Decompress())
	e.Use(middlewares.SessionWithCookies)
	e.POST("/api/user/register", handler.PostUserRegister)
	e.POST("/api/user/login", handler.PostUserLogin)
	e.POST("/api/user/orders", handler.PostUserOrders)
	e.GET("/api/user/orders", handler.GetUserOrders)
	e.GET("/api/user/balance", handler.GetUserBalance)
	e.POST("/api/user/balance/withdraw", handler.PostUserBalanceWithdraw)
	e.GET("/api/user/balance/withdrawals", handler.GetUserBalanceWithdrawals)

	return e
}
