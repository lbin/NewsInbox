package request

import (
	"fmt"
	"reflect"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	enTranslations "github.com/go-playground/validator/v10/translations/en"
	zhTranslations "github.com/go-playground/validator/v10/translations/zh"
)

var Trans *ut.Translator

func TransInit(lang string) (err error) {
	if lang == "" {
		lang = "zh"
	}

	valid, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		return
	}

	zhT := zh.New() //chinese
	enT := en.New() //english
	uni := ut.New(enT, zhT, enT)

	trans, o := uni.GetTranslator(lang)
	if !o {
		return fmt.Errorf("uni.GetTranslator(%s) failed", lang)
	}

	valid.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("label")
		return name
	})

	// 注册翻译器
	switch lang {
	case "zh":
		err = zhTranslations.RegisterDefaultTranslations(valid, trans)
	default:
		err = enTranslations.RegisterDefaultTranslations(valid, trans)
	}

	Trans = &trans

	return
}
