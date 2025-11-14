package gb

import (
	"errors"

	"gorm.io/gorm"
)

// IsErrRecordNotFound 函数用于处理IsErrRecordNotFound相关逻辑。
func IsErrRecordNotFound(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}

// redisClientNilErr 函数用于处理redisClientNilErr相关逻辑。
func redisClientNilErr() error {
	return errors.New("RedisClient为空,需要先使用gb.InitRedis()进行初始化")
}

// IsErrMysqlOne 函数用于处理IsErrMysqlOne相关逻辑。
func IsErrMysqlOne(err error) bool {
	if err.Error() == "duplicated key not allowed" {
		return true
	}
	return false
}
