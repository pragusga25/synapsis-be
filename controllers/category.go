package controllers

import (
	"errors"
	"synapsis/models"
	"synapsis/services"
	"synapsis/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgconn"
)

type CategoryController struct {
	categoryService services.CategoryService
}

func NewCategoryController(categoryService services.CategoryService) *CategoryController {
	return &CategoryController{categoryService: categoryService}
}

func (c *CategoryController) GetAllCategories(ctx *fiber.Ctx) error {
	categories, err := c.categoryService.GetAllCategories()
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, categories)
}

func (c *CategoryController) GetCategoryByID(ctx *fiber.Ctx) error {
	categoryID, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid category ID")
	}
	category, err := c.categoryService.GetCategoryByID(uint(categoryID))
	if err != nil {

		return utils.ErrorResponseJSON(ctx, fiber.StatusNotFound, "Category not found")
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, category)
}

func (c *CategoryController) CreateCategory(ctx *fiber.Ctx) error {
	category := new(models.Category)
	if err := ctx.BodyParser(category); err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, err.Error())
	}

	if err := c.categoryService.CreateCategory(category); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return utils.ErrorResponseJSON(ctx, fiber.StatusConflict, "Name already exists")
			}
		}

		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusCreated, category)
}

func (c *CategoryController) DeleteCategoryByID(ctx *fiber.Ctx) error {
	categoryID, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid category ID")
	}
	if err := c.categoryService.DeleteCategoryByID(uint(categoryID)); err != nil {
		cerr, ok := err.(*utils.CustomError)
		if ok {
			return utils.ErrorResponseJSON(ctx, cerr.Status, cerr.Inner.Error())
		}
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}

	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, nil)
}

func (c *CategoryController) ListProductsByCategoryID(ctx *fiber.Ctx) error {
	categoryID, err := ctx.ParamsInt("id")
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusBadRequest, "Invalid category ID")
	}
	products, err := c.categoryService.ListProductsByCategoryID(uint(categoryID))
	if err != nil {
		return utils.ErrorResponseJSON(ctx, fiber.StatusInternalServerError, utils.ErrDatabaseOperationFailed.Error())
	}
	return utils.SuccessResponseJSON(ctx, fiber.StatusOK, products)
}
