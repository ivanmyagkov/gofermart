package dto

import "time"

type User struct {
	UserID   int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdrawals struct {
	UserID      int    `json:"user_id"`
	Order       string `json:"order"`
	Sum         int    `json:"sum"`
	ProcessedAt string `json:"processed_at,omitempty"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

const StatusNew = "NEW"

type Order struct {
	Number     string  `json:"number"`
	Status     string  `json:"status"`
	Accrual    float64 `json:"accrual"`
	UploadedAt time.Time
}

type AccrualResponse struct {
	NumOrder    string  `json:"order"`
	OrderStatus string  `json:"status"`
	Accrual     float64 `json:"accrual"`
}
