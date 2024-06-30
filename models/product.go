package models

import (
	"time"
)

type Product struct {
	ID          uint       `gorm:"primarykey"`
	Name        string     `gorm:"type:varchar(120);not null" json:",omitempty"` //omitempty
	Description string     `gorm:"type:text;not null" json:",omitempty"`
	Price       float64    `gorm:"type:decimal(10,2)" json:",omitempty"`
	Quantity    uint       `gorm:"not null" json:",omitempty"`
	Categories  []Category `gorm:"many2many:product_categories;"`

	CreatedAt time.Time `json:",omitempty"`
	UpdatedAt time.Time `json:",omitempty"`
}

type ProductCategory struct {
	CategoryID uint `gorm:"primaryKey"`
	Category   Category
	ProductID  uint `gorm:"primaryKey"`
	Product    Product
	CreatedAt  time.Time
}

type ListProductResponse struct {
	ID          uint
	Name        string
	Description string
	Price       float64
	Quantity    uint
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type GetProductCategoryResponse struct {
	Name string
	ID   uint
}

type GetProductResponse struct {
	ListProductResponse
	Categories []GetProductCategoryResponse
}
