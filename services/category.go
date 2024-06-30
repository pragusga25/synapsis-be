package services

import (
	"synapsis/models"
	"synapsis/repositories"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

type CategoryService interface {
	CreateCategory(category *models.Category) error
	GetCategoryByID(id uint) (*models.Category, error)
	GetAllCategories() ([]models.Category, error)
	DeleteCategoryByID(id uint) error
	ListProductsByCategoryID(categoryID uint) ([]models.ListProductResponse, error)
}

type categoryService struct {
	categoryRepository repositories.CategoryRepository
}

func NewCategoryService(categoryRepo repositories.CategoryRepository) CategoryService {
	return &categoryService{categoryRepository: categoryRepo}
}

func (s *categoryService) CreateCategory(category *models.Category) error {
	return s.categoryRepository.CreateCategory(category)
}

func (s *categoryService) GetCategoryByID(id uint) (*models.Category, error) {
	return s.categoryRepository.GetCategoryByID(id)
}

func (s *categoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepository.GetAllCategories()
}

func (s *categoryService) DeleteCategoryByID(id uint) error {
	return s.categoryRepository.DeleteCategoryByID(id)
}

func (s *categoryService) ListProductsByCategoryID(categoryID uint) ([]models.ListProductResponse, error) {
	products, err := s.categoryRepository.ListProductsByCategoryID(categoryID)
	if err != nil {
		return nil, utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
	}

	resP := make([]models.ListProductResponse, 0)
	for _, product := range products {
		resP = append(resP, utils.ProductModelToListProductResponse(&product))
	}

	return resP, nil
}
