package validators

import (
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

const(
	titleTag = "validtitle"
	authorTag = "validauthor"
	yearTag = "validyear"
	requiredTag = "required"
)

var validYear validator.Func = func(fl validator.FieldLevel) bool {
	year := fl.Field().Int()
	if year > int64(time.Now().Year()) || year < -868 {
		return false
	}
	return true
}

var trslValidYear validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(yearTag, fmt.Sprintf("{0} should be between -868 and %v", time.Now().Year()), true)
}

var requiredMessage validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(requiredTag, "{0} cannot be empty", true)
}

var errTranslator ut.Translator

func RegisterBookValidators(){
	trsEntity := en.New()
	uni := ut.New(trsEntity, trsEntity)

	var ok bool
	errTranslator, ok = uni.GetTranslator("en")
	if !ok {
		log.Fatalf("Translator not found")
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(yearTag, validYear)

		v.RegisterTranslation(yearTag, errTranslator, trslValidYear, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(yearTag, fe.Field())
			return t
		})

		v.RegisterTranslation(requiredTag, errTranslator, requiredMessage, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(requiredTag, fe.Field())
			return t
		})


	}else {
		log.Panicf("Could not register validators: %v", ok)
	}
}

func Translate(err error) string {
	errArray, ok := err.(validator.ValidationErrors)
	if !ok {
		return err.Error()
	}
	var result string

	for i, e := range errArray {
		if i != 0 {
			result += ", "
		}
		result += e.Translate(errTranslator)
	}
	return result
}