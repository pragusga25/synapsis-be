package routes

import (
	"synapsis/controllers"
	"synapsis/middleware"
	"synapsis/utils"
	"synapsis/validations"

	"github.com/gofiber/fiber/v2"
)

func Setup(app *fiber.App, ctrl *controllers.Controllers,
) {
	app.Get("/", func(c *fiber.Ctx) error {
		return utils.SuccessResponseJSON(c, fiber.StatusOK, &fiber.Map{"message": "Welcome to Real World!"})
	})
	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/register", middleware.ValidateBody(&validations.RegisterDto{}), ctrl.AuthController.Register)
	auth.Post("/login", middleware.ValidateBody(&validations.LoginDto{}), ctrl.AuthController.Login)

	product := api.Group("/products")
	product.Get("/", ctrl.ProductController.GetAllProducts)
	product.Get("/:id", ctrl.ProductController.GetProductByID)
	product.Post("/", middleware.ValidateBody(&validations.CreateProductDto{}), middleware.AuthenticateUser, middleware.AdminOnly, ctrl.ProductController.CreateProduct)
	product.Put("/:id", middleware.ValidateBody(&validations.UpdateProductDto{}), middleware.AuthenticateUser, middleware.AdminOnly, ctrl.ProductController.UpdateProductByID)
	product.Put("/:product_id/categories/:category_id", middleware.AuthenticateUser, middleware.AdminOnly, ctrl.ProductController.AddCategoryToProduct)
	product.Delete("/:product_id/categories/:category_id", middleware.AuthenticateUser, middleware.AdminOnly, ctrl.ProductController.RemoveCategoryFromProduct)

	category := api.Group("/categories")
	category.Get("/", ctrl.CategoryController.GetAllCategories)
	category.Get("/:id", ctrl.CategoryController.GetCategoryByID)
	category.Post("/", middleware.ValidateBody(&validations.CreateCategoryDto{}), middleware.AuthenticateUser, middleware.AdminOnly, ctrl.CategoryController.CreateCategory)
	category.Delete("/:id", middleware.AuthenticateUser, middleware.AdminOnly, ctrl.CategoryController.DeleteCategoryByID)
	category.Get("/:id/products", ctrl.CategoryController.ListProductsByCategoryID)

	cart := api.Group("/cart")
	cart.Get("", middleware.AuthenticateUser, ctrl.CartController.GetCart)
	cart.Post("", middleware.ValidateBody(&validations.AddToCartDto{}), middleware.AuthenticateUser, ctrl.CartController.AddToCart)
	cart.Put("", middleware.ValidateBody(&validations.UpdateCartDto{}), middleware.AuthenticateUser, ctrl.CartController.UpdateCart)

	order := api.Group("/orders")
	order.Post("/checkout", middleware.AuthenticateUser, ctrl.OrderController.Checkout)
	order.Post("/confirm", middleware.AuthenticateUser, middleware.AdminOnly, ctrl.OrderController.ConfirmOrder)
	order.Get("", middleware.AuthenticateUser, ctrl.OrderController.GetOrders)
	order.Get("/:id", middleware.AuthenticateUser, ctrl.OrderController.GetOrder)

	trx := api.Group("/transactions")
	trx.Post("", ctrl.TransactionController.CreateTransaction)
	trx.Get("", middleware.AuthenticateUser, ctrl.TransactionController.GetTransactions)
}
