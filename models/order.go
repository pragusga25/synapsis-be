package models

import (
	"time"

	"gorm.io/gorm"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusConfirmed OrderStatus = "confirmed"
	OrderStatusFailed    OrderStatus = "failed"
)

func (r OrderStatus) IsValid() bool {
	switch r {
	case OrderStatusPending, OrderStatusConfirmed:
		return true
	}
	return false
}

func (r OrderStatus) String() string {
	return string(r)
}

type Order struct {
	ID         string `gorm:"primarykey;type:varchar(80);not null"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
	UserID     uint
	User       User
	TotalPrice float64
	Status     OrderStatus

	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	OrderID   string
	Order     Order
	ProductID uint
	Product   Product

	// Create the snapshot of the product at the time of order
	// to keep the product price at the time of order
	ProductName        string
	ProductPrice       float64
	ProductDescription string
	Quantity           int
}

type ListOrderItemResponse struct {
	ID                 uint
	ProductName        string
	ProductPrice       float64
	ProductDescription string
	Quantity           int
}

type ListOrderResponse struct {
	ID         string
	Status     string
	TotalPrice float64
	OrderItems []ListOrderItemResponse
}
