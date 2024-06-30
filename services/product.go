package services

import (
	"synapsis/models"
	"synapsis/repositories"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

type ProductService interface {
	CreateProduct(product *models.Product) error
	GetProductByID(id uint) (*models.GetProductResponse, error)
	GetProductsByCategoryID(categoryID uint) ([]models.ListProductResponse, error)
	GetAllProducts() ([]models.ListProductResponse, error)
	AddCategoryToProduct(productID, categoryID uint) error
	RemoveCategoryFromProduct(productID, categoryID uint) error
	UpdateProduct(productID uint, updates map[string]interface{}) error
}

type productService struct {
	productRepository  repositories.ProductRepository
	categoryRepository repositories.CategoryRepository
}

func NewProductService(productRepo repositories.ProductRepository, categoryRepo repositories.CategoryRepository) ProductService {
	return &productService{productRepository: productRepo, categoryRepository: categoryRepo}
}

func (s *productService) CreateProduct(product *models.Product) error {
	return s.productRepository.CreateProduct(product)
}

func (s *productService) GetProductByID(id uint) (*models.GetProductResponse, error) {
	product, err := s.productRepository.GetProductByID(id)

	if err != nil {
		return nil, utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
	}

	resProduct := models.GetProductResponse{
		ListProductResponse: utils.ProductModelToListProductResponse(product),
	}

	resProduct.Categories = make([]models.GetProductCategoryResponse, 0)
	if len(product.Categories) > 0 {
		for _, category := range product.Categories {
			resProduct.Categories = append(resProduct.Categories, models.GetProductCategoryResponse{
				ID:   category.ID,
				Name: category.Name,
			})
		}
	}

	return &resProduct, nil
}

func (s *productService) GetProductsByCategoryID(categoryID uint) ([]models.ListProductResponse, error) {
	products, err := s.productRepository.GetProductsByCategoryID(categoryID)
	if err != nil {
		return nil, utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
	}

	resP := make([]models.ListProductResponse, 0)
	for _, product := range products {
		resP = append(resP, utils.ProductModelToListProductResponse(&product))
	}
	return resP, nil
}

func (s *productService) GetAllProducts() ([]models.ListProductResponse, error) {
	products, err := s.productRepository.GetAllProducts()

	if err != nil {
		return nil, utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
	}

	resP := make([]models.ListProductResponse, 0)
	for _, product := range products {
		resP = append(resP, utils.ProductModelToListProductResponse(&product))
	}
	return resP, nil
}

func (s *productService) AddCategoryToProduct(productID, categoryID uint) error {
	category, err := s.categoryRepository.GetCategoryByID(categoryID)

	if err != nil {
		return utils.WrapWithCustomeError(utils.ErrCategoryNotFound, fiber.StatusNotFound)
	}

	product, err := s.productRepository.GetProductByID(productID)
	if err != nil {
		return utils.WrapWithCustomeError(utils.ErrProductNotFound, fiber.StatusNotFound)
	}

	return s.productRepository.AddCategoryToProduct(product, category)

}

func (s *productService) RemoveCategoryFromProduct(productID, categoryID uint) error {
	category, err := s.categoryRepository.GetCategoryByID(categoryID)

	if err != nil {
		return utils.WrapWithCustomeError(utils.ErrCategoryNotFound, fiber.StatusNotFound)
	}

	product, err := s.productRepository.GetProductByID(productID)
	if err != nil {
		return utils.WrapWithCustomeError(utils.ErrProductNotFound, fiber.StatusNotFound)
	}

	return s.productRepository.RemoveCategoryFromProduct(product, category)

}

func (s *productService) UpdateProduct(productID uint, updates map[string]interface{}) error {
	return s.productRepository.UpdateProductByID(productID, updates)
}
