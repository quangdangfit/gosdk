package main

import (
	"github.com/quangdangfit/gosdk/errors"
	"github.com/quangdangfit/gosdk/utils/logger"
	"github.com/quangdangfit/gosdk/validator"
)

func main() {
	type Car struct {
		Name  string `json:"name" validate:"required"`
		Brand string `json:"brand" validate:"required"`
	}

	car := Car{Name: "VinFast"}

	validate := validator.New()
	err := validate.Validate(car)

	if err != nil {
		logger.Error(errors.Stack(err))
	}
}
