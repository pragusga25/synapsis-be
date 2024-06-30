package validations

type CreateProductDto struct {
	Name        string  `json:"name" validate:"required,min=3,max=120"`
	Price       float64 `json:"price" validate:"required,gt=0"`
	Description string  `json:"description" validate:"required,min=3,max=1000"`
	Quantity    uint    `json:"quantity" validate:"required,gt=0"`
}

type UpdateProductDto struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=3,max=100"`
	Description *string  `json:"description,omitempty" validate:"omitempty,min=3,max=1000"`
	Price       *float64 `json:"price,omitempty" validate:"omitempty,gt=0"`
	Quantity    *uint    `json:"quantity,omitempty" validate:"omitempty,gt=0"`
}
