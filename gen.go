package gb

import (
	"fmt"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"strings"
)

type GenFieldType struct {
	ColumnName string
	ColumnType string
}

type GenTag struct {
	ColumnName string
	Tags       map[string]string
}

type GenConfig struct {
	OutFilePath string
	fieldType   []GenFieldType
	tags        []GenTag
}

type WithGenConfig func(*GenConfig)

func WithGenOutFilePath(outFilePath string) WithGenConfig {
	return func(gc *GenConfig) {
		gc.OutFilePath = outFilePath
	}
}

func WithGenFieldType(fields []GenFieldType) WithGenConfig {
	return func(gc *GenConfig) {
		gc.fieldType = fields
	}
}

func WithGenJsonTag(tags []GenTag) WithGenConfig {
	return func(gc *GenConfig) {
		gc.tags = tags
	}
}

func (db *GormClient) Gen(opts ...WithGenConfig) {
	var genConfig = new(GenConfig)
	for i := range opts {
		opts[i](genConfig)
	}

	if genConfig.OutFilePath == "" {
		genConfig.OutFilePath = "gen/query"
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:        fmt.Sprintf(genConfig.OutFilePath),
		FieldCoverable: false,
		Mode:           gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
	})

	var dataMap = map[string]func(columnType gorm.ColumnType) (dataType string){
		// 整数类型
		"tinyint": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
					return "*uint8"
				}
				return "*int8"
			}
			if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
				return "uint8"
			}
			return "int8"
		},

		"smallint": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
					return "*uint16"
				}
				return "*int16"
			}
			if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
				return "uint16"
			}
			return "int16"
		},

		"mediumint": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
					return "*uint32"
				}
				return "*int32"
			}
			if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
				return "uint32"
			}
			return "int32"
		},

		"int": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
					return "*uint32"
				}
				return "*int32"
			}
			if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
				return "uint32"
			}
			return "int32"
		},

		"integer": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
					return "*uint32"
				}
				return "*int32"
			}
			if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
				return "uint32"
			}
			return "int32"
		},

		"bigint": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
					return "*uint64"
				}
				return "*int64"
			}
			if unsigned, ok := columnType.ColumnType(); ok && strings.Contains(strings.ToLower(unsigned), "unsigned") {
				return "uint64"
			}
			return "int64"
		},

		"bit": func(columnType gorm.ColumnType) (dataType string) {
			if length, ok := columnType.Length(); ok && length == 1 {
				if nullable, ok := columnType.Nullable(); ok && nullable {
					return "*bool"
				}
				return "bool"
			}
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		// 浮点类型
		"float": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*float32"
			}
			return "float32"
		},

		"double": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*float64"
			}
			return "float64"
		},

		"real": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*float64"
			}
			return "float64"
		},

		"decimal": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*decimal.Decimal"
			}
			return "decimal.Decimal"
		},

		"numeric": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		// 字符串类型
		"char": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		"varchar": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		"tinytext": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		"text": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		"mediumtext": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		"longtext": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		// 二进制类型
		"binary": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"varbinary": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"tinyblob": func(columnType gorm.ColumnType) (dataType string) {
			return "[]byte"
		},

		"blob": func(columnType gorm.ColumnType) (dataType string) {
			return "[]byte"
		},

		"mediumblob": func(columnType gorm.ColumnType) (dataType string) {
			return "[]byte"
		},

		"longblob": func(columnType gorm.ColumnType) (dataType string) {
			return "[]byte"
		},

		// 日期时间类型
		"date": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*gb.DateOnly"
			}
			return "gb.DateOnly"
		},

		"time": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*gb.TimeOnly"
			}
			return "gb.TimeOnly"
		},

		"datetime": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*time.Time"
			}
			return "time.Time"
		},

		"timestamp": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*time.Time"
			}
			return "time.Time"
		},

		"year": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*int"
			}
			return "int"
		},

		// JSON 类型 (MySQL 5.7+)
		"json": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*datatypes.JSON"
			}
			return "datatypes.JSON"
		},

		// 枚举和集合类型
		"enum": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		"set": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*string"
			}
			return "string"
		},

		// 空间数据类型
		"geometry": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"point": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"linestring": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"polygon": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"multipoint": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"multilinestring": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"multipolygon": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		"geometrycollection": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*[]byte"
			}
			return "[]byte"
		},

		// 布尔类型 (通常用 TINYINT(1) 表示)
		"boolean": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*bool"
			}
			return "bool"
		},

		"bool": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*bool"
			}
			return "bool"
		},
	}

	g.WithDataTypeMap(dataMap)
	g.UseDB(db.DB)

	var fieldTypes []gen.ModelOpt
	for _, item := range genConfig.fieldType {
		fieldTypes = append(fieldTypes, gen.FieldType(item.ColumnName, item.ColumnType))
	}

	for _, item := range genConfig.tags {
		gen.FieldTag(item.ColumnName, func(tag field.Tag) field.Tag {
			for k, v := range item.Tags {
				tag.Set(k, v)
			}
			return tag
		})
	}

	g.ApplyBasic(g.GenerateAllTable(fieldTypes...)...)
	g.Execute()
}
