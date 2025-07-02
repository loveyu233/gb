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
