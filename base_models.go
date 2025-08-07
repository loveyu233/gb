package gb

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
)

// BaseModel 基础模型，包含共同字段
type BaseModel struct {
	ID              int64     `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt       time.Time `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"-"`
	UpdatedAt       time.Time `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"-"`
	CreatedAtFormat string    `json:"created_at" gorm:"-"`
	UpdatedAtFormat string    `json:"updated_at" gorm:"-"`
}

// AfterFind 在查询后自动调用，格式化时间
func (m *BaseModel) AfterFind(tx *gorm.DB) error {
	m.CreatedAtFormat = DateTimeToString(m.CreatedAt)
	m.UpdatedAtFormat = DateTimeToString(m.UpdatedAt)
	return nil
}

// BaseDeleteAt 如果要使用复合索引则一定使用BaseDeleteAtContainsIndex而不是BaseDeleteAt,因为mysql中null值不能作为唯一值判定
type BaseDeleteAt struct {
	DeletedAt       gorm.DeletedAt `json:"-"`
	DeletedAtFormat string         `gorm:"-" json:"deleted_at"`
}

// AfterFind 在查询后自动调用，格式化时间
func (m *BaseDeleteAt) AfterFind(tx *gorm.DB) error {
	t := m.DeletedAt.Time
	m.DeletedAtFormat = DateTimeToString(t)
	return nil
}

// BaseDeleteAtContainsIndex 包含复合索引deleted_unique_index
type BaseDeleteAtContainsIndex struct {
	BaseDeleteAt
	DeletedAtFlag soft_delete.DeletedAt `gorm:"softDelete:flag;default:0;uniqueIndex:deleted_unique_index" json:"-"`
}
