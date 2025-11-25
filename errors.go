package gb

import (
	"errors"

	"gorm.io/gorm"
)

// IsErrRecordNotFound 用来判断错误是否表示记录不存在。
func IsErrRecordNotFound(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return true
	}
	return false
}

// redisClientNilErr 用来生成 Redis 客户端未初始化的错误。
func redisClientNilErr() error {
	return errors.New("RedisClient为空,需要先使用gb.InitRedis()进行初始化")
}

// IsErrMysqlOne 用来判断错误是否由唯一键冲突引起。
func IsErrMysqlOne(err error) bool {
	if err.Error() == "duplicated key not allowed" {
		return true
	}
	return false
}
