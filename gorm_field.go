package gb

import (
	"gorm.io/datatypes"
	"gorm.io/gen/field"
	"gorm.io/gorm/schema"
)

// GenJSONArrayQuery 函数用于处理GenJSONArrayQuery相关逻辑。
func GenJSONArrayQuery(column field.IColumnName) *datatypes.JSONArrayExpression {
	return datatypes.JSONArrayQuery(column.ColumnName().String())
}

// GenNewTime 函数用于处理GenNewTime相关逻辑。
func GenNewTime(table schema.Tabler, column field.IColumnName) field.Time {
	return field.NewTime(table.TableName(), column.ColumnName().String())
}

// GenNewUnsafeFieldRaw 函数用于处理GenNewUnsafeFieldRaw相关逻辑。
func GenNewUnsafeFieldRaw(rawSQL string, vars ...interface{}) field.Field {
	return field.NewUnsafeFieldRaw(rawSQL, vars...)
}
