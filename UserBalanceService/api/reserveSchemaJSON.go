package api

type ReserveSchemaJSON struct {
	UserID    int     `json:"user_id" binding:"required,min=1" example:"1"`
	ServiceID int     `json:"service_id" binding:"required,min=1" example:"1"`
	OrderID   int     `json:"order_id" binding:"required,min=1" example:"1"`
	Price     float64 `json:"price" binding:"required,min=0" example:"100"`
}
