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
	ErrWrongOrder    = errors.New("wrong order number")
	ErrMoney         = errors.New("no enough money")
	ErrOtherUser     = errors.New("other user have this order")
)

type DB interface {
	UserRegister(user *dto.User) error
	UserLogin(user *dto.User) error
	SaveOrder(number string, userID int) (dto.Order, error)
	GetOrders(userID int) ([]dto.Order, error)
	UserBalance(userID int) (dto.Balance, error)
	BalanceWithdraw(userID int, withdraw dto.Withdrawals) error
	GetUserWithdrawals(userID int) ([]dto.Withdrawals, error)
	UpdateAccrualOrder(ac dto.Order) error
	SelectNewOrders() ([]dto.Order, error)
	Ping() error
	Close() error
}
