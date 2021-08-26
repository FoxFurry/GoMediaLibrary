package validators

import (
	"fmt"
	"github.com/foxfurry/simple-rest/internal/common/server/common_translators"
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
	emptyFieldMsg = "cannot be empty"
)

var(
	invalidYearMsg = fmt.Sprintf("Year should be between -868 and %v", time.Now().Year())

	FieldTitleEmpty = common_translators.FieldError{
		Field: "Title",
		Msg:   "Title " + emptyFieldMsg,
	}
	FieldAuthorEmpty = common_translators.FieldError{
		Field: "Author",
		Msg:   "Author " + emptyFieldMsg,
	}
	FieldYearEmpty = common_translators.FieldError{
		Field: "Year",
		Msg:   "Year " + emptyFieldMsg,
	}
	FieldYearInvalid = common_translators.FieldError{
		Field: "Year",
		Msg:   fmt.Sprintf("Year should be between -868 and %v", time.Now().Year()),
	}
)

var validID validator.Func = func(fl validator.FieldLevel) bool {
	id := fl.Field().Int()
	if id < 1 {
		return false
	}
	return true
}

var validYear validator.Func = func(fl validator.FieldLevel) bool {
	year := fl.Field().Int()
	if year > int64(time.Now().Year()) || year < -868 {
		return false
	}
	return true
}

var trslValidYear validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(yearTag, invalidYearMsg, true)
}

var trslValidID validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(idTag, "{0} should be positive non-null number", true)
}

var requiredMessage validator.RegisterTranslationsFunc = func(ut ut.Translator) error {
	return ut.Add(requiredTag, "{0} " +emptyFieldMsg, true)
}

func RegisterBookValidators() {
	errTranslator := common_translators.GetTranslator()

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
		log.Panicf("Could not register common_translators: %v", ok)
	}
}
