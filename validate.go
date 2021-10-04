package ktools

import (
	"fmt"
	"strings"

	// "github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"

	// en_translations "github.com/go-playground/validator/v10/translations/en"
	zh_tran "github.com/go-playground/validator/v10/translations/zh"
)

// use a single instance , it caches struct info
var (
	uni      *ut.UniversalTranslator
	validate *validator.Validate
)

func Trans() {

	// NOTE: ommitting allot of error checking for brevity

	// en := en.New()
	// uni = ut.New(en, en)
	zh := zh.New()
	uni = ut.New(zh, zh)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("zh")

	validate = validator.New()
	zh_tran.RegisterDefaultTranslations(validate, trans)

	translateAll(trans)
	translateIndividual(trans)
	translateOverride(trans) // yep you can specify your own in whatever locale you want!
}

func TranVali(ruleData interface{}, zhMap map[string]string) (res map[string]string) {

	zh := zh.New()
	uni = ut.New(zh, zh)

	// this is usually know or extracted from http 'Accept-Language' header
	// also see uni.FindTranslator(...)
	trans, _ := uni.GetTranslator("zh")

	validate = validator.New()
	zh_tran.RegisterDefaultTranslations(validate, trans)

	err := validate.Struct(ruleData)

	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		res = errs.Translate(trans)

		res = tranAttr(res, zhMap)

		return
	}

	return

}

func tranAttr(errMap map[string]string, zhMap map[string]string) (tranMap map[string]string) {

	for key, errString := range errMap {

		attrKey := strings.Split(key, ".")[1]

		val, has := zhMap[attrKey]

		if has {
			tran := strings.Replace(errString, attrKey, val, 1)

			errMap[key] = tran
		}
	}

	tranMap = errMap

	return
}

///==================== validate example =====================

func translateAll(trans ut.Translator) {

	type User struct {
		Username string `validate:"required"`
		Tagline  string `validate:"required,lt=10"`
		Tagline2 string `validate:"required,gt=1"`
	}

	user := User{
		Username: "Joeybloggs",
		Tagline:  "This tagline is way too long.",
		Tagline2: "1",
	}

	err := validate.Struct(user)
	if err != nil {

		// translate all error at once
		errs := err.(validator.ValidationErrors)

		// returns a map with key = namespace & value = translated error
		// NOTICE: 2 errors are returned and you'll see something surprising
		// translations are i18n aware!!!!
		// eg. '10 characters' vs '1 character'
		fmt.Println(errs.Translate(trans))
	}
}

func translateIndividual(trans ut.Translator) {

	type User struct {
		Username string `validate:"required"`
		Some     int    `validate:"required,gt=6"`
	}

	var user User

	err := validate.Struct(user)
	if err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			fmt.Println(e.Translate(trans))
		}
	}
}

func translateOverride(trans ut.Translator) {

	validate.RegisterTranslation("required", trans, func(ut ut.Translator) error {
		return ut.Add("required", "{0} 重写：必须填写!", true) // see universal-translator for details
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", fe.Field())

		return t
	})

	type User struct {
		Username string `validate:"required"`
	}

	var user User

	err := validate.Struct(user)
	if err != nil {

		errs := err.(validator.ValidationErrors)

		for _, e := range errs {
			// can translate each error one at a time.
			fmt.Println(e.Translate(trans))
		}
	}
}
