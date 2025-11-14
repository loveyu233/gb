package gb

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/xuri/excelize/v2"
)

// ExcelMapper Excel数据映射器
type ExcelMapper struct {
	// 配置选项
	SheetName    string // 工作表名称
	SheetIndex   int    // 工作表索引 (0-based)
	HeaderRow    int    // 表头行号 (1-based)
	DataStartRow int    // 数据开始行号 (1-based)
	StrictMode   bool   // 严格模式，遇到错误即停止

	// 内部缓存
	fieldCache map[reflect.Type]*structInfo
	cacheMutex sync.RWMutex

	// 错误收集
	errors []MappingError
}

// MappingError 映射错误信息
type MappingError struct {
	Row     int    // 行号
	Column  string // 列名
	Field   string // 字段名
	Value   string // 原始值
	Message string // 错误信息
}

// Error 方法用于处理Error相关逻辑。
func (e MappingError) Error() string {
	return fmt.Sprintf("row %d, column %s, field %s: %s (value: %s)",
		e.Row, e.Column, e.Field, e.Message, e.Value)
}

// structInfo 结构体信息缓存
type structInfo struct {
	fields []fieldInfo
}

// fieldInfo 字段映射信息
type fieldInfo struct {
	index     int            // 字段索引
	name      string         // 字段名
	tag       string         // excel标签值
	fieldType reflect.Type   // 字段类型
	isPointer bool           // 是否为指针类型
	converter valueConverter // 值转换器
}

// valueConverter 值转换器接口
type valueConverter interface {
	Convert(value string) (interface{}, error)
}

// 预定义转换器
var (
	stringConverter  = &stringConv{}
	intConverter     = &intConv{}
	int64Converter   = &int64Conv{}
	float64Converter = &float64Conv{}
	boolConverter    = &boolConv{}
	timeConverter    = &timeConv{}
)

type WithExcelMapperOption func(m *ExcelMapper)

// WithExcelMapperSheetName 函数用于处理WithExcelMapperSheetName相关逻辑。
func WithExcelMapperSheetName(sheetName string) WithExcelMapperOption {
	return func(m *ExcelMapper) {
		m.SheetName = sheetName
	}
}

// WithExcelMapperSheetIndex 函数用于处理WithExcelMapperSheetIndex相关逻辑。
func WithExcelMapperSheetIndex(sheetIndex int) WithExcelMapperOption {
	return func(m *ExcelMapper) {
		m.SheetIndex = sheetIndex
	}
}

// WithExcelMapperHeaderRow 函数用于处理WithExcelMapperHeaderRow相关逻辑。
func WithExcelMapperHeaderRow(headerRow int) WithExcelMapperOption {
	return func(m *ExcelMapper) {
		m.HeaderRow = headerRow
	}
}

// WithExcelMapperDataStartRow 函数用于处理WithExcelMapperDataStartRow相关逻辑。
func WithExcelMapperDataStartRow(dataStartRow int) WithExcelMapperOption {
	return func(m *ExcelMapper) {
		m.DataStartRow = dataStartRow
	}
}

// WithExcelMapperStrictMode 函数用于处理WithExcelMapperStrictMode相关逻辑。
func WithExcelMapperStrictMode(strictMode bool) WithExcelMapperOption {
	return func(m *ExcelMapper) {
		m.StrictMode = strictMode
	}
}

// InitExcelMapper 函数用于处理InitExcelMapper相关逻辑。
func InitExcelMapper(options ...WithExcelMapperOption) *ExcelMapper {
	excelMapper := new(ExcelMapper)
	for i := range options {
		options[i](excelMapper)
	}
	if excelMapper.HeaderRow == 0 {
		excelMapper.HeaderRow = 1
	}
	if excelMapper.DataStartRow == 0 {
		excelMapper.DataStartRow = 2
	}
	excelMapper.fieldCache = make(map[reflect.Type]*structInfo)
	excelMapper.errors = make([]MappingError, 0)
	return excelMapper
}

// MapToStructs 方法用于处理MapToStructs相关逻辑。
func (m *ExcelMapper) MapToStructs(filePath string, result interface{}) error {
	// 参数验证
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr || resultValue.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("结果必须是指向切片的指针")
	}

	sliceValue := resultValue.Elem()
	elemType := sliceValue.Type().Elem()

	// 获取结构体信息
	structInfo, err := m.getStructInfo(elemType)
	if err != nil {
		return err
	}

	// 打开Excel文件
	file, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("打开excel文件失败: %v", err)
	}
	defer file.Close()

	// 确定工作表名称
	sheetName := m.SheetName
	if sheetName == "" {
		sheets := file.GetSheetList()
		if len(sheets) == 0 {
			return fmt.Errorf("在excel文件中找不到工作表")
		}
		if m.SheetIndex >= len(sheets) {
			return fmt.Errorf("工作表索引 %d 超出范围", m.SheetIndex)
		}
		sheetName = sheets[m.SheetIndex]
	}

	// 获取数据范围
	rows, err := file.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("无法获取行: %v", err)
	}

	if len(rows) < m.DataStartRow {
		return fmt.Errorf("excel文件中的行不足")
	}

	// 构建列映射
	columnMap, err := m.buildColumnMap(rows, structInfo)
	if err != nil {
		return err
	}

	// 预分配切片容量
	dataRows := rows[m.DataStartRow-1:]
	sliceValue.Set(reflect.MakeSlice(sliceValue.Type(), 0, len(dataRows)))

	// 批量处理数据
	return m.processDataRows(dataRows, sliceValue, elemType, structInfo, columnMap)
}

