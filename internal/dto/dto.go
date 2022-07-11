package dto

type User struct {
	UserID   int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdrawals struct {
	UserID      int     `json:"user_id"`
	Order       string  `json:"order"`
	Sum         float64 `json:"sum"`
	ProcessedAt string  `json:"processed_at,omitempty"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn float64 `json:"withdrawn"`
}

//const StatusNew = "NEW"

type Order struct {
	Number     string      `json:"order"`
	Status     OrderStatus `json:"status"`
	Accrual    float64     `json:"accrual"`
	UploadedAt string      `json:"uploaded_at,omitempty"`
}

type OrderStatus string

const (
	StatusNew        OrderStatus = "NEW"
	StatusRegistered OrderStatus = "REGISTERED"
	StatusInvalid    OrderStatus = "INVALID"
	StatusProcessing OrderStatus = "PROCESSING"
	StatusProcessed  OrderStatus = "PROCESSED"
)

type AccrualResponse struct {
	NumOrder    string      `json:"order"`
	OrderStatus OrderStatus `json:"status"`
	Accrual     float64     `json:"accrual"`
}
