package repositories

import (
	"synapsis/models"

	"gorm.io/gorm"
)

type OrderRepository interface {
	CreateOrder(order *models.Order) error
	CreateOrderItem(orderItem *models.OrderItem) error
	GetOrdersByUserID(userID uint) ([]models.Order, error)
	GetOrderByID(orderID string) (*models.Order, error)
	UpdateOrder(order *models.Order) error
	GetOrderByIDAndUserID(orderID string, userID uint) (*models.Order, error)
}

type orderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) OrderRepository {
	return &orderRepository{db: db}
}

func (r *orderRepository) CreateOrder(order *models.Order) error {
	return r.db.Create(order).Error
}

func (r *orderRepository) CreateOrderItem(orderItem *models.OrderItem) error {
	return r.db.Create(orderItem).Error
}

func (r *orderRepository) GetOrdersByUserID(userID uint) ([]models.Order, error) {
	var orders []models.Order
	err := r.db.Where("user_id = ?", userID).Preload("OrderItems").Find(&orders).Error
	return orders, err
}

func (r *orderRepository) GetOrderByID(orderID string) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ?", orderID).Preload("OrderItems").First(&order).Error
	return &order, err
}

func (r *orderRepository) UpdateOrder(order *models.Order) error {
	return r.db.Save(order).Error
}

func (r *orderRepository) GetOrderByIDAndUserID(orderID string, userID uint) (*models.Order, error) {
	var order models.Order
	err := r.db.Where("id = ? AND user_id = ?", orderID, userID).Preload("OrderItems").First(&order).Error
	return &order, err
}
