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
	ColumnType       string            // 字段类型,时间,日期默认使用gb实现,其他类型写对应go包路径,例如:model.User
	IsJsonStatusType bool              // 默认false设置为true自动添加标签:gorm:column:ColumnName;serializer:json,如果为true且在Tags中设置了gorm则会忽略,需要自行添加serializer:json
	Tags             map[string]string // 可以设置生成后字段的标签,key为标签名,value为标签值
}

type GenConfig struct {
	outFilePath            string
	globalColumnType       map[string]func(gorm.ColumnType) string
	globalSimpleColumnType []GenFieldType
	useTablesName          []string
	tableColumnType        map[string][]GenFieldType
}

type WithGenConfig func(*GenConfig)

func WithGenOutFilePath(outFilePath string) WithGenConfig {
	return func(gc *GenConfig) {
		gc.outFilePath = outFilePath
	}
}

// WithGenTableColumnType 只有使用了WithGenTableColumnType该方法才有效,map的key为表名称,value为自定义数据类型,方式和全局定义一样
func WithGenTableColumnType(value map[string][]GenFieldType) WithGenConfig {
	return func(gc *GenConfig) {
		gc.tableColumnType = value
	}
}

func WithGenUseTablesName(tablesName ...string) WithGenConfig {
	return func(gc *GenConfig) {
		gc.useTablesName = tablesName
	}
}

func WithGenGlobalSimpleColumnType(fields []GenFieldType) WithGenConfig {
	return func(gc *GenConfig) {
		gc.globalSimpleColumnType = append(gc.globalSimpleColumnType, fields...)
	}
}

// WithGenGlobalSimpleColumnTypeAddJsonSliceType 简化的使用方式,例如:WithGenGlobalSimpleColumnTypeAddJsonSliceType("arrFieldName","int64")
func WithGenGlobalSimpleColumnTypeAddJsonSliceType(sliceFieldName, sliceType string) WithGenConfig {
	return func(gc *GenConfig) {
		gc.globalSimpleColumnType = append(gc.globalSimpleColumnType, GenFieldType{
			ColumnName:       sliceFieldName,
			ColumnType:       fmt.Sprintf("datatypes.JSONSlice[%s]", sliceType),
			IsJsonStatusType: true,
		})
	}
}

func WithGenGlobalSimpleColumnTypeAddJsonType(sliceFieldName, sliceType string) WithGenConfig {
	return func(gc *GenConfig) {
		gc.globalSimpleColumnType = append(gc.globalSimpleColumnType, GenFieldType{
			ColumnName:       sliceFieldName,
			ColumnType:       sliceType,
			IsJsonStatusType: true,
		})
	}
}

func WithGenGlobalColumnType(value map[string]func(gorm.ColumnType) string) WithGenConfig {
	return func(gc *GenConfig) {
		if len(gc.globalColumnType) == 0 {
			gc.globalColumnType = value
		} else {
			for k, v := range value {
				gc.globalColumnType[k] = v
			}
		}
	}
}

// WithGenGlobalColumnTypeAddDatatypes 使用官方的date和time,官方包的date输出的日期包含了时间部分,自行取舍
func WithGenGlobalColumnTypeAddDatatypes() WithGenConfig {
	return func(gc *GenConfig) {
		if len(gc.globalColumnType) == 0 {
			gc.globalColumnType = map[string]func(gorm.ColumnType) string{
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
			}
		} else {
			gc.globalColumnType["date"] = func(columnType gorm.ColumnType) (dataType string) {
				if nullable, ok := columnType.Nullable(); ok && nullable {
					return "*datatypes.Date"
				}
				return "datatypes.Date"
			}
			gc.globalColumnType["time"] = func(columnType gorm.ColumnType) (dataType string) {
				if nullable, ok := columnType.Nullable(); ok && nullable {
					return "*datatypes.Time"
				}
				return "datatypes.Time"
			}
		}

	}
}

// Gen 默认字段如果是date和time类型生成go的类型为gb实现的,如果要使用官方库的时间和日期类型,请使用WithGenGlobalColumnTypeAddDatatypes,输出路径不设置默认为:gen/query
func (db *GormClient) Gen(opts ...WithGenConfig) {
	var genConfig = new(GenConfig)
	for i := range opts {
		opts[i](genConfig)
	}

	if genConfig.outFilePath == "" {
		genConfig.outFilePath = "gen/query"
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:        fmt.Sprintf(genConfig.outFilePath),
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
				return "*gb.DateTime"
			}
			return "gb.DateTime"
		},

		"timestamp": func(columnType gorm.ColumnType) (dataType string) {
			if nullable, ok := columnType.Nullable(); ok && nullable {
				return "*gb.DateTime"
			}
			return "gb.DateTime"
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

	for k, v := range genConfig.globalColumnType {
		dataMap[k] = v
	}

	g.WithDataTypeMap(dataMap)
	g.UseDB(db.DB)

	var fieldTypes []gen.ModelOpt
	for _, item := range genConfig.globalSimpleColumnType {
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

	if len(genConfig.useTablesName) > 0 {
		var gms []interface{}
		for _, table := range genConfig.useTablesName {
			var opts []gen.ModelOpt
			if genFieldTypes, ok := genConfig.tableColumnType[table]; ok {
				for _, fieldType := range genFieldTypes {
					if fieldType.ColumnName == "" {
						panic("column_name不能为空")
					}
					if fieldType.ColumnType != "" {
						opts = append(opts, gen.FieldType(fieldType.ColumnName, fieldType.ColumnType))
					}
					if fieldType.IsJsonStatusType {
						if len(fieldType.Tags) == 0 {
							fieldType.Tags = map[string]string{
								"gorm": fmt.Sprintf("column:%s;serializer:json", fieldType.ColumnName),
							}
						} else {
							if _, ok := fieldType.Tags["gorm"]; !ok {
								fieldType.Tags["gorm"] = fmt.Sprintf("column:%s;serializer:json", fieldType.ColumnName)
							}
						}
					}
					if len(fieldType.Tags) > 0 {
						opts = append(opts, gen.FieldTag(fieldType.ColumnName, func(tag field.Tag) field.Tag {
							for k, v := range fieldType.Tags {
								tag.Set(k, v)
							}
							return tag
						}))
					}
				}
			}
			gms = append(gms, g.GenerateModel(table, opts...))
		}
		g.ApplyBasic(gms...)
	} else {
		g.ApplyBasic(g.GenerateAllTable(fieldTypes...)...)
	}
	g.Execute()
}
