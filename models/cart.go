package models

import "time"

type Cart struct {
	ID        uint `gorm:"primarykey"`
	UserID    uint `gorm:"uniqueIndex:idx_user_product"`
	ProductID uint `gorm:"uniqueIndex:idx_user_product"`
	Product   Product
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ListCartProductResponse struct {
	ID          uint
	Name        string
	Price       float64
	Description string
}

type ListCartResponse struct {
	ID        uint
	Quantity  int
	CreatedAt time.Time
	UpdatedAt time.Time
	Product   ListCartProductResponse
}
