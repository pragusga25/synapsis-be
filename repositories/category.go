package repositories

import (
	"synapsis/models"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CategoryRepository interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	DeleteCategoryByID(id uint) error
	ListProductsByCategoryID(categoryID uint) ([]models.Product, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) CreateCategory(category *models.Category) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) GetCategoryByID(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Products").First(&category, id).Error
	return &category, err
}

func (r *categoryRepository) GetAllCategories() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) DeleteCategoryByID(id uint) error {
	db := r.db.Delete(&models.Category{}, id)

	err := db.Error
	if err != nil {
		return err
	}

	if db.RowsAffected == 0 {
		return utils.WrapWithCustomeError(utils.ErrCategoryNotFound, fiber.StatusNotFound)
	}

	return nil
}

func (r *categoryRepository) ListProductsByCategoryID(categoryID uint) ([]models.Product, error) {
	var products []models.Product
	// err := r.db.Where("category_id = ?", categoryID).Find(&products).Error
	err := r.db.Joins("JOIN product_categories ON products.id = product_categories.product_id").
		Where("product_categories.category_id = ?", categoryID).Find(&products).Error
	return products, err
}
