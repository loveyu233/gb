package gb

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"
	"time"
)

// ScopeOrderDesc 倒序
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

// SelectByID 查询指定id数据,obj必须为指针
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

// ScopePaginationFromGin 从gin的上下文中解析出分页参数并返回分页的scope
func (db *GormClient) ScopePaginationFromGin(c *gin.Context) func(db *gorm.DB) *gorm.DB {
	page, size := ParsePaginationParams(c)

	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopePagination 分页,默认从1开始
func (db *GormClient) ScopePagination(page, size int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if size == -1 {
			return db
		}

		return db.Offset((page - 1) * size).Limit(size)
	}
}

// ScopeFilterID 根据ID过滤
func (db *GormClient) ScopeFilterID(id int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id = ?", id)
	}
}

// ScopeFilterStatus 根据状态过滤
func (db *GormClient) ScopeFilterStatus(status any) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("status = ?", status)
	}
}

// ScopeDateRange 根据时间范围过滤
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

// ScopeFilterKeyword 根据关键字过滤
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

// SelectForUpdateTx 获取事务并加锁
func (db *GormClient) SelectForUpdateTx() *gorm.DB {
	return db.DB.Clauses(clause.Locking{Strength: "UPDATE"})
}

// ScopeTime 查询时间范围, start <= column < end
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

// ScopeToday 今天
func (db *GormClient) ScopeToday(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetTodayInterval()
	return db.ScopeTime(start, end, columns...)
}

// ScopeYesterday 昨天
func (db *GormClient) ScopeYesterday(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetYesterdayInterval()
	return db.ScopeTime(start, end, columns...)
}

// ScopeLastMonth 上个月
func (db *GormClient) ScopeLastMonth(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetLastMonthInterval()
	return db.ScopeTime(start, end, columns...)
}

// ScopeCurrentMonth 本月
func (db *GormClient) ScopeCurrentMonth(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetCurrentMonthInterval()
	return db.ScopeTime(start, end, columns...)
}

// ScopeNextMonth 下个月
func (db *GormClient) ScopeNextMonth(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetNextMonthInterval()
	return db.ScopeTime(start, end, columns...)
}

// ScopeLastYears 去年
func (db *GormClient) ScopeLastYears(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetLastYearsInterval()
	return db.ScopeTime(start, end, columns...)
}

// ScopeCurrentYears 今年
func (db *GormClient) ScopeCurrentYears(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetCurrentYearsInterval()
	return db.ScopeTime(start, end, columns...)
}

// ScopeNextYears 明年
func (db *GormClient) ScopeNextYears(columns ...string) func(db *gorm.DB) *gorm.DB {
	start, end := GetNextYearsInterval()
	return db.ScopeTime(start, end, columns...)
}

// Transaction 开启事务(db.Transaction)
func (db *GormClient) Transaction(tx func(tx *gorm.DB) error) error {
	return db.DB.Transaction(tx)
}

// Lock 加锁
func (db *GormClient) Lock() *gorm.DB {
	return db.Clauses(clause.Locking{Strength: "UPDATE"})
}
