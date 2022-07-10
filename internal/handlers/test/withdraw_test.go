package test

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"ivanmyagkov/gofermart/internal/config"
	"ivanmyagkov/gofermart/internal/handlers"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/middlewares"
	"ivanmyagkov/gofermart/internal/storage"
)

func TestHandler_PostUserBalanceWithdraw(t *testing.T) {
	type args struct {
		db     *interfaces.DB
		cfg    *config.Config
		qu     chan string
		cookie string
	}
	type want struct {
		code int
	}
	tests := []struct {
		name  string
		args  args
		value string
		want  want
	}{
		{
			name:  "body without token",
			value: `{"order": "2345","sum": 1}`,
			args: args{
				qu:     make(chan string, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "123454546565gdrrgr",
			},
			want: want{code: 401},
		},
		{
			name:  "wrong order",
			value: `{"order": "23451","sum": 1}`,
			args: args{
				qu:     make(chan string, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoieWFuMTIiLCJ1c2VySUQiOjF9.vrpDzuAw8sTMKQWFMPqM03oFrMAbFYx_h0G84-3jNi0",
			},
			want: want{code: 422},
		},
		{
			name:  "success",
			value: `{"order": "2345","sum": 1}`,
			args: args{
				qu:     make(chan string, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoieWFuMTIiLCJ1c2VySUQiOjF9.vrpDzuAw8sTMKQWFMPqM03oFrMAbFYx_h0G84-3jNi0",
			},
			want: want{code: 200},
		},
		{
			name:  "doesnt have any money",
			value: `{"order": "2345","sum": 1500}`,
			args: args{
				qu:     make(chan string, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoieWFuMTIiLCJ1c2VySUQiOjF9.vrpDzuAw8sTMKQWFMPqM03oFrMAbFYx_h0G84-3jNi0",
			},
			want: want{code: 402},
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
			e.POST("/api/user/balance/withdraw", s.PostUserBalanceWithdraw)
			req := httptest.NewRequest(http.MethodPost, "/api/user/balance/withdraw", strings.NewReader(tt.value))
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

func TestHandler_GetUserBalanceWithdrawals(t *testing.T) {
	type args struct {
		db     *interfaces.DB
		cfg    *config.Config
		qu     chan string
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
				qu:     make(chan string, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "123454546565gdrrgr",
			},
			want: want{code: 401},
		},
		{
			name: "not found",
			args: args{
				qu:     make(chan string, 100),
				cfg:    config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
				cookie: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyIjoieWFuMTIzIiwidXNlcklEIjo2fQ.YEpE706fNF6PoIqGBXh6345Z0WlrEjJf94jpB1VJgmI",
			},
			want: want{code: 204},
		},
		{
			name: "success",
			args: args{
				qu:     make(chan string, 100),
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
			e.GET("/api/user/withdrawals", s.GetUserBalanceWithdrawals)
			req := httptest.NewRequest(http.MethodGet, "/api/user/withdrawals", nil)
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
