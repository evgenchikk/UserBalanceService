package models

type Users struct {
	UserID          int     `json:"user_id" binding:"required"`
	RealBalance     float64 `json:"real_balance" binding:"required"`
	ReservedBalance float64 `json:"reserved_balance"`
}
