package api

type TransferSchemaJSON struct {
	FromUserID int     `json:"from_user_id" binding:"required,min=1" example:"1"`
	ToUserID   int     `json:"to_user_id" binding:"required,min=1" example:"2"`
	Amount     float64 `json:"amount" binding:"required,min=0" example:"100"`
}
