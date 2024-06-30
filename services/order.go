package services

import (
	"synapsis/models"
	"synapsis/repositories"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

type OrderService interface {
	Checkout(userID uint) (*models.Order, error)
	GetOrdersByUserID(userID uint) ([]models.ListOrderResponse, error)
	GetOrderByIDAndUserID(userID uint, orderID string) (*models.Order, error)
	ConfirmOrder(orderID string) error
}

type orderService struct {
	orderRepository repositories.OrderRepository
	cartRepository  repositories.CartRepository
}

func NewOrderService(orderRepo repositories.OrderRepository, cartRepo repositories.CartRepository) OrderService {
	return &orderService{orderRepository: orderRepo, cartRepository: cartRepo}
}

func (s *orderService) Checkout(userID uint) (*models.Order, error) {
	carts, err := s.cartRepository.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}

	if len(carts) == 0 {
		return nil, utils.WrapWithCustomeError(utils.ErrEmptyCart, fiber.StatusBadRequest)
	}

	var totalPrice float64
	for _, cart := range carts {
		totalPrice += cart.Product.Price * float64(cart.Quantity)
	}

	order := &models.Order{
		UserID:     userID,
		TotalPrice: totalPrice,
		Status:     models.OrderStatusPending,
		ID:         utils.GenerateID(16),
	}

	err = s.orderRepository.CreateOrder(order)
	if err != nil {
		return nil, err
	}

	for _, cart := range carts {
		orderItem := &models.OrderItem{
			OrderID:            order.ID,
			ProductID:          cart.ProductID,
			ProductName:        cart.Product.Name,
			ProductPrice:       cart.Product.Price,
			ProductDescription: cart.Product.Description,
			Quantity:           cart.Quantity,
		}
		err := s.orderRepository.CreateOrderItem(orderItem)
		if err != nil {
			return nil, err
		}
		// Remove item from cart after order is created
		s.cartRepository.DeleteCart(&cart)
	}

	return order, nil

}

func (s *orderService) GetOrdersByUserID(userID uint) ([]models.ListOrderResponse, error) {
	orders, err := s.orderRepository.GetOrdersByUserID(userID)

	if err != nil {
		return nil, utils.WrapWithCustomeError(utils.ErrOrderNotFound, fiber.StatusNotFound)
	}

	return utils.OrderModelsToListOrderResponses(orders), nil
}

func (s *orderService) GetOrderByIDAndUserID(userID uint, orderID string) (*models.Order, error) {
	return s.orderRepository.GetOrderByIDAndUserID(orderID, userID)
}

func (s *orderService) ConfirmOrder(orderID string) error {
	order, err := s.orderRepository.GetOrderByID(orderID)
	if err != nil {
		return err
	}

	order.Status = models.OrderStatusConfirmed
	return s.orderRepository.UpdateOrder(order)
}
