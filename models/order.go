package models

type CreateOrderRequest struct {
	ProductName string `json:"product_name" binding:"required"`
	Quantity    int    `json:"quantity" binding:"required"`
	UserID      string `json:"user_id"`
}
