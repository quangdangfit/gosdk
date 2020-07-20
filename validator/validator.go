package validator

import (
	govalidator "github.com/go-playground/validator/v10"

	"github.com/quangdangfit/gosdk/utils/logger"
)

type Validator interface {
	Validate(data interface{}, opts ...Option) error
}

type validator struct {
	engine *govalidator.Validate
}

func (v *validator) Validate(data interface{}, opts ...Option) error {
	opt := getOption(opts...)

	switch opt.typeOfObject {
	case Struct:
		return v.engine.Struct(data)
	case Field:
		return v.engine.Var(data, opt.tag)
	default:
		logger.Errorw("invalid type. not validate, error will nil")
		return nil
	}
}

func New() Validator {
	engine := govalidator.New()
	return &validator{engine: engine}
}
