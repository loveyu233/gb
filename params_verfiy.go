package gb

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zhtranslations "github.com/go-playground/validator/v10/translations/zh"
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
	registerIDCarValidator(v)
	registerGUIDValidator(v)
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
		return ValidateChineseMobile(phone)
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

// registerIDCarValidator 注册身份证号验证器
func registerIDCarValidator(v *validator.Validate) {
	v.RegisterValidation("idcar", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return ValidateChineseIDCard(phone)
	})

	// 注册手机号翻译
	v.RegisterTranslation("idcar", validatorTrans,
		// 注册翻译器
		func(ut ut.Translator) error {
			return ut.Add("idcar", "身份证号格式不正确", true)
		},
		// 自定义翻译函数
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("idcar", fe.Field())
			return t
		},
	)
}

func registerGUIDValidator(v *validator.Validate) {
	v.RegisterValidation("guid", func(fl validator.FieldLevel) bool {
		phone := fl.Field().String()
		return ValidateCustomGUIDRegex(phone)
	})

	// 注册手机号翻译
	v.RegisterTranslation("guid", validatorTrans,
		// 注册翻译器
		func(ut ut.Translator) error {
			return ut.Add("guid", "guid格式不正确", true)
		},
		// 自定义翻译函数
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("guid", fe.Field())
			return t
		},
	)
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

// CreateRequiredError 创建一个key不能为空的验证错误
func CreateRequiredError(key string) error {
	fieldError := &mockFieldError{
		tag:   "required",
		field: key,
		param: "",
	}

	validationErrors := validator.ValidationErrors{fieldError}
	return validationErrors
}

// CreateTypeError 创建一个参数类型错误
func CreateTypeError(key, value string, originalErr error) error {
	fieldError := &mockFieldError{
		tag:   "type",
		field: key,
		param: value,
		err:   originalErr,
	}

	validationErrors := validator.ValidationErrors{fieldError}
	return validationErrors
}

// mockFieldError 模拟 validator.FieldError 接口（扩展版本）
type mockFieldError struct {
	tag   string
	field string
	param string
	err   error // 添加原始错误
}

func (m *mockFieldError) Tag() string             { return m.tag }
func (m *mockFieldError) ActualTag() string       { return m.tag }
func (m *mockFieldError) Namespace() string       { return m.field }
func (m *mockFieldError) StructNamespace() string { return m.field }
func (m *mockFieldError) Field() string           { return m.field }
func (m *mockFieldError) StructField() string     { return m.field }
func (m *mockFieldError) Value() interface{}      { return m.param }
func (m *mockFieldError) Param() string           { return m.param }
func (m *mockFieldError) Kind() reflect.Kind      { return reflect.String }
func (m *mockFieldError) Type() reflect.Type      { return reflect.TypeOf("") }
func (m *mockFieldError) Error() string {
	if m.tag == "type" {
		return fmt.Sprintf("%s type conversion failed", m.field)
	}
	return fmt.Sprintf("%s is %s", m.field, m.tag)
}

// Translate 实现翻译方法，兼容现有的翻译器
func (m *mockFieldError) Translate(trans ut.Translator) string {
	switch m.tag {
	case "required":
		// 尝试使用翻译器翻译，如果失败则使用默认中文
		if trans != nil {
			if t, err := trans.T("required", m.field); err == nil {
				return t
			}
		}
		return fmt.Sprintf("%s不能为空", m.field)

	case "type":
		// 尝试使用翻译器翻译 type 标签
		if trans != nil {
			if t, err := trans.T("type", m.field); err == nil {
				return t
			}
		}

		// 根据原始错误类型返回更具体的错误信息
		if m.err != nil {
			switch m.err.(type) {
			case *strconv.NumError:
				return fmt.Sprintf("%s必须是有效的数字", m.field)
			default:
				if strings.Contains(m.err.Error(), "bool") {
					return fmt.Sprintf("%s必须是有效的布尔值(true/false)", m.field)
				}
			}
		}
		return fmt.Sprintf("%s参数格式错误", m.field)

	default:
		return fmt.Sprintf("%s验证失败", m.field)
	}
}

// convertToType 将字符串转换为指定类型
func convertToType[T any](value string) (T, error) {
	var result any
	var err error
	var zero T

	// 使用类型断言确定目标类型
	switch any(zero).(type) {
	case string:
		result = value

	case int:
		result, err = strconv.Atoi(value)

	case int8:
		temp, parseErr := strconv.ParseInt(value, 10, 8)
		result, err = int8(temp), parseErr

	case int16:
		temp, parseErr := strconv.ParseInt(value, 10, 16)
		result, err = int16(temp), parseErr

	case int32:
		temp, parseErr := strconv.ParseInt(value, 10, 32)
		result, err = int32(temp), parseErr

	case int64:
		result, err = strconv.ParseInt(value, 10, 64)

	case uint:
		temp, parseErr := strconv.ParseUint(value, 10, 0)
		result, err = uint(temp), parseErr

	case uint8:
		temp, parseErr := strconv.ParseUint(value, 10, 8)
		result, err = uint8(temp), parseErr

	case uint16:
		temp, parseErr := strconv.ParseUint(value, 10, 16)
		result, err = uint16(temp), parseErr

	case uint32:
		temp, parseErr := strconv.ParseUint(value, 10, 32)
		result, err = uint32(temp), parseErr

	case uint64:
		result, err = strconv.ParseUint(value, 10, 64)

	case float32:
		temp, parseErr := strconv.ParseFloat(value, 32)
		result, err = float32(temp), parseErr

	case float64:
		result, err = strconv.ParseFloat(value, 64)

	case bool:
		result, err = strconv.ParseBool(value)

	case []string:
		// 支持逗号分隔的字符串数组
		result = strings.Split(value, ",")

	case []int:
		// 支持逗号分隔的整数数组
		parts := strings.Split(value, ",")
		intSlice := make([]int, len(parts))
		for i, part := range parts {
			intSlice[i], err = strconv.Atoi(strings.TrimSpace(part))
			if err != nil {
				return zero, err
			}
		}
		result = intSlice

	default:
		// 不支持的类型，尝试直接断言
		if converted, ok := any(value).(T); ok {
			result = converted
		} else {
			return zero, fmt.Errorf("不支持的类型转换: %T", zero)
		}
	}

	if err != nil {
		return zero, err
	}

	return result.(T), nil
}
