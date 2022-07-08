package storage

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"time"

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
        "current" float not null default 0 CHECK ("current">=0),
        withdrawn float not null  default 0
        );
		CREATE TABLE IF NOT EXISTS orders (
		    "number" text primary key unique,
		    user_id int not null references users(id),
		    status text not null,
		    accrual float,
		    uploaded_at timestamp
		);
		CREATE TABLE IF NOT EXISTS withdrawals (
	    	user_id int not null references users(id),
			"number" text not null unique,
			"sum" float not null,
			processed_at timestamp
	);

		
	`
	_, err := db.Exec(query)
	if err != nil {
		return err
	}
	return nil
}
func (D *Storage) SelectNewOrders() ([]string, error) {
	query := `SELECT "number" from orders where status !=$1 and status !=$2`
	rows, err := D.db.Query(query, dto.StatusInvalid, dto.StatusProcessed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}
	var orders []string
	var order string
	for rows.Next() {
		if err = rows.Scan(&order); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	return orders, nil
}
func (D *Storage) UserRegister(user *dto.User) error {
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		return err
	}
	query := `INSERT INTO users (login, password) VALUES ($1, $2) RETURNING id`
	err = D.db.QueryRow(query, user.Login, hash).Scan(&user.UserID)
	if err != nil {
		errCode := err.(*pq.Error).Code
		if pgerrcode.IsIntegrityConstraintViolation(string(errCode)) {
			return interfaces.ErrAlreadyExists
		}
		return err
	}
	return nil
}

func (D *Storage) UserLogin(user *dto.User) error {
	var u dto.User
	query := `SELECT id,password FROM users WHERE login=$1`
	err := D.db.QueryRow(query, user.Login).Scan(&user.UserID, &u.Password)
	if err != nil {
		return err
	}
	check := utils.CheckPasswordHash(user.Password, u.Password)
	if !check {
		return interfaces.ErrBadPassword
	}
	return nil
}

func (D *Storage) SaveOrder(number string, userID int) error {
	var order dto.Order
	order.Number = number
	order.Status = dto.StatusNew
	order.Accrual = 0
	time := time.Now().Format(time.RFC3339)

	insertQuery := `INSERT INTO orders (number, user_id, status, accrual, uploaded_at) VALUES ($1,$2,$3,$4,$5)`
	_, err := D.db.Exec(insertQuery, order.Number, userID, order.Status, order.Accrual, time)

	if err != nil {
		errCode := err.(*pq.Error).Code
		if pgerrcode.IsIntegrityConstraintViolation(string(errCode)) {
			var user int
			selectOrder := `SELECT user_id FROM orders WHERE number=$1`
			err = D.db.QueryRow(selectOrder, number).Scan(&user)
			if err != nil {
				return err
			}

			if user == userID {
				return interfaces.ErrAlreadyExists
			}
			return interfaces.ErrOtherUser
		}
	}
	return nil
}

func (D *Storage) GetOrders(userID int) ([]dto.Order, error) {
	var order dto.Order
	var ordersArr []dto.Order
	query := `SELECT number, status, accrual, uploaded_at FROM orders WHERE user_id = $1`
	rows, err := D.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(&order.Number, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}
		ordersArr = append(ordersArr, order)
	}
	if len(ordersArr) == 0 {
		return nil, interfaces.ErrNotFound
	}
	return ordersArr, nil
}

func (D *Storage) UserBalance(userID int) (dto.Balance, error) {
	var balance dto.Balance
	query := `SELECT "current",withdrawn FROM users WHERE id=$1`
	err := D.db.QueryRow(query, userID).Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		return balance, err
	}
	return balance, nil
}
func (D *Storage) BalanceWithdraw(userID int, withdraw dto.Withdrawals) error {
	var money float64
	var check bool
	query := `SELECT "current" FROM users WHERE id=$1`
	err := D.db.QueryRow(query, userID).Scan(&money)
	if err != nil {
		return err
	}
	if money < withdraw.Sum {
		return interfaces.ErrMoney
	}
	orderQuerty := `SELECT true FROM withdrawals WHERE "number"=$1 and user_id=$2`
	err = D.db.QueryRow(orderQuerty, withdraw.Order, userID).Scan(&check)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tx, err := D.db.Begin()
			if err != nil {
				return err
			}

			insetrQuery := `INSERT INTO withdrawals (user_id, "number", sum, processed_at) VALUES ($1,$2,$3,$4)`
			_, err = tx.Exec(insetrQuery, userID, withdraw.Order, withdraw.Sum, withdraw.ProcessedAt)
			if err != nil {
				return err
			}
			updateBalance := `update users set current = "current"-$1,withdrawn = withdrawn+$1 where id= $2 `
			_, err = tx.Exec(updateBalance, withdraw.Sum, userID)
			if err != nil {
				return err
			}
			return nil
		}
		return err
	}

	return interfaces.ErrWrongOrder
}
func (D *Storage) UpdateAccrualOrder(ac dto.AccrualResponse) error {
	var userID int
	tx, err := D.db.Begin()
	if err != nil {
		return err
	}
	query := `UPDATE orders SET status = $1, accrual = $2 WHERE number = $3 RETURNING user_id`
	err = tx.QueryRow(query, ac.OrderStatus, ac.Accrual, ac.NumOrder).Scan(&userID)
	if err != nil {
		return err
	}
	update := `UPDATE users SET current = current + $1 WHERE id = $2`
	_, err = tx.Exec(update, ac.Accrual, userID)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (D *Storage) GetUserWithdrawals(userID int) ([]dto.Withdrawals, error) {
	var withdrawalsArr []dto.Withdrawals
	var withdrawl dto.Withdrawals
	query := `select w.number, w.sum, w.processed_at from withdrawals w left join users u on u.id = w.user_id where u.id=$1`
	rows, err := D.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for rows.Next() {
		if err = rows.Scan(&withdrawl.Order, &withdrawl.Sum, &withdrawl.ProcessedAt); err != nil {
			return nil, err
		}
		withdrawalsArr = append(withdrawalsArr, withdrawl)
	}
	if len(withdrawalsArr) == 0 {
		return nil, interfaces.ErrNotFound
	}

	return withdrawalsArr, nil
}

func (D *Storage) Ping() error {
	return D.db.Ping()
}
func (D *Storage) Close() error {
	err := D.db.Close()
	return err
}
