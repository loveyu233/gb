package gb

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ScopeOrderDesc 方法用于处理ScopeOrderDesc相关逻辑。
func (db *GormClient) ScopeOrderDesc(columnName ...string) func(db *gorm.DB) *gorm.DB {
	var order string
	if len(columnName) > 0 {
		order = fmt.Sprintf("%s desc", columnName[0])
	} else {
		order = fmt.Sprintf("created_at desc")
	}
	return func(db *gorm.DB) *gorm.DB {
		db.Order(order)
		return db
	}
}

// SelectByID 方法用于处理SelectByID相关逻辑。
func (db *GormClient) SelectByID(obj schema.Tabler, id any) error {
	if !IsPtr(obj) {
		return errors.New("obj必须是指针类型")
	}

	// 使用GORM根据ID查询数据
	result := db.DB.First(obj, id)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

// ScopePaginationFromGin 方法用于处理ScopePaginationFromGin相关逻辑。
func (db *GormClient) ScopePaginationFromGin(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, size := ParsePaginationParams(c)

	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopePagination 方法用于处理ScopePagination相关逻辑。
func (db *GormClient) ScopePagination(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopeFilterID 方法用于处理ScopeFilterID相关逻辑。
func (db *GormClient) ScopeFilterID(id int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// ScopeFilterStatus 方法用于处理ScopeFilterStatus相关逻辑。
func (db *GormClient) ScopeFilterStatus(status any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// ScopeDateRange 方法用于处理ScopeDateRange相关逻辑。
func (db *GormClient) ScopeDateRange(field string, start, end *time.Time) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if start == nil && end == nil {
			return db
		}

		if start != nil && end == nil {
			return db.Where(field+" >= ?", start)
		}

		if start == nil && end != nil {
			return db.Where(field+" <= ?", end)
		}

		return db.Where(field+" BETWEEN ? AND ?", start, end)
	}
}

// ScopeFilterKeyword 方法用于处理ScopeFilterKeyword相关逻辑。
func (db *GormClient) ScopeFilterKeyword(keyword string, columns ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(columns) == 0 {
			return db
		}

		// 使用括号将所有 OR 条件包裹起来
		var query string
		args := make([]interface{}, len(columns))
		for i, column := range columns {
			if i == 0 {
				query += column + " LIKE ?"
			} else {
				query += " OR " + column + " LIKE ?"
			}
			args[i] = "%" + keyword + "%"
		}

		// 使用 db.Where 传递完整的查询字符串和参数
		return db.Where(query, args...)
	}
}

// SelectForUpdateTx 方法用于处理SelectForUpdateTx相关逻辑。
func (db *GormClient) SelectForUpdateTx() *gorm.DB {
	return db.DB.Clauses(clause.Locking{Strength: "UPDATE"})
}

// ScopeTime 方法用于处理ScopeTime相关逻辑。
func (db *GormClient) ScopeTime(start, end string, columns ...string) func(db *gorm.DB) *gorm.DB {
	var column string
	if len(columns) > 0 {
		column = columns[0]
	} else {
		column = "created_at"
	}
	return func(db *gorm.DB) *gorm.DB {
		return db.Where(fmt.Sprintf("%s >= '%s' and %s < '%s'", column, start, column, end))
	}
}

// Transaction 方法用于处理Transaction相关逻辑。
func (db *GormClient) Transaction(tx func(tx *gorm.DB) error) error {
	return db.DB.Transaction(tx)
}

// Lock 方法用于处理Lock相关逻辑。
func (db *GormClient) Lock() *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}
