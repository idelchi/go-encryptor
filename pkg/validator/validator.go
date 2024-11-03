// Package validator provides a wrapper around go-playground/validator with
// built-in translation support for validation error messages. It simplifies
// struct validation and error handling by providing formatted, human-readable
// error messages in English.
package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
)

// ErrValidation is a sentinel error value used to identify validation errors.
var ErrValidation = errors.New("validation error")

// FieldLevel is a type alias for validator.FieldLevel.
type FieldLevel = validator.FieldLevel

// Validator wraps the go-playground/validator functionality with translation support.
// It holds both the validator instance and a translator for converting validation
// errors into human-readable messages.
type Validator struct {
	validate   *validator.Validate
	translator ut.Translator
}

// Validator returns the underlying validator instance.
// This method provides access to the raw validator when needed for advanced usage.
func (v *Validator) Validator() *validator.Validate {
	return v.validate
}

// NewValidator creates and initializes a new Validator instance with English translations.
// It sets up the universal translator with English locale and registers the default
// English translations for validation error messages.
func NewValidator() *Validator {
	en := en.New()
	uni := ut.New(en, en)
	trans, _ := uni.GetTranslator("en")

	validate := validator.New()
	if err := en_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		panic(err)
	}

	return &Validator{
		validate:   validate,
		translator: trans,
	}
}

// FormatErrors converts validator.ValidationErrors into a slice of standard errors.
// It translates the validation errors into human-readable messages using the
// configured translator. Returns nil if there are no errors to format.
func (v *Validator) FormatErrors(err error) []error {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		// If the error is not a ValidationErrors type, return it as a single error
		return []error{err}
	}

	if len(validationErrors) == 0 {
		return nil
	}

	errMap := validationErrors.Translate(v.translator)

	totalErrors := make([]error, 0, len(errMap))

	for _, err := range errMap {
		totalErrors = append(totalErrors, fmt.Errorf("%w: %s", ErrValidation, err))
	}

	return totalErrors
}

// Validate performs validation on the provided struct and returns any validation errors.
func (v *Validator) Validate(c any) []error {
	if err := v.validate.Struct(c); err != nil {
		return v.FormatErrors(err)
	}

	return nil
}

// RegisterValidationAndTranslation registers both a validation function and its error message translation.
// It simplifies the process of adding custom validations with proper error messages.
//
// Parameters:
//   - tag: the validation tag to use in struct field tags
//   - fn: the validation function that implements the validation logic
//   - msgTemplate: the error message template (use {0} for the field name and {1} for the parameter)
//
// Example usage:
//
//	validator.RegisterValidationAndTranslation(
//	    "multiple",
//	    validateMultiple,
//	    "{0} must be a multiple of {1}"
//	)
func (v *Validator) RegisterValidationAndTranslation(tag string, fn validator.Func, msgTemplate string) error {
	// Register the validation function
	if err := v.validate.RegisterValidation(tag, fn); err != nil {
		return fmt.Errorf("registering validation: %w", err)
	}

	// Register the translation
	if err := v.validate.RegisterTranslation(tag, v.translator,
		func(ut ut.Translator) error {
			if err := ut.Add(tag, msgTemplate, true); err != nil {
				return fmt.Errorf("adding translation: %w", err)
			}

			return nil
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			params := strings.Split(fe.Param(), "/")

			param1 := params[0]
			param2 := ""
			if len(params) > 1 {
				param2 = params[1]
			}

			t, _ := ut.T(tag, fe.Field(), param1, param2)

			return t
		},
	); err != nil {
		return fmt.Errorf("registering translation: %w", err)
	}

	return nil
}
