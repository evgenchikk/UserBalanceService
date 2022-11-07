package api

type DereserveSchemaJSON struct {
	OrderID int `json:"order_id" binding:"required,min=1" example:"1"`
}
