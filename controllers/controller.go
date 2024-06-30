package controllers

type Controllers struct {
	AuthController        *AuthController
	OrderController       *OrderController
	ProductController     *ProductController
	CartController        *CartController
	CategoryController    *CategoryController
	TransactionController *TransactionController
}
