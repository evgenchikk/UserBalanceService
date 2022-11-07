package api

type BalanceSchemaJSON struct {
	UserID int `json:"user_id" binding:"required,min=1" example:"1"`
}
