package services

import (
	"synapsis/models"
	"synapsis/repositories"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type CartService interface {
	AddToCart(userID uint, productID uint, quantity int) error
	GetCartByUserID(userID uint) ([]models.ListCartResponse, error)
	UpdateCart(userID uint, productID uint, quantity int) error
	DeleteCart(cart *models.Cart) error
}

type cartService struct {
	cartRepository    repositories.CartRepository
	productRepository repositories.ProductRepository
}

func NewCartService(cartRepo repositories.CartRepository, productRepo repositories.ProductRepository) CartService {
	return &cartService{cartRepository: cartRepo, productRepository: productRepo}
}

func (s *cartService) AddToCart(userID uint, productID uint, quantity int) error {

	product, err := s.productRepository.GetProductByID(productID)

	if err != nil {
		return utils.WrapWithCustomeError(utils.ErrProductNotFound, fiber.StatusNotFound)
	}

	if quantity > int(product.Quantity) {
		return utils.WrapWithCustomeError(utils.ErrProductQuantityExceeded, fiber.StatusBadRequest)
	}

	cart, err := s.cartRepository.GetCartByUserIDAndProductID(userID, productID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return s.cartRepository.AddToCart(&models.Cart{UserID: userID, ProductID: productID, Quantity: quantity})
		}

		return utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
	}

	if quantity+cart.Quantity > int(product.Quantity) {
		return utils.WrapWithCustomeError(utils.ErrProductQuantityExceeded, fiber.StatusBadRequest)
	}

	cart.Quantity += quantity
	return s.cartRepository.UpdateCart(cart)
}

func (s *cartService) GetCartByUserID(userID uint) ([]models.ListCartResponse, error) {
	cart, err := s.cartRepository.GetCartByUserID(userID)

	if err != nil {
		return nil, utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
	}

	return utils.CartModelsToListCartResponses(cart), nil
}

func (s *cartService) UpdateCart(userID uint, productID uint, quantity int) error {
	cart, err := s.cartRepository.GetCartByUserIDAndProductID(userID, productID)

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return utils.WrapWithCustomeError(utils.ErrProductNotFound, fiber.StatusNotFound)
		}

		return utils.WrapWithCustomeError(utils.ErrDatabaseOperationFailed, fiber.StatusInternalServerError)
	}

	product, err := s.productRepository.GetProductByID(productID)

	if err != nil {
		return utils.WrapWithCustomeError(utils.ErrProductNotFound, fiber.StatusNotFound)
	}

	if quantity > int(product.Quantity) {
		return utils.WrapWithCustomeError(utils.ErrProductQuantityExceeded, fiber.StatusBadRequest)
	}

	cart.Quantity = quantity
	return s.cartRepository.UpdateCart(cart)
}

func (s *cartService) DeleteCart(cart *models.Cart) error {
	return s.cartRepository.DeleteCart(cart)
}
