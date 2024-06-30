package repositories

import (
	"synapsis/models"

	"gorm.io/gorm"
)

type CartRepository interface {
	AddToCart(cart *models.Cart) error
	GetCartByUserID(userID uint) ([]models.Cart, error)
	GetCartByUserIDAndProductID(userID, productID uint) (*models.Cart, error)
	UpdateCart(cart *models.Cart) error
	DeleteCart(cart *models.Cart) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) AddToCart(cart *models.Cart) error {
	return r.db.Create(cart).Error
}

func (r *cartRepository) GetCartByUserID(userID uint) ([]models.Cart, error) {
	var carts []models.Cart
	err := r.db.Where("user_id = ?", userID).Preload("Product").Find(&carts).Error
	return carts, err
}

func (r *cartRepository) GetCartByUserIDAndProductID(userID, productID uint) (*models.Cart, error) {
	var cart models.Cart
	err := r.db.Where("user_id = ? AND product_id = ?", userID, productID).First(&cart).Error
	return &cart, err
}

func (r *cartRepository) UpdateCart(cart *models.Cart) error {
	if cart.Quantity == 0 {
		return r.db.Delete(cart).Error
	}
	return r.db.Save(cart).Error
}

func (r *cartRepository) DeleteCart(cart *models.Cart) error {
	return r.db.Delete(cart).Error
}
