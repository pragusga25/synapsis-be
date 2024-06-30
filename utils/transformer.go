package utils

import "synapsis/models"

func ProductModelToListProductResponse(product *models.Product) models.ListProductResponse {
	return models.ListProductResponse{
		ID:          product.ID,
		Name:        product.Name,
		Price:       product.Price,
		Quantity:    product.Quantity,
		Description: product.Description,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
	}
}

func ProductModelsToListProductResponses(products []models.Product) []models.ListProductResponse {

	listProductResponses := make([]models.ListProductResponse, 0, len(products))

	for _, product := range products {
		listProductResponses = append(listProductResponses, ProductModelToListProductResponse(&product))
	}

	return listProductResponses
}

func CartModelToListCartResponse(cart *models.Cart) models.ListCartResponse {
	return models.ListCartResponse{
		ID:        cart.ID,
		Quantity:  cart.Quantity,
		CreatedAt: cart.CreatedAt,
		UpdatedAt: cart.UpdatedAt,
		Product: models.ListCartProductResponse{
			ID:          cart.Product.ID,
			Name:        cart.Product.Name,
			Price:       cart.Product.Price,
			Description: cart.Product.Description,
		},
	}
}

func OrderModelToListOrderResponse(order *models.Order) models.ListOrderResponse {
	listOrderItemResponses := make([]models.ListOrderItemResponse, 0, len(order.OrderItems))

	for _, orderItem := range order.OrderItems {
		listOrderItemResponses = append(listOrderItemResponses, models.ListOrderItemResponse{
			ID:                 orderItem.ID,
			ProductName:        orderItem.ProductName,
			ProductPrice:       orderItem.ProductPrice,
			ProductDescription: orderItem.ProductDescription,
			Quantity:           orderItem.Quantity,
		})
	}

	return models.ListOrderResponse{
		ID:         order.ID,
		Status:     string(order.Status),
		TotalPrice: order.TotalPrice,
		OrderItems: listOrderItemResponses,
	}
}

func CartModelsToListCartResponses(carts []models.Cart) []models.ListCartResponse {
	listCartResponses := make([]models.ListCartResponse, 0, len(carts))

	for _, cart := range carts {
		listCartResponses = append(listCartResponses, CartModelToListCartResponse(&cart))
	}

	return listCartResponses
}

func OrderModelsToListOrderResponses(orders []models.Order) []models.ListOrderResponse {
	listOrderResponses := make([]models.ListOrderResponse, 0, len(orders))

	for _, order := range orders {
		listOrderResponses = append(listOrderResponses, OrderModelToListOrderResponse(&order))
	}

	return listOrderResponses
}
