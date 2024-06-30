package controllers

import (
	"synapsis/models"
	"synapsis/services"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
)

type ProductController struct {
	productService services.ProductService
}

func NewProductController(ProductService services.ProductService) *ProductController {
	return &ProductController{productService: ProductService}
}

func (c *ProductController) GetAllProducts(ctx *fiber.Ctx) error {
	products, err := c.productService.GetAllProducts()
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, products)
}

func (c *ProductController) GetProductByID(ctx *fiber.Ctx) error {
	productID, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid product ID")
	}
	product, err := c.productService.GetProductByID(uint(productID))
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusNotFound, "Product not found")
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, product)

}

func (c *ProductController) GetProductsByCategoryID(ctx *fiber.Ctx) error {
	categoryID, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid category ID")
	}
	products, err := c.productService.GetProductsByCategoryID(uint(categoryID))
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, products)
}

func (c *ProductController) CreateProduct(ctx *fiber.Ctx) error {
	product := new(models.Product)
	if err := ctx.BodyParser(product); err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	if err := c.productService.CreateProduct(product); err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	product.Categories = []models.Category{}

	return utils.SuccessResponseJSON(ctx, fiber.StatusCreated, product)
}

func (c *ProductController) UpdateProductByID(ctx *fiber.Ctx) error {
	productID, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	updates := make(map[string]interface{})
	if err := ctx.BodyParser(&updates); err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	if err := c.productService.UpdateProduct(uint(productID), updates); err != nil {
		cerr, ok := err.(*utils.CustomError)
		if ok {
			return utils.ErrorResponseJSON(ctx, cerr.Status, cerr.Inner.Error())
		}
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)
}

func (c *ProductController) AddCategoryToProduct(ctx *fiber.Ctx) error {
	productID, err := ctx.ParamsInt("product_id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid product ID")
	}

	categoryID, err := ctx.ParamsInt("category_id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid category ID")
	}

	err = c.productService.AddCategoryToProduct(uint(productID), uint(categoryID))
	if err != nil {
		cerr, ok := err.(*utils.CustomError)
		if ok {
			return utils.ErrorResponseJSON(ctx, cerr.Status, cerr.Error())
		}

		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)
}

func (c *ProductController) RemoveCategoryFromProduct(ctx *fiber.Ctx) error {
	productID, err := ctx.ParamsInt("product_id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid product ID")
	}

	categoryID, err := ctx.ParamsInt("category_id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid category ID")
	}

	err = c.productService.RemoveCategoryFromProduct(uint(productID), uint(categoryID))
	if err != nil {
		cerr, ok := err.(*utils.CustomError)
		if ok {
			return utils.ErrorResponseJSON(ctx, cerr.Status, cerr.Error())
		}

		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)
}
