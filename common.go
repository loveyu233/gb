package gb

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
)

// ScopeOrderDesc 用来构造按指定列倒序排序的 GORM Scope。
func (db *GormClient) ScopeOrderDesc(columnName ...string) func(db *gorm.DB) *gorm.DB {
	orderColumn := safeColumnNameFromArgs(columnName, "created_at")
	order := fmt.Sprintf("%s desc", orderColumn)
	return func(db *gorm.DB) *gorm.DB {
		db.Order(order)
		return db
	}
}

// SelectByID 用来根据主键查询记录并填充目标模型。
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

// ScopePaginationFromGin 用来根据 gin 上下文中的分页参数设置偏移和限制。
func (db *GormClient) ScopePaginationFromGin(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, size := ParsePaginationParams(c)

	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopePagination 用来依据给定页码和条数限制查询结果。
func (db *GormClient) ScopePagination(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopeFilterID 用来在查询中添加按 id 精确匹配的条件。
func (db *GormClient) ScopeFilterID(id int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// ScopeFilterStatus 用来为查询追加按状态字段过滤的条件。
func (db *GormClient) ScopeFilterStatus(status any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// ScopeDateRange 用来对时间字段追加开始结束区间的过滤。
func (db *GormClient) ScopeDateRange(field string, start, end *time.Time) func(db *gorm.DB) *gorm.DB {
	column := safeColumnName(field, "created_at")
	return func(db *gorm.DB) *gorm.DB {
		if column == "" {
			return db
		}
		if start == nil && end == nil {
			return db
		}

		if start != nil && end == nil {
			return db.Where(column+" >= ?", start)
		}

		if start == nil && end != nil {
			return db.Where(column+" <= ?", end)
		}

		return db.Where(column+" BETWEEN ? AND ?", start, end)
	}
}

// ScopeFilterKeyword 用来在多列上构造关键字模糊匹配条件。
func (db *GormClient) ScopeFilterKeyword(keyword string, columns ...string) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if keyword == "" || len(columns) == 0 {
			return db
		}

		var queryParts []string
		var args []interface{}
		for _, column := range columns {
			safeCol := safeColumnName(column, "")
			if safeCol == "" {
				continue
			}
			queryParts = append(queryParts, safeCol+" LIKE ?")
			args = append(args, "%"+keyword+"%")
		}

		if len(queryParts) == 0 {
			return db
		}

		query := strings.Join(queryParts, " OR ")
		return db.Where(query, args...)
	}
}

// SelectForUpdateTx 用来返回带 FOR UPDATE 锁的事务查询。
func (db *GormClient) SelectForUpdateTx() *gorm.DB {
	return db.DB.Clauses(clause.Locking{Strength: "UPDATE"})
}

// ScopeTime 用来生成指定列在给定范围内的过滤条件。
func (db *GormClient) ScopeTime(start, end string, columns ...string) func(db *gorm.DB) *gorm.DB {
	column := safeColumnNameFromArgs(columns, "created_at")
	return func(db *gorm.DB) *gorm.DB {
		if column == "" {
			return db
		}
		return db.Where(fmt.Sprintf("%s >= ? and %s < ?", column, column), start, end)
	}
}

// Transaction 用来包装 GORM 事务并执行回调函数。
func (db *GormClient) Transaction(tx func(tx *gorm.DB) error) error {
	return db.DB.Transaction(tx)
}

// Lock 用来为当前查询添加 UPDATE 锁。
func (db *GormClient) Lock() *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}

var columnNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_\.]+$`)

func safeColumnNameFromArgs(columns []string, fallback string) string {
	if len(columns) == 0 {
		return safeColumnName("", fallback)
	}
	return safeColumnName(columns[0], fallback)
}

func safeColumnName(column string, fallback string) string {
	column = strings.TrimSpace(column)
	if column == "" {
		column = fallback
	}
	if column == "" {
		return ""
	}
	if columnNamePattern.MatchString(column) {
		return column
	}
	if fallback != "" && columnNamePattern.MatchString(fallback) {
		return fallback
	}
	return ""
}
