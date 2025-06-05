package gb

import (
	"crypto/rand"
	"github.com/google/uuid"
	"github.com/loveyu233/gb/snowflake"
	"github.com/rs/xid"
	"math/big"
)

// GetUUID 长度为36的字符串
func GetUUID() string {
	return uuid.NewString()
}

// GetXID 长度为20的字符串
func GetXID() string {
	return xid.New().String()
}

// GetSnowflakeID 长度为18的数字
func GetSnowflakeID() int64 {
	return snowflake.GetId()
}

// RandomString 获取指定长度的随机字符串
func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			panic(err)
		}
		result[i] = charset[num.Int64()]
	}
	return string(result)
}

// RandomStringWithPrefix 带指定前后缀的指定长度字符串
func RandomStringWithPrefix(length int, prefix, suffix string) string {
	if length <= len(prefix)+len(suffix) {
		return prefix + suffix
	}

	randomLength := length - len(prefix) - len(suffix)
	randomPart := RandomString(randomLength)

	return prefix + randomPart + suffix
}
