package gb

import (
	"gorm.io/gorm"
	"gorm.io/plugin/soft_delete"
	"time"
)

// BaseModel 基础模型，包含共同字段
type BaseModel struct {
	ID              int64                 `gorm:"column:id;type:bigint(20);primary_key;AUTO_INCREMENT" json:"id"`
	CreatedAt       time.Time             `gorm:"column:created_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"-"`
	UpdatedAt       time.Time             `gorm:"column:updated_at;type:datetime;default:CURRENT_TIMESTAMP;NOT NULL" json:"-"`
	DeletedAt       soft_delete.DeletedAt `gorm:"column:deleted_at;type:datetime;default:0" json:"-"`
	CreatedAtFormat string                `json:"created_at" gorm:"-"`
	UpdatedAtFormat string                `json:"updated_at" gorm:"-"`
}

// AfterFind 在查询后自动调用，格式化时间
func (m *BaseModel) AfterFind(tx *gorm.DB) error {
	m.CreatedAtFormat = DateTimeToString(m.CreatedAt)
	m.UpdatedAtFormat = DateTimeToString(m.UpdatedAt)
	return nil
}
