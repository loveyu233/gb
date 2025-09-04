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

// IsErrMysqlOne 是否是唯一键错误
func IsErrMysqlOne(err error) bool {
	if err.Error() == "duplicated key not allowed" {
		return true
	}
	return false
}
