package test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"ivanmyagkov/gofermart/internal/config"
	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/handlers"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/middlewares"
	"ivanmyagkov/gofermart/internal/storage"
)

func TestHandler_GetUserBalance(t *testing.T) {
	type args struct {
		db     *interfaces.DB
		cfg    *config.Config
		qu     chan dto.AccrualResponse
		cookie string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "body without token",
			args: args{
				qu:     make(chan dto.AccrualResponse, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "123454546565gdrrgr",
			},
			want: want{code: 401},
		},
		{
			name: "success",
			args: args{
				qu:     make(chan dto.AccrualResponse, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoieWFuMTIiLCJ1c2VySUQiOjF9.vrpDzuAw8sTMKQWFMPqM03oFrMAbFYx_h0G84-3jNi0",
			},
			want: want{code: 200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			e.Use(middlewares.SessionWithCookies)
			db, err := storage.NewDB(tt.args.cfg.GetDatabaseURI(), context.Background())
			if err != nil {
				log.Fatalf("Failed to create db %e", err)
			}
			s := handlers.New(db, tt.args.qu)
			e.GET("/api/user/balance", s.GetUserBalance)
			req := httptest.NewRequest(http.MethodGet, "/api/user/balance", nil)
			rec := httptest.NewRecorder()

			cookies := new(http.Cookie)
			cookies.Name = "token"
			cookies.Path = "/"
			cookies.Value = tt.args.cookie
			req.AddCookie(cookies)
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.want.code, rec.Code)

		})
	}
}
