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
	"github.com/stretchr/testify/require"

	"ivanmyagkov/gofermart/internal/config"
	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/handlers"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/storage"
)

func TestHandler_PostUserRegister(t *testing.T) {
	type args struct {
		db  *interfaces.DB
		cfg *config.Config
		qu  chan dto.AccrualResponse
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
			name:  "body is empty",
			value: "",
			args: args{
				qu:  make(chan dto.AccrualResponse, 100),
				cfg: config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
			},
			want: want{code: 400},
		},
		{
			name:  "already exists",
			value: `{"login":"yan12", "password":"yan12"}`,
			args: args{
				qu:  make(chan dto.AccrualResponse, 100),
				cfg: config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
			},
			want: want{code: 409},
		},
		{
			name:  "success",
			value: `{"login":"yan123", "password":"yan123"}`,
			args: args{
				qu:  make(chan dto.AccrualResponse, 100),
				cfg: config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
			},
			want: want{code: 200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			db, err := storage.NewDB(tt.args.cfg.GetDatabaseURI(), context.Background())
			if err != nil {
				log.Fatalf("Failed to create db %e", err)
			}
			s := handlers.New(db, tt.args.qu)

			req := httptest.NewRequest(http.MethodPost, "/api/user/register", strings.NewReader(tt.value))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := s.PostUserRegister(c)
			if assert.NoError(t, h) {
				require.Equal(t, tt.want.code, rec.Code)
			}
		})
	}
}

func TestHandler_PostUserLogin(t *testing.T) {
	type args struct {
		db  *interfaces.DB
		cfg *config.Config
		qu  chan dto.AccrualResponse
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
			name:  "body is empty",
			value: "",
			args: args{
				qu:  make(chan dto.AccrualResponse, 100),
				cfg: config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
			},
			want: want{code: 400},
		},
		{
			name:  "login or password is wrong",
			value: `{"login":"yan12", "password":"1234"}`,
			args: args{
				qu:  make(chan dto.AccrualResponse, 100),
				cfg: config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
			},
			want: want{code: 401},
		},
		{
			name:  "success",
			value: `{"login":"yan123", "password":"yan123"}`,
			args: args{
				qu:  make(chan dto.AccrualResponse, 100),
				cfg: config.NewConfig(":8080", "postgres://ivanmyagkov@localhost:5432/postgres?sslmode=disable", "http://localhost:8080"),
			},
			want: want{code: 200},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			db, err := storage.NewDB(tt.args.cfg.GetDatabaseURI(), context.Background())
			if err != nil {
				log.Fatalf("Failed to create db %e", err)
			}
			s := handlers.New(db, tt.args.qu)

			req := httptest.NewRequest(http.MethodPost, "/api/user/login", strings.NewReader(tt.value))
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			h := s.PostUserLogin(c)
			if assert.NoError(t, h) {
				require.Equal(t, tt.want.code, rec.Code)
			}
		})
	}
}
