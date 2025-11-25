package gb

import (
	"gorm.io/datatypes"
	"gorm.io/gen/field"
	"gorm.io/gorm/schema"
)

// GenJSONArrayQuery 用来针对 JSON 数组列构建查询表达式。
func GenJSONArrayQuery(column field.IColumnName) *datatypes.JSONArrayExpression {
	return datatypes.JSONArrayQuery(column.ColumnName().String())
}

// GenNewTime 用来为指定表字段创建时间字段对象。
func GenNewTime(table schema.Tabler, column field.IColumnName) field.Time {
	return field.NewTime(table.TableName(), column.ColumnName().String())
}

// GenNewUnsafeFieldRaw 用来创建原始 SQL 字段引用。
func GenNewUnsafeFieldRaw(rawSQL string, vars ...interface{}) field.Field {
	return field.NewUnsafeFieldRaw(rawSQL, vars...)
}
