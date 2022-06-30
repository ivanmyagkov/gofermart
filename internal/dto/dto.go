package dto

type User struct {
	UserID   int    `json:"user_id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type Withdrawals struct {
	Order       string `json:"order"`
	Sum         int    `json:"sum"`
	ProcessedAt string `json:"processed_at,omitempty"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}
