package gb

import (
	"gorm.io/datatypes"
	"gorm.io/gen/field"
)

// GenJSONArrayQuery datatypes.JSONArrayQuery(columnName)
func GenJSONArrayQuery(columnName string) *datatypes.JSONArrayExpression {
	return datatypes.JSONArrayQuery(columnName)
}

// GenNewTime field.NewTime(tableName, columnName)
func GenNewTime(tableName string, columnName string) field.Time {
	return field.NewTime(tableName, columnName)
}

// GenNewUnsafeFieldRaw field.NewUnsafeFieldRaw(rawSQL, vars...)
func GenNewUnsafeFieldRaw(rawSQL string, vars ...interface{}) field.Field {
	return field.NewUnsafeFieldRaw(rawSQL, vars...)
}
