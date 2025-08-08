package domain

import "time"

type OrderStatus string

const (
	StatusPending   OrderStatus = "pending"
	StatusCompleted OrderStatus = "completed"
	StatusCancelled OrderStatus = "cancelled"
)

type OrderItem struct {
	ProductID  string `json:"product_id" bson:"product_id"`
	Quantity   int    `json:"quantity" bson:"quantity"`
	PriceCents int64  `json:"price_cents" bson:"price_cents"` // snapshot
}

type Order struct {
	ID         string      `json:"id" bson:"_id,omitempty"`
	UserID     string      `json:"user_id" bson:"user_id"`
	Items      []OrderItem `json:"items" bson:"items"`
	TotalCents int64       `json:"total_cents" bson:"total_cents"`
	Status     OrderStatus `json:"status" bson:"status"`
	CreatedAt  time.Time   `json:"created_at" bson:"created_at"`
	UpdatedAt  time.Time   `json:"updated_at" bson:"updated_at"`
}

type CreateOrderRequest struct {
	UserID        string      `json:"user_id" binding:"required"`
	Items         []OrderItem `json:"items" binding:"required"`
	PaymentMethod string      `json:"payment_method"`
}
