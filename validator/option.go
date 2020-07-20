package validator

type Option interface {
	apply(*option)
}

type Type int

const (
	Struct Type = 1
	Field  Type = 2
)

type option struct {
	typeOfObject Type
	tag          string
}

type optionFn func(*option)

func (optFn optionFn) apply(opt *option) {
	optFn(opt)
}

// WithType set type of object will validate, default type is Struct
func WithType(t Type) Option {
	return optionFn(func(opt *option) {
		opt.typeOfObject = t
	})
}

// WithTag will available with `type = Field`
func WithTag(tag string) Option {
	return optionFn(func(opt *option) {
		opt.tag = tag
	})
}

func getOption(opts ...Option) option {
	opt := option{
		typeOfObject: Struct,
		tag:          "",
	}

	for _, o := range opts {
		o.apply(&opt)
	}

	return opt
}
