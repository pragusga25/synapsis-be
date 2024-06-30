package main

import (
	"fmt"
	"synapsis/config"
	"synapsis/controllers"
	"synapsis/database"
	"synapsis/middleware"
	"synapsis/repositories"
	"synapsis/routes"
	"synapsis/services"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		panic("Failed to load configuration")
	}

	// Initialize database
	database.InitDB(cfg)

	// Initialize repositories
	userRepo := repositories.NewUserRepository(database.DB)
	orderRepo := repositories.NewOrderRepository(database.DB)
	cartRepo := repositories.NewCartRepository(database.DB)
	productRepo := repositories.NewProductRepository(database.DB)
	categoryRepo := repositories.NewCategoryRepository(database.DB)
	transactionRepo := repositories.NewTransactionRepository(database.DB)

	// Initialize services
	authService := services.NewAuthService(userRepo)
	orderService := services.NewOrderService(orderRepo, cartRepo)
	productService := services.NewProductService(productRepo, categoryRepo)
	cartService := services.NewCartService(cartRepo, productRepo)
	categoryService := services.NewCategoryService(categoryRepo)
	midtransService := services.NewMidtransService(orderRepo, transactionRepo)

	transactionService := services.NewTransactionService(transactionRepo, orderRepo)

	// Initialize controllers
	authController := controllers.NewAuthController(authService)
	orderController := controllers.NewOrderController(orderService, midtransService)
	productController := controllers.NewProductController(productService)
	cartController := controllers.NewCartController(cartService)
	categoryController := controllers.NewCategoryController(categoryService)

	transactionController := controllers.NewTransactionController(transactionService, midtransService)

	ctrl := &controllers.Controllers{
		AuthController:        authController,
		OrderController:       orderController,
		ProductController:     productController,
		CartController:        cartController,
		CategoryController:    categoryController,
		TransactionController: transactionController,
	}

	// Setup Fiber
	app := fiber.New()

	// Add middleware
	app.Use(middleware.ToSnakeCaseMiddleware)

	// Setup routes
	routes.Setup(app, ctrl)

	// Start server
	app.Listen(fmt.Sprintf(":%d", cfg.Port))
}
