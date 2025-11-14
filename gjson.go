package gb

import (
	"errors"

	"github.com/tidwall/gjson"
)

// JsonGetValue 函数用于处理JsonGetValue相关逻辑。
func JsonGetValue(jsonStr string, key string) (any, error) {
	if !gjson.Valid(jsonStr) {
		return nil, errors.New("invalid json")
	}
	return gjson.Get(jsonStr, key), nil
}
