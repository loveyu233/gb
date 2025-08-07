package gb

import (
	"errors"

	"gorm.io/gorm"
)

// IsErrRecordNotFound 判断err是否为gorm.ErrRecordNotFound
func IsErrRecordNotFound(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}
func redisClientNilErr() error {
	return errors.New("RedisClient为空,需要先使用gb.InitRedis()进行初始化")
}
