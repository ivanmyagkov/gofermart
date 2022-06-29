package storage

import (
	"context"
	"database/sql"
	"log"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"

	"ivanmyagkov/gofermart/internal/dto"
	"ivanmyagkov/gofermart/internal/interfaces"
	"ivanmyagkov/gofermart/internal/utils"
)

type Storage struct {
	db  *sql.DB
	ctx context.Context
}

func NewDB(psqlConn string, ctx context.Context) (*Storage, error) {
	db, err := sql.Open("postgres", psqlConn)
	if err != nil {
		return nil, interfaces.ErrDBConn
	}

	if err = db.Ping(); err != nil {
		return nil, interfaces.ErrPingDB
	}
	log.Println("Connected to DB!")
	if err = createTable(db); err != nil {
		log.Println(err)
		return nil, interfaces.ErrCreateTable
	}
	return &Storage{
		db:  db,
		ctx: ctx,
	}, nil
}

func createTable(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS users (
		id serial primary key,
		login text not null unique,
		password text not null,
        "current" float not null default 0,
        withdrawn int not null  default 0);
	`
	_, err := db.Exec(query)
	log.Println(err)
	if err != nil {
		return err
	}
	return nil
}

func (D *Storage) UserRegister(user dto.User) error {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (login, password) VALUES ($1, $2) `
	_, err = D.db.Exec(query, user.Login, hash)
	if err != nil {
		errCode := err.(*pq.Error).Code
		if pgerrcode.IsIntegrityConstraintViolation(string(errCode)) {
			return interfaces.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (D *Storage) UserLogin(user dto.User) error {
	var u dto.User
	query := `SELECT login,password FROM users WHERE login=$1`
	D.db.QueryRow(query, user.Login).Scan(&u.Login, &u.Password)
	check := utils.CheckPasswordHash(user.Password, u.Password)
	if !check {
		return interfaces.ErrBadPassword
	}
	return nil
}

func (D *Storage) Ping() error {
	return D.db.Ping()
}
func (D *Storage) Close() error {
	err := D.db.Close()
	return err
}
