package helper

import (
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTrans "github.com/go-playground/validator/v10/translations/en"
)

type Validator interface {
	Validate(i interface{}) error
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func NewValidator() Validator {
	langEn := en.New()
	uni := ut.New(langEn, langEn)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	_ = enTrans.RegisterDefaultTranslations(validate, trans)

	return &customValidator{
		validator: validate,
	}
}
