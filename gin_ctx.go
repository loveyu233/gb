package gb

import (
	"encoding/json"
	"errors"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// GetGinContextValue 从gin.Context中获取指定类型的值,可以是基础类型,也可以是结构体
func GetGinContextValue[T any](c *gin.Context, key string) (T, bool) {
	var zero T

	value, exists := c.Get(key)
	if !exists {
		return zero, false
	}

	// 如果值为nil，返回零值
	if value == nil {
		return zero, true
	}

	return ConvertAnyType[T](value)
}

func ConvertAnyType[T any](value any) (T, bool) {
	var zero T
	// 获取目标类型
	targetType := reflect.TypeOf(zero)
	// 处理基础类型
	switch any(zero).(type) {
	case string:
		return any(cast.ToString(value)).(T), true
	case int:
		return any(cast.ToInt(value)).(T), true
	case int64:
		return any(cast.ToInt64(value)).(T), true
	case int32:
		return any(cast.ToInt32(value)).(T), true
	case int16:
		return any(cast.ToInt16(value)).(T), true
	case int8:
		return any(cast.ToInt8(value)).(T), true
	case uint:
		return any(cast.ToUint(value)).(T), true
	case uint64:
		return any(cast.ToUint64(value)).(T), true
	case uint32:
		return any(cast.ToUint32(value)).(T), true
	case uint16:
		return any(cast.ToUint16(value)).(T), true
	case uint8:
		return any(cast.ToUint8(value)).(T), true
	case float64:
		return any(cast.ToFloat64(value)).(T), true
	case float32:
		return any(cast.ToFloat32(value)).(T), true
	case bool:
		return any(cast.ToBool(value)).(T), true
	default:
		// 处理结构体和其他复杂类型
		return ConvertToStruct[T](value, targetType)
	}
}

// ConvertToStruct 将值转换为结构体类型
func ConvertToStruct[T any](value interface{}, targetType reflect.Type) (T, bool) {
	var zero T

	// 获取value的反射值和类型
	valueReflect := reflect.ValueOf(value)
	valueType := reflect.TypeOf(value)

	// 如果类型完全匹配，直接返回
	if valueType == targetType {
		return value.(T), true
	}

	// 处理指针类型
	if targetType.Kind() == reflect.Ptr {
		if valueType.Kind() == reflect.Ptr {
			// 两个都是指针，递归处理元素类型
			if valueReflect.IsNil() {
				return zero, true
			}
			elemResult, ok := ConvertToStruct[T](valueReflect.Elem().Interface(), targetType.Elem())
			if !ok {
				return zero, false
			}
			// 创建指针
			result := reflect.New(targetType.Elem())
			result.Elem().Set(reflect.ValueOf(elemResult).Elem())
			return result.Interface().(T), true
		} else {
			// value不是指针，target是指针
			elemResult, ok := ConvertToStruct[T](value, targetType.Elem())
			if !ok {
				return zero, false
			}
			result := reflect.New(targetType.Elem())
			result.Elem().Set(reflect.ValueOf(elemResult))
			return result.Interface().(T), true
		}
	}

	// 如果value是指针，获取其元素
	if valueType.Kind() == reflect.Ptr {
		if valueReflect.IsNil() {
			return zero, true
		}
		return ConvertToStruct[T](valueReflect.Elem().Interface(), targetType)
	}

	// 处理slice类型
	if targetType.Kind() == reflect.Slice && valueType.Kind() == reflect.Slice {
		return convertSlice[T](value, targetType)
	}

	// 处理map类型
	if targetType.Kind() == reflect.Map && valueType.Kind() == reflect.Map {
		return convertMap[T](value, targetType)
	}

	// 尝试JSON序列化/反序列化转换（适用于结构体）
	if targetType.Kind() == reflect.Struct ||
		(targetType.Kind() == reflect.Slice && targetType.Elem().Kind() == reflect.Struct) ||
		(targetType.Kind() == reflect.Map && (targetType.Key().Kind() == reflect.String || targetType.Elem().Kind() == reflect.Struct)) {

		// 先序列化为JSON
		jsonData, err := json.Marshal(value)
		if err != nil {
			return zero, false
		}

		// 创建目标类型的实例
		result := reflect.New(targetType).Interface()

		// 反序列化到目标类型
		err = json.Unmarshal(jsonData, result)
		if err != nil {
			return zero, false
		}

		// 如果目标类型不是指针，获取其值
		if targetType.Kind() != reflect.Ptr {
			return reflect.ValueOf(result).Elem().Interface().(T), true
		}

		return result.(T), true
	}

	// 尝试reflect转换
	if valueReflect.Type().ConvertibleTo(targetType) {
		converted := valueReflect.Convert(targetType)
		return converted.Interface().(T), true
	}

	return zero, false
}

// convertSlice 处理slice类型转换
func convertSlice[T any](value interface{}, targetType reflect.Type) (T, bool) {
	var zero T

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Kind() != reflect.Slice {
		return zero, false
	}

	elemType := targetType.Elem()
	result := reflect.MakeSlice(targetType, valueReflect.Len(), valueReflect.Cap())

	for i := 0; i < valueReflect.Len(); i++ {
		elem := valueReflect.Index(i).Interface()
		convertedElem, ok := ConvertToStruct[interface{}](elem, elemType)
		if !ok {
			return zero, false
		}
		result.Index(i).Set(reflect.ValueOf(convertedElem))
	}

	return result.Interface().(T), true
}

// convertMap 处理map类型转换
func convertMap[T any](value interface{}, targetType reflect.Type) (T, bool) {
	var zero T

	valueReflect := reflect.ValueOf(value)
	if valueReflect.Kind() != reflect.Map {
		return zero, false
	}

	keyType := targetType.Key()
	elemType := targetType.Elem()
	result := reflect.MakeMap(targetType)

	for _, key := range valueReflect.MapKeys() {
		// 转换key
		convertedKey, ok := ConvertToStruct[interface{}](key.Interface(), keyType)
		if !ok {
			return zero, false
		}

		// 转换value
		val := valueReflect.MapIndex(key).Interface()
		convertedVal, ok := ConvertToStruct[interface{}](val, elemType)
		if !ok {
			return zero, false
		}

		result.SetMapIndex(reflect.ValueOf(convertedKey), reflect.ValueOf(convertedVal))
	}

	return result.Interface().(T), true
}

// GetGinContextTokenLoadData 获取Gin上下文中的token的自定义数据
func GetGinContextTokenLoadData[T any](c *gin.Context) (T, error) {
	var zero T
	tokenInfo, exists := c.Get("token_info")
	if !exists {
		return zero, errors.New("token信息不存在")
	}
	claims, ok := tokenInfo.(*Claims)
	if !ok {
		return zero, errors.New("token信息类型错误")
	}
	value, exists := ConvertAnyType[T](claims.Data)
	if !exists {
		return zero, errors.New("token信息类型错误")
	}
	return value, nil
}
func GetGinContextTokenString(c *gin.Context) string {
	return c.GetString("token")
}
