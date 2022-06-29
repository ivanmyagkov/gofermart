package interfaces

import (
	"errors"

	"ivanmyagkov/gofermart/internal/dto"
)

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrBadPassword   = errors.New("wrong password")
	ErrDBConn        = errors.New("DB connection error")
	ErrCreateTable   = errors.New("create tables error")
	ErrPingDB        = errors.New("ping Db error")
	ErrWasDeleted    = errors.New("was deleted")
)

type DB interface {
	UserRegister(user dto.User) error
	UserLogin(user dto.User) error
	Ping() error
	Close() error
}
