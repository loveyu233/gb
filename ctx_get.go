package gb

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

// GetGinContextValue 使用泛型获取上下文中的值并转换为指定类型
func GetGinContextValue[T any](c *gin.Context, key string) T {
	var zero T

	value, exists := c.Get(key)
	if !exists {
		return zero
	}

	switch any(zero).(type) {
	case string:
		return any(cast.ToString(value)).(T)
	case int:
		return any(cast.ToInt(value)).(T)
	case int64:
		return any(cast.ToInt64(value)).(T)
	case int32:
		return any(cast.ToInt32(value)).(T)
	case int16:
		return any(cast.ToInt16(value)).(T)
	case int8:
		return any(cast.ToInt8(value)).(T)
	case uint:
		return any(cast.ToUint(value)).(T)
	case uint64:
		return any(cast.ToUint64(value)).(T)
	case uint32:
		return any(cast.ToUint32(value)).(T)
	case uint16:
		return any(cast.ToUint16(value)).(T)
	case uint8:
		return any(cast.ToUint8(value)).(T)
	case float64:
		return any(cast.ToFloat64(value)).(T)
	case float32:
		return any(cast.ToFloat32(value)).(T)
	case bool:
		return any(cast.ToBool(value)).(T)
	default:
		// 对于不支持的类型，返回零值
		return zero
	}
}
