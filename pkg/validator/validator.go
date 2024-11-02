package validator

import (
	"errors"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

func (v *Validator) Validator() *validator.Validate {
	return v.validate
}

func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	en_translations.RegisterDefaultTranslations(validate, trans)

	return &Validator{
		validate:   validate,
		translator: trans,
	}
}

func (v *Validator) FormatErrors(err error) []error {
	errs := err.(validator.ValidationErrors)

	if len(errs) == 0 {
		return nil
	}

	errMap := errs.Translate(v.translator)

	totalErrors := make([]error, 0, len(errMap))

	for _, err := range errMap {
		totalErrors = append(totalErrors, errors.New(err))
	}

	return totalErrors
}

func (v *Validator) Validate(c any) []error {
	if err := v.validate.Struct(c); err != nil {
		return v.FormatErrors(err)
	}

	return nil
}
