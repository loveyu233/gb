package gb

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

var (
	validatorTrans ut.Translator
)

// InitValidator 初始化验证器
func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		panic("无法找到验证器")
	}

	// 注册翻译器
	var err error
	validatorTrans, err = registerTranslator(v)
	if err != nil {
		panic(err)
	}
	registerTagNameFunc(v)
	registerPhoneValidator(v)
	registerDecimalPlacesValidator(v)
}

// TranslateError 翻译错误信息，只返回第一个错误
func TranslateError(err error) error {
	switch typedErr := err.(type) {
	case *json.SyntaxError:
		return fmt.Errorf("JSON语法错误: %s", typedErr.Error())
	case *json.UnmarshalTypeError:
		return fmt.Errorf("参数类型错误: 字段 '%s' 应为 %s 类型", typedErr.Field, typedErr.Type)
	case validator.ValidationErrors:
		if len(typedErr) > 0 {
			return errors.New(typedErr[0].Translate(validatorTrans))
		}
	case *validator.InvalidValidationError:
		return typedErr

	case *strconv.NumError:
		return fmt.Errorf("参数类型解析错误: '%s' %s", typedErr.Num, typedErr.Err)
	}

	return err
}

func registerTagNameFunc(v *validator.Validate) {
	v.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}
		return name
	})
}

// registerPhoneValidator 注册手机号验证器
func registerPhoneValidator(v *validator.Validate) {
	v.RegisterValidation("phone", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return IsPhone(phone)
	})

	// 注册手机号翻译
	v.RegisterTranslation("phone", validatorTrans,
		// 注册翻译器
		func(ut ut.Translator) error {
			return ut.Add("phone", "手机号格式不正确", true)
		},
		// 自定义翻译函数
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("phone", fe.Field())
			return t
		},
	)
}

const (
	// 大陆手机号正则
	PhoneRegex = "^1[3-9]\\d{9}$"
)

// IsPhone 判断是否为大陆手机号
func IsPhone(phone string) bool {
	reg, err := regexp.Compile(PhoneRegex)
	if err != nil {
		return false
	}
	return reg.MatchString(phone)
}

// registerDecimalPlacesValidator 注册小数点位数验证器
func registerDecimalPlacesValidator(v *validator.Validate) {
	v.RegisterValidation("decimal_places", func(fl validator.FieldLevel) bool {
		param := fl.Param() // 获取参数值，如 "2"
		places, err := strconv.Atoi(param)
		if err != nil {
			return false
		}

		value := fl.Field().Float()
		multiplier := math.Pow10(places)
		return value == float64(int64(value*multiplier))/multiplier
	})

	// 注册翻译
	v.RegisterTranslation("decimal_places", validatorTrans,
		// 注册翻译器
		func(ut ut.Translator) error {
			return ut.Add("decimal_places", "{0}最多支持{1}位小数", true)
		},
		// 自定义翻译函数
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("decimal_places", fe.Field(), fe.Param())
			return t
		},
	)
}

// 注册翻译
func registerTranslator(v *validator.Validate) (trans ut.Translator, err error) {
	// 初始化中文翻译器
	zhTrans := zh.New()
	uni := ut.New(zhTrans, zhTrans)

	trans, found := uni.GetTranslator("zh")
	if !found {
		return nil, errors.New("无法找到中文翻译器")
	}
	// ValidatorTrans = trans

	// 注册默认的中文翻译
	if err := zhtranslations.RegisterDefaultTranslations(v, trans); err != nil {
		return nil, fmt.Errorf("注册默认翻译失败: %w", err)
	}

	// 注册 unique 标签的翻译
	v.RegisterTranslation("unique", trans, func(ut ut.Translator) error {
		return ut.Add("unique", "{0}不能包含重复值", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("unique", fe.Field())
		return t
	})

	return trans, nil
}
