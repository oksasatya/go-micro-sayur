package validator

import (
	"errors"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
)

type Validator struct {
	Validator  *validator.Validate
	Translator ut.Translator
}

func NewValidator() *Validator {
	en := en.New()        // Create a new instance of the validator
	uni := ut.New(en, en) // Create a universal translator with the English locale
	trans, found := uni.GetTranslator("en")
	if !found {
		log.Fatalf("[NewValidator-1] NewValidator: translator not found")
	}

	validate := validator.New()

	return &Validator{
		Validator:  validate,
		Translator: trans, // Initialize translator if needed
	}
}

func (v *Validator) Validate(i interface{}) error {
	err := v.Validator.Struct(i)
	if err != nil {
		var object validator.ValidationErrors
		_ = errors.As(err, &object)
		for _, e := range object {
			log.Errorf("[Validator-1] Validate: %v", e.Translate(v.Translator))
			return errors.New(e.Translate(v.Translator))
		}
	}
	return nil
}
