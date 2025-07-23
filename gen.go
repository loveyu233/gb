package gb

import (
	"fmt"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"strings"
)

type GenFieldType struct {
	ColumnName       string            // 字段名称,表中字段的名称不是结构体的名称
	ColumnType       string            // 字段类型,时间,日期默认为datatypes.Time和datatypes.Date,如果是json类型的数组自行设置为datatypes.JSONSlice[类型]无其他需求不需要设置IsJsonStatusType和Tags,如果是结构体则写入结构体路径,例如:model.User,无其他需求设置IsJsonStatusType为true即可
	IsJsonStatusType bool              // 默认false设置为true自动添加标签:gorm:column:ColumnName;serializer:json,如果为true且在Tags中设置了gorm则会忽略,需要自行添加serializer:json
	Tags             map[string]string // 可以设置生成后字段的标签,key为标签名,value为标签值
}

type GenConfig struct {
	OutFilePath string
	fieldType   []GenFieldType
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
				return "*datatypes.Date"
			}
			return "datatypes.Date"
		},

		"time": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*datatypes.Time"
			}
			return "datatypes.Time"
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
		if item.ColumnName == "" {
			panic("column_name不能为空")
		}
		if item.ColumnType != "" {
			fieldTypes = append(fieldTypes, gen.FieldType(item.ColumnName, item.ColumnType))
		}
		if item.IsJsonStatusType {
			if len(item.Tags) == 0 {
				item.Tags = map[string]string{
					"gorm": fmt.Sprintf("column:%s;serializer:json", item.ColumnName),
				}
			} else {
				if _, ok := item.Tags["gorm"]; !ok {
					item.Tags["gorm"] = fmt.Sprintf("column:%s;serializer:json", item.ColumnName)
				}
			}
		}
		if len(item.Tags) > 0 {
			fieldTypes = append(fieldTypes, gen.FieldTag(item.ColumnName, func(tag field.Tag) field.Tag {
				for k, v := range item.Tags {
					tag.Set(k, v)
				}
				return tag
			}))
		}
	}

	g.ApplyBasic(g.GenerateAllTable(fieldTypes...)...)
	g.Execute()
}