// getStructInfo 方法用于处理getStructInfo相关逻辑。
func (m *ExcelMapper) getStructInfo(elemType reflect.Type) (*structInfo, error) {
	m.cacheMutex.RLock()
	if info, exists := m.fieldCache[elemType]; exists {
		m.cacheMutex.RUnlock()
		return info, nil
	}
	m.cacheMutex.RUnlock()

	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	// 双重检查
	if info, exists := m.fieldCache[elemType]; exists {
		return info, nil
	}

	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("目标类型必须为struct")
	}

	info := &structInfo{
		fields: make([]fieldInfo, 0, elemType.NumField()),
	}

	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("excel")
		if tag == "" || tag == "-" {
			continue
		}

		fieldType := field.Type
		isPointer := fieldType.Kind() == reflect.Ptr
		if isPointer {
			fieldType = fieldType.Elem()
		}

		converter := m.getConverter(fieldType)
		if converter == nil {
			return nil, fmt.Errorf("不支持的字段类型: %v", fieldType)
		}

		info.fields = append(info.fields, fieldInfo{
			index:     i,
			name:      field.Name,
			tag:       tag,
			fieldType: fieldType,
			isPointer: isPointer,
			converter: converter,
		})
	}

	m.fieldCache[elemType] = info
	return info, nil
}

// getConverter 方法用于处理getConverter相关逻辑。
func (m *ExcelMapper) getConverter(fieldType reflect.Type) valueConverter {
	switch fieldType.Kind() {
	case reflect.String:
		return stringConverter
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return intConverter
	case reflect.Int64:
		return int64Converter
	case reflect.Float32, reflect.Float64:
		return float64Converter
	case reflect.Bool:
		return boolConverter
	}

	// 检查是否为time.Time
	if fieldType == reflect.TypeOf(time.Time{}) {
		return timeConverter
	}

	return nil
}

// buildColumnMap 方法用于处理buildColumnMap相关逻辑。
func (m *ExcelMapper) buildColumnMap(rows [][]string, structInfo *structInfo) (map[string]int, error) {
	if len(rows) < m.HeaderRow {
		return nil, fmt.Errorf("标题行 %d 未找到", m.HeaderRow)
	}

	headers := rows[m.HeaderRow-1]
	columnMap := make(map[string]int, len(structInfo.fields))

	for _, fieldInfo := range structInfo.fields {
		columnIndex := m.findColumnIndex(headers, fieldInfo.tag)
		if columnIndex == -1 {
			if m.StrictMode {
				return nil, fmt.Errorf("未找到字段的列 %s 带标签 %s",
					fieldInfo.name, fieldInfo.tag)
			}
			continue
		}
		columnMap[fieldInfo.tag] = columnIndex
	}

	return columnMap, nil
}

// findColumnIndex 方法用于处理findColumnIndex相关逻辑。
func (m *ExcelMapper) findColumnIndex(headers []string, tag string) int {
	// 优先精确匹配
	for i, header := range headers {
		if strings.TrimSpace(header) == tag {
			return i
		}
	}

	// 忽略大小写匹配
	tagLower := strings.ToLower(strings.TrimSpace(tag))
	for i, header := range headers {
		if strings.ToLower(strings.TrimSpace(header)) == tagLower {
			return i
		}
	}

	// 检查是否为列索引或列名
	if colIndex, err := strconv.Atoi(tag); err == nil && colIndex >= 0 && colIndex < len(headers) {
		return colIndex
	}

	// 检查是否为Excel列名 (A, B, AA等)
	if colIndex := m.parseExcelColumn(tag); colIndex >= 0 && colIndex < len(headers) {
		return colIndex
	}

	return -1
}

// parseExcelColumn 方法用于处理parseExcelColumn相关逻辑。
func (m *ExcelMapper) parseExcelColumn(col string) int {
	col = strings.ToUpper(strings.TrimSpace(col))
	if col == "" {
		return -1
	}

	result := 0
	for _, char := range col {
		if char < 'A' || char > 'Z' {
			return -1
		}
		result = result*26 + int(char-'A'+1)
	}
	return result - 1
}

