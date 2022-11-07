package api

type AddSchemaJSON struct {
	UserID int     `json:"user_id" binding:"required,min=1" example:"1"`
	Amount float64 `json:"amount" binding:"required,min=0" example:"100"`
}
