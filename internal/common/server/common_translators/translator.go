package common_translators

import (
	"fmt"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"sync"
)

var translatorSingleton ut.Translator
var once sync.Once

type FieldError struct {
	Field string `json:"field"`
	Msg   string `json:"msg"`
}

func (f FieldError) Error() string {
	return fmt.Sprintf("field: %v, msg: %v", f.Field, f.Msg)
}

func CreateFieldError(field string, msg string) FieldError {
	return FieldError{
		Field: field,
		Msg:   msg,
	}
}

func GetTranslator() ut.Translator {
	once.Do(func() {
		trsEntity := en.New()
		uni := ut.New(trsEntity, trsEntity)

		var ok bool
		translatorSingleton, ok = uni.GetTranslator("en")
		if !ok {
			log.Fatalf("Translator not found")
		}
	})

	return translatorSingleton
}

func Translate(validatorError error) []FieldError {
	translator := GetTranslator()

	errArray, ok := validatorError.(validator.ValidationErrors)
	if !ok {
		log.Fatalf("Could not translate: %v", ok)
	}
	var result []FieldError

	for _, err := range errArray {
		result = append(result, FieldError{
			Field: err.Field(),
			Msg:   err.Translate(translator),
		})
	}
	return result
}
