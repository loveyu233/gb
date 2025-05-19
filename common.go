package gb

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// ScopePaginationFromGin 从gin的上下文中解析出分页参数并返回分页的scope
func ScopePaginationFromGin(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, size := ParsePaginationParams(c)

	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopePagination 分页
func ScopePagination(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopeFilterID 根据ID过滤
func ScopeFilterID(id int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// ScopeFilterStatus 根据状态过滤
func ScopeFilterStatus(status any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// ScopeDateRange 根据时间范围过滤
func ScopeDateRange(start, end *time.Time, field string) func(db *gorm.DB) *gorm.DB {
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

// ScopeFilterKeyword 根据关键字过滤
func ScopeFilterKeyword(keyword string, columns ...string) func(db *gorm.DB) *gorm.DB {
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

// SelectForUpdateTx 获取事务并加锁
func SelectForUpdateTx(tx *gorm.DB) *gorm.DB {
	return tx.Clauses(clause.Locking{Strength: "UPDATE"})
}