// processDataRows 方法用于处理processDataRows相关逻辑。
func (m *ExcelMapper) processDataRows(dataRows [][]string, sliceValue reflect.Value,
	elemType reflect.Type, structInfo *structInfo, columnMap map[string]int) error {

	// 批量创建结构体实例
	instances := make([]reflect.Value, 0, len(dataRows))

	for rowIndex, row := range dataRows {
		actualRow := m.DataStartRow + rowIndex

		// 跳过空行
		if m.isEmptyRow(row) {
			continue
		}

		instance := reflect.New(elemType).Elem()
		hasError := false

		// 设置字段值
		for _, fieldInfo := range structInfo.fields {
			columnIndex, exists := columnMap[fieldInfo.tag]
			if !exists {
				continue
			}

			var cellValue string
			if columnIndex < len(row) {
				cellValue = strings.TrimSpace(row[columnIndex])
			}

			if err := m.setFieldValue(instance, fieldInfo, cellValue, actualRow, fieldInfo.tag); err != nil {
				if m.StrictMode {
					return err
				}
				hasError = true
			}
		}

		if !hasError || !m.StrictMode {
			instances = append(instances, instance)
		}
	}

	// 批量添加到切片
	for _, instance := range instances {
		sliceValue.Set(reflect.Append(sliceValue, instance))
	}

	return nil
}

// setFieldValue 方法用于处理setFieldValue相关逻辑。
func (m *ExcelMapper) setFieldValue(instance reflect.Value, fieldInfo fieldInfo,
	cellValue string, row int, column string) error {

	field := instance.Field(fieldInfo.index)

	// 处理空值
	if cellValue == "" {
		if fieldInfo.isPointer {
			field.Set(reflect.Zero(field.Type()))
			return nil
		}
		// 非指针类型的空值使用零值
		return nil
	}

	// 转换值
	convertedValue, err := fieldInfo.converter.Convert(cellValue)
	if err != nil {
		mappingErr := MappingError{
			Row:     row,
			Column:  column,
			Field:   fieldInfo.name,
			Value:   cellValue,
			Message: err.Error(),
		}
		m.errors = append(m.errors, mappingErr)
		return mappingErr
	}

	// 设置值
	valueReflect := reflect.ValueOf(convertedValue)
	if fieldInfo.isPointer {
		ptrValue := reflect.New(fieldInfo.fieldType)
		ptrValue.Elem().Set(valueReflect)
		field.Set(ptrValue)
	} else {
		field.Set(valueReflect)
	}

	return nil
}

// isEmptyRow 方法用于处理isEmptyRow相关逻辑。
func (m *ExcelMapper) isEmptyRow(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}

// GetErrors 方法用于处理GetErrors相关逻辑。
func (m *ExcelMapper) GetErrors() []MappingError {
	return m.errors
}

// ClearErrors 方法用于处理ClearErrors相关逻辑。
func (m *ExcelMapper) ClearErrors() {
	m.errors = m.errors[:0]
}

// 值转换器实现
type stringConv struct{}

// Convert 方法用于处理Convert相关逻辑。
func (c *stringConv) Convert(value string) (interface{}, error) { return value, nil }

type intConv struct{}

// Convert 方法用于处理Convert相关逻辑。
func (c *intConv) Convert(value string) (interface{}, error) {
	result, err := strconv.Atoi(value)
	return result, err
}

type int64Conv struct{}

// Convert 方法用于处理Convert相关逻辑。
func (c *int64Conv) Convert(value string) (interface{}, error) {
	result, err := strconv.ParseInt(value, 10, 64)
	return result, err
}

type float64Conv struct{}

// Convert 方法用于处理Convert相关逻辑。
func (c *float64Conv) Convert(value string) (interface{}, error) {
	result, err := strconv.ParseFloat(value, 64)
	return result, err
}

type boolConv struct{}

// Convert 方法用于处理Convert相关逻辑。
func (c *boolConv) Convert(value string) (interface{}, error) {
	value = strings.ToLower(value)
	switch value {
	case "true", "1", "yes", "y", "是", "真":
		return true, nil
	case "false", "0", "no", "n", "否", "假":
		return false, nil
	default:
		return strconv.ParseBool(value)
	}
}

type timeConv struct{}

// Convert 方法用于处理Convert相关逻辑。
func (c *timeConv) Convert(value string) (interface{}, error) {
	// 常见时间格式
	formats := []string{
		"2006-01-02",
		"2006/01/02",
		"2006-01-02 15:04:05",
		"2006/01/02 15:04:05",
		"01/02/2006",
		"02/01/2006",
		"2006年01月02日",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		}
	}

	// 尝试Excel数值日期
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		// Excel日期起始点：1900年1月1日（但Excel错误地认为1900年是闰年）
		excelEpoch := time.Date(1899, 12, 30, 0, 0, 0, 0, ShangHaiTimeLocation)
		return excelEpoch.AddDate(0, 0, int(f)), nil
	}

	return nil, fmt.Errorf("无法解析时间: %s", value)
}
