package validators

import (
	"fmt"
	"github.com/foxfurry/medialib/internal/common/server/translator"
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"log"
	"time"
)

const (
	idTag       = "validID"
	yearTag     = "validYear"
	requiredTag = "required"
)

var(
	emptyFieldMsg = "cannot be empty"
	invalidIDMsg = "should be positive non-null number"
	invalidYearMsg = fmt.Sprintf("Year should be between -868 and %v", time.Now().Year())

	FieldTitleEmpty = translator.FieldError{
		Field: "Title",
		Msg:   "Title " + emptyFieldMsg,
	}
	FieldAuthorEmpty = translator.FieldError{
		Field: "Author",
		Msg:   "Author " + emptyFieldMsg,
	}
	FieldYearEmpty = translator.FieldError{
		Field: "Year",
		Msg:   "Year " + emptyFieldMsg,
	}
	FieldYearInvalid = translator.FieldError{
		Field: "Year",
		Msg:   invalidYearMsg,
	}
)

var validID validator.Func = func(fl validator.FieldLevel) bool {
	id := fl.Field().Int()
	if id >= 1 {
		return true
	}
	return false
}

var validYear validator.Func = func(fl validator.FieldLevel) bool {
	year := fl.Field().Int()
	if -868 <= year && year <= int64(time.Now().Year()) {
		return true
	}
	return false
}

var trslValidYear validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(yearTag, invalidYearMsg, true)
}

var trslValidID validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(idTag, "{0} " + invalidIDMsg, true)
}

var requiredMessage validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(requiredTag, "{0} " + emptyFieldMsg, true)
}

func RegisterBookValidators() {
	errTranslator := translator.GetInstance()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation(yearTag, validYear)
		v.RegisterTranslation(yearTag, errTranslator, trslValidYear, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(yearTag, fe.Field())
			return t
		})

		v.RegisterValidation(idTag, validID)
		v.RegisterTranslation(idTag, errTranslator, trslValidID, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(yearTag, fe.Field())
			return t
		})

		v.RegisterTranslation(requiredTag, errTranslator, requiredMessage, func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T(requiredTag, fe.Field())
			return t
		})
	} else {
		log.Panicf("Could not register translator: %v", ok)
	}
}
