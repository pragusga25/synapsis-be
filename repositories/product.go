package repositories

import (
	"synapsis/models"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type ProductRepository interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uint) (*models.Product, error)
	GetProductsByCategoryID(categoryID uint) ([]models.Product, error)
	GetAllProducts() ([]models.Product, error)
	DeleteProductByID(id uint) error
	UpdateProductByID(productID uint, updates map[string]interface{}) error
	AddCategoryToProduct(product *models.Product, category *models.Category) error
	RemoveCategoryFromProduct(product *models.Product, category *models.Category) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) CreateProduct(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) GetProductByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.Preload("Categories").First(&product, id).Error

	return &product, err
}

func (r *productRepository) GetProductsByCategoryID(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	err := r.db.Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Where("product_categories.category_id = ?", categoryID).Find(&products).Error
	return products, err
}

func (r *productRepository) GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error

	return products, err
}

func (r *productRepository) DeleteProductByID(id uint) error {
	return r.db.Delete(&models.Product{}, id).Error
}

func (r *productRepository) UpdateProductByID(productID uint, updates map[string]interface{}) error {
	db := r.db.Model(&models.Product{}).Where("id = ?", productID).Updates(updates)

	err := db.Error
	if err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return utils.WrapWithCustomeError(utils.ErrProductNotFound, fiber.StatusNotFound)
	}

	return nil
}

func (r *productRepository) AddCategoryToProduct(product *models.Product, category *models.Category) error {
	product.Categories = append(product.Categories, *category)
	return r.db.Save(product).Error
}

func (r *productRepository) RemoveCategoryFromProduct(product *models.Product, category *models.Category) error {
	return r.db.Model(product).Association("Categories").Delete(category)
}
