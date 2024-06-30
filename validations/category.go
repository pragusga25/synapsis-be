package validations

type CreateCategoryDto struct {
	Name string `json:"name" validate:"required,min=3,max=32"`
}
