package models

import "time"

type Orders struct {
	OrderID     int       `json:"order_id"`
	UserID      int       `json:"user_id"`
	ServiceID   int       `json:"service_id"`
	ActivatedAt time.Time `json:"activated_at"`
	Price       float64   `json:"price"`
}
