package gb

import (
	"errors"

	"github.com/tidwall/gjson"
)

// JsonGetValue 用来根据 key 路径从 JSON 字符串提取值。
func JsonGetValue(jsonStr string, key string) (any, error) {
	if !gjson.Valid(jsonStr) {
		return nil, errors.New("invalid json")
	}
	return gjson.Get(jsonStr, key), nil
}
