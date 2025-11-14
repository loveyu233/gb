package gb

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
	"unsafe"

	"github.com/xuri/excelize/v2"
)

// ExcelExporter Excel数据导出器
type ExcelExporter struct {
	// 配置选项
	SheetName        string             // 工作表名称
	IncludeHeaderPtr *bool              // 是否包含表头
	includeHeader    bool               // 是否包含表头,方法内使用
	HeaderRow        int                // 表头行号 (1-based)
	DataStartRow     int                // 数据开始行号 (1-based)
	ColumnWidths     map[string]float64 // 列宽设置 map[列名]宽度
	HeaderStyle      *HeaderStyle       // 表头样式
	DataStyle        *DataStyle         // 数据样式

	// 内部缓存
	structCache map[reflect.Type]*exportStructInfo
	cacheMutex  sync.RWMutex

	// 性能统计
	exportedRows int
	exportedCols int
}

// HeaderStyle 表头样式
type HeaderStyle struct {
	Bold            bool
	BackgroundColor string // 十六进制颜色，如 "FFFF00"
	FontColor       string
	FontSize        int
	Alignment       string // "left", "center", "right"
}

// DataStyle 数据样式
type DataStyle struct {
	FontSize     int
	Alignment    string
	NumberFormat string // 数字格式，如 "0.00", "#,##0"
	DateFormat   string // 日期格式，如 "yyyy-mm-dd"
}

// exportStructInfo 导出结构体信息
type exportStructInfo struct {
	fields []exportFieldInfo
}

// exportFieldInfo 导出字段信息
type exportFieldInfo struct {
	index       int            // 字段索引
	name        string         // 字段名
	tag         string         // excel标签值
	columnTitle string         // 列标题
	fieldType   reflect.Type   // 字段类型
	isPointer   bool           // 是否为指针类型
	formatter   valueFormatter // 值格式化器
	columnIndex int            // Excel列索引
}

// valueFormatter 值格式化器接口
type valueFormatter interface {
	Format(value interface{}) string
}

// 预定义格式化器
var (
	stringFormatter  = &stringFormat{}
	intFormatter     = &intFormat{}
	floatFormatter   = &floatFormat{}
	boolFormatter    = &boolFormat{}
	timeFormatter    = &timeFormat{}
	pointerFormatter = &pointerFormat{}
)

type WithExcelExporterOption func(*ExcelExporter)

// WithExcelExporterSheetName 函数用于处理WithExcelExporterSheetName相关逻辑。
func WithExcelExporterSheetName(sheetName string) WithExcelExporterOption {
	return func(e *ExcelExporter) {
		e.SheetName = sheetName
	}
}

// WithExcelExporterHeaderRow 函数用于处理WithExcelExporterHeaderRow相关逻辑。
func WithExcelExporterHeaderRow(headerRow int) WithExcelExporterOption {
	return func(e *ExcelExporter) {
		e.HeaderRow = headerRow
	}
}

// WithExcelExporterDataStartRow 函数用于处理WithExcelExporterDataStartRow相关逻辑。
func WithExcelExporterDataStartRow(dataStartRow int) WithExcelExporterOption {
	return func(e *ExcelExporter) {
		e.DataStartRow = dataStartRow
	}
}

// WithExcelExporterIncludeHeader 函数用于处理WithExcelExporterIncludeHeader相关逻辑。
func WithExcelExporterIncludeHeader(includeHeader *bool) WithExcelExporterOption {
	return func(e *ExcelExporter) {
		e.IncludeHeaderPtr = includeHeader
	}
}

// WithExcelExporterColumnWidths 函数用于处理WithExcelExporterColumnWidths相关逻辑。
func WithExcelExporterColumnWidths(columnWidths map[string]float64) WithExcelExporterOption {
	return func(e *ExcelExporter) {
		e.ColumnWidths = columnWidths
	}
}

// WithExcelExporterHeaderStyle 函数用于处理WithExcelExporterHeaderStyle相关逻辑。
func WithExcelExporterHeaderStyle(headerStyle *HeaderStyle) WithExcelExporterOption {
	return func(e *ExcelExporter) {
		e.HeaderStyle = headerStyle
	}
}

// WithExcelExporterDataStyle 函数用于处理WithExcelExporterDataStyle相关逻辑。
func WithExcelExporterDataStyle(dataStyle *DataStyle) WithExcelExporterOption {
	return func(e *ExcelExporter) {
		e.DataStyle = dataStyle
	}
}

// InitExcelExporter 函数用于处理InitExcelExporter相关逻辑。
func InitExcelExporter(options ...WithExcelExporterOption) *ExcelExporter {
	excelExporter := new(ExcelExporter)
	for i := range options {
		options[i](excelExporter)
	}
	if excelExporter.SheetName == "" {
		excelExporter.SheetName = "Sheet1"
	}

	if excelExporter.HeaderRow == 0 {
		excelExporter.HeaderRow = 1
	}

	if excelExporter.DataStartRow == 0 {
		excelExporter.DataStartRow = 2
	}

	if excelExporter.IncludeHeaderPtr == nil {
		excelExporter.includeHeader = true
	} else {
		excelExporter.includeHeader = *excelExporter.IncludeHeaderPtr
	}

	if excelExporter.HeaderStyle == nil {
		excelExporter.HeaderStyle = &HeaderStyle{
			Bold:            true,
			BackgroundColor: "D3D3D3", // 浅灰色
			FontColor:       "000000", // 黑色
			FontSize:        12,
			Alignment:       "center",
		}
	}

	if excelExporter.DataStyle == nil {
		excelExporter.DataStyle = &DataStyle{
			FontSize:     11,
			Alignment:    "left",
			NumberFormat: "General",
			DateFormat:   "yyyy-mm-dd",
		}
	}

	excelExporter.structCache = make(map[reflect.Type]*exportStructInfo)

	return excelExporter
}

// ExportToFile 方法用于处理ExportToFile相关逻辑。
func (e *ExcelExporter) ExportToFile(data interface{}, filePath string) error {
	// 创建Excel文件
	file := excelize.NewFile()

	// 导出到工作表
	err := e.ExportToSheet(data, file, e.SheetName)
	if err != nil {
		return err
	}
	if e.SheetName != "Sheet1" {
		file.DeleteSheet("Sheet1")
	}
	fileNameType := GetFileNameType(filePath)
	if fileNameType != "xlsx" {
		filePath += ".xlsx"
	}
	// 保存文件
	return file.SaveAs(filePath)
}

// ExportToBuffer 方法用于处理ExportToBuffer相关逻辑。
func (e *ExcelExporter) ExportToBuffer(data interface{}) (*bytes.Buffer, error) {
	// 创建Excel文件
	file := excelize.NewFile()

	// 导出到工作表
	err := e.ExportToSheet(data, file, e.SheetName)
	if err != nil {
		return nil, err
	}
	if e.SheetName != "Sheet1" {
		file.DeleteSheet("Sheet1")
	}

	return file.WriteToBuffer()
}

// ExportToExcelizeFile 方法用于处理ExportToExcelizeFile相关逻辑。
func (e *ExcelExporter) ExportToExcelizeFile(data interface{}) (*excelize.File, error) {
	// 创建Excel文件
	file := excelize.NewFile()

	// 导出到工作表
	err := e.ExportToSheet(data, file, e.SheetName)
	if err != nil {
		return nil, err
	}
	if e.SheetName != "Sheet1" {
		file.DeleteSheet("Sheet1")
	}

	return file, nil
}

// ExportToSheet 方法用于处理ExportToSheet相关逻辑。
func (e *ExcelExporter) ExportToSheet(data interface{}, file *excelize.File, sheetName string) error {
	// 参数验证
	dataValue := reflect.ValueOf(data)
	if dataValue.Kind() != reflect.Slice && dataValue.Kind() != reflect.Array {
		return fmt.Errorf("数据必须是切片或数组")
	}

	if dataValue.Len() == 0 {
		return fmt.Errorf("数据为空")
	}

	elemType := dataValue.Type().Elem()
	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	// 获取结构体信息
	structInfo, err := e.getExportStructInfo(elemType)
	if err != nil {
		return err
	}

	// 创建或获取工作表
	if sheetName != "Sheet1" {
		index, err := file.NewSheet(sheetName)
		if err != nil {
			return fmt.Errorf("创建工作表失败: %v", err)
		}
		file.SetActiveSheet(index)
	}

	// 写入表头
	currentRow := e.HeaderRow
	if e.includeHeader {
		err = e.writeHeaders(file, sheetName, structInfo, currentRow)
		if err != nil {
			return err
		}
		currentRow = e.DataStartRow
	}

	// 批量写入数据
	err = e.writeData(file, sheetName, dataValue, structInfo, currentRow)
	if err != nil {
		return err
	}

	// 应用列宽
	e.applyColumnWidths(file, sheetName, structInfo)

	// 设置统计信息
	e.exportedRows = dataValue.Len()
	e.exportedCols = len(structInfo.fields)

	return nil
}

// getExportStructInfo 方法用于处理getExportStructInfo相关逻辑。
func (e *ExcelExporter) getExportStructInfo(elemType reflect.Type) (*exportStructInfo, error) {
	e.cacheMutex.RLock()
	if info, exists := e.structCache[elemType]; exists {
		e.cacheMutex.RUnlock()
		return info, nil
	}
	e.cacheMutex.RUnlock()

	e.cacheMutex.Lock()
	defer e.cacheMutex.Unlock()

	// 双重检查
	if info, exists := e.structCache[elemType]; exists {
		return info, nil
	}

	if elemType.Kind() != reflect.Struct {
		return nil, fmt.Errorf("元素类型必须为struct")
	}

	info := &exportStructInfo{
		fields: make([]exportFieldInfo, 0, elemType.NumField()),
	}

	columnIndex := 0
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		tag := field.Tag.Get("excel")

		if tag == "" || tag == "-" {
			continue
		}

		// 解析标签，支持 "列名" 或 "列名,title:显示标题"
		parts := strings.Split(tag, ",")
		columnName := strings.TrimSpace(parts[0])
		columnTitle := columnName

		// 解析额外选项
		for _, part := range parts[1:] {
			if strings.HasPrefix(part, "title:") {
				columnTitle = strings.TrimSpace(strings.TrimPrefix(part, "title:"))
			}
		}

		fieldType := field.Type
		isPointer := fieldType.Kind() == reflect.Ptr
		if isPointer {
			fieldType = fieldType.Elem()
		}

		formatter := e.getFormatter(fieldType, isPointer)

		info.fields = append(info.fields, exportFieldInfo{
			index:       i,
			name:        field.Name,
			tag:         columnName,
			columnTitle: columnTitle,
			fieldType:   fieldType,
			isPointer:   isPointer,
			formatter:   formatter,
			columnIndex: columnIndex,
		})

		columnIndex++
	}

	e.structCache[elemType] = info
	return info, nil
}

// getFormatter 方法用于处理getFormatter相关逻辑。
func (e *ExcelExporter) getFormatter(fieldType reflect.Type, isPointer bool) valueFormatter {
	if isPointer {
		return pointerFormatter
	}

	switch fieldType.Kind() {
	case reflect.String:
		return stringFormatter
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intFormatter
	case reflect.Float32, reflect.Float64:
		return floatFormatter
	case reflect.Bool:
		return boolFormatter
	}

	// 检查是否为time.Time
	if fieldType == reflect.TypeOf(time.Time{}) {
		return timeFormatter
	}

	return stringFormatter // 默认转换为字符串
}

// writeHeaders 方法用于处理writeHeaders相关逻辑。
func (e *ExcelExporter) writeHeaders(file *excelize.File, sheetName string,
	structInfo *exportStructInfo, row int) error {

	// 创建表头样式
	headerStyleID, err := e.createHeaderStyle(file)
	if err != nil {
		return err
	}

	// 写入表头数据
	for _, fieldInfo := range structInfo.fields {
		cellName := e.getCellName(row, fieldInfo.columnIndex)

		// 设置单元格值
		err = file.SetCellValue(sheetName, cellName, fieldInfo.columnTitle)
		if err != nil {
			return fmt.Errorf("设置标题单元格失败 %s: %v", cellName, err)
		}

		// 应用表头样式
		if headerStyleID != 0 {
			err = file.SetCellStyle(sheetName, cellName, cellName, headerStyleID)
			if err != nil {
				return fmt.Errorf("设置页眉样式失败: %v", err)
			}
		}
	}

	return nil
}

// writeData 方法用于处理writeData相关逻辑。
func (e *ExcelExporter) writeData(file *excelize.File, sheetName string,
	dataValue reflect.Value, structInfo *exportStructInfo, startRow int) error {

	// 创建数据样式
	dataStyleID, err := e.createDataStyle(file)
	if err != nil {
		return err
	}

	// 预分配单元格数据
	batchSize := 1000 // 批量处理大小
	cellData := make([][]interface{}, 0, batchSize)

	for i := 0; i < dataValue.Len(); i++ {
		item := dataValue.Index(i)
		if item.Kind() == reflect.Ptr {
			if item.IsNil() {
				continue
			}
			item = item.Elem()
		}

		rowData := make([]interface{}, len(structInfo.fields))

		// 提取字段值
		for j, fieldInfo := range structInfo.fields {
			fieldValue := item.Field(fieldInfo.index)
			formattedValue := fieldInfo.formatter.Format(fieldValue.Interface())
			rowData[j] = formattedValue
		}

		cellData = append(cellData, rowData)

		// 批量写入
		if len(cellData) >= batchSize || i == dataValue.Len()-1 {
			err = e.writeBatch(file, sheetName, cellData, structInfo,
				startRow+i-len(cellData)+1, dataStyleID)
			if err != nil {
				return err
			}
			cellData = cellData[:0] // 重置切片
		}
	}

	return nil
}

// writeBatch 方法用于处理writeBatch相关逻辑。
func (e *ExcelExporter) writeBatch(file *excelize.File, sheetName string,
	cellData [][]interface{}, structInfo *exportStructInfo, startRow int, styleID int) error {

	for i, rowData := range cellData {
		currentRow := startRow + i

		for j, value := range rowData {
			cellName := e.getCellName(currentRow, j)

			// 设置单元格值
			err := file.SetCellValue(sheetName, cellName, value)
			if err != nil {
				return fmt.Errorf("设置单元格失败 %s: %v", cellName, err)
			}

			// 应用数据样式
			if styleID != 0 {
				err = file.SetCellStyle(sheetName, cellName, cellName, styleID)
				if err != nil {
					return fmt.Errorf("设置数据样式失败: %v", err)
				}
			}
		}
	}

	return nil
}

// getCellName 方法用于处理getCellName相关逻辑。
func (e *ExcelExporter) getCellName(row, col int) string {
	return e.getColumnName(col) + strconv.Itoa(row)
}

// getColumnName 方法用于处理getColumnName相关逻辑。
func (e *ExcelExporter) getColumnName(col int) string {
	if col < 26 {
		return string(rune('A' + col))
	}

	var result []byte
	for col >= 0 {
		remainder := col % 26
		result = append([]byte{byte('A' + remainder)}, result...)
		col = col/26 - 1
		if col < 0 {
			break
		}
	}

	return *(*string)(unsafe.Pointer(&result))
}

// createHeaderStyle 方法用于处理createHeaderStyle相关逻辑。
func (e *ExcelExporter) createHeaderStyle(file *excelize.File) (int, error) {
	if e.HeaderStyle == nil {
		return 0, nil
	}

	style := &excelize.Style{
		Font: &excelize.Font{
			Bold:   e.HeaderStyle.Bold,
			Size:   float64(e.HeaderStyle.FontSize),
			Color:  e.HeaderStyle.FontColor,
			Family: "Arial",
		},
		Alignment: &excelize.Alignment{
			Horizontal: e.HeaderStyle.Alignment,
			Vertical:   "center",
		},
	}

	if e.HeaderStyle.BackgroundColor != "" {
		style.Fill = excelize.Fill{
			Type:    "pattern",
			Color:   []string{e.HeaderStyle.BackgroundColor},
			Pattern: 1,
		}
	}

	return file.NewStyle(style)
}

// createDataStyle 方法用于处理createDataStyle相关逻辑。
func (e *ExcelExporter) createDataStyle(file *excelize.File) (int, error) {
	if e.DataStyle == nil {
		return 0, nil
	}

	style := &excelize.Style{
		Font: &excelize.Font{
			Size:   float64(e.DataStyle.FontSize),
			Family: "Arial",
		},
		Alignment: &excelize.Alignment{
			Horizontal: e.DataStyle.Alignment,
			Vertical:   "center",
		},
		NumFmt: 0, // 使用默认格式
	}

	return file.NewStyle(style)
}

// applyColumnWidths 方法用于处理applyColumnWidths相关逻辑。
func (e *ExcelExporter) applyColumnWidths(file *excelize.File, sheetName string,
	structInfo *exportStructInfo) {

	for _, fieldInfo := range structInfo.fields {
		columnName := e.getColumnName(fieldInfo.columnIndex)

		// 检查是否有自定义宽度
		if width, exists := e.ColumnWidths[fieldInfo.tag]; exists {
			file.SetColWidth(sheetName, columnName, columnName, width)
			continue
		}

		// 根据字段类型设置默认宽度
		defaultWidth := e.getDefaultColumnWidth(fieldInfo.fieldType)
		file.SetColWidth(sheetName, columnName, columnName, defaultWidth)
	}
}

// getDefaultColumnWidth 方法用于处理getDefaultColumnWidth相关逻辑。
func (e *ExcelExporter) getDefaultColumnWidth(fieldType reflect.Type) float64 {
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32:
		return 10.0
	case reflect.Int64:
		return 15.0
	case reflect.Float32, reflect.Float64:
		return 12.0
	case reflect.Bool:
		return 8.0
	case reflect.String:
		return 20.0
	}

	if fieldType == reflect.TypeOf(time.Time{}) {
		return 18.0
	}

	return 15.0 // 默认宽度
}

// GetStats 方法用于处理GetStats相关逻辑。
func (e *ExcelExporter) GetStats() (rows, cols int) {
	return e.exportedRows, e.exportedCols
}

// 值格式化器实现
type stringFormat struct{}

// Format 方法用于处理Format相关逻辑。
func (f *stringFormat) Format(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

type intFormat struct{}

// Format 方法用于处理Format相关逻辑。
func (f *intFormat) Format(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%d", value)
}

type floatFormat struct{}

// Format 方法用于处理Format相关逻辑。
func (f *floatFormat) Format(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%.2f", value)
}

type boolFormat struct{}

// Format 方法用于处理Format相关逻辑。
func (f *boolFormat) Format(value interface{}) string {
	if value == nil {
		return ""
	}
	if b, ok := value.(bool); ok {
		if b {
			return "是"
		}
		return "否"
	}
	return fmt.Sprintf("%v", value)
}

type timeFormat struct{}

// Format 方法用于处理Format相关逻辑。
func (f *timeFormat) Format(value interface{}) string {
	if value == nil {
		return ""
	}
	if t, ok := value.(time.Time); ok {
		if t.IsZero() {
			return ""
		}
		return t.In(ShangHaiTimeLocation).Format("2006-01-02 15:04:05")
	}
	return fmt.Sprintf("%v", value)
}

type pointerFormat struct{}

// Format 方法用于处理Format相关逻辑。
func (f *pointerFormat) Format(value interface{}) string {
	rv := reflect.ValueOf(value)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ""
	}

	actual := rv.Elem().Interface()

	// 根据实际类型选择格式化器
	switch rv.Elem().Kind() {
	case reflect.String:
		return stringFormatter.Format(actual)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intFormatter.Format(actual)
	case reflect.Float32, reflect.Float64:
		return floatFormatter.Format(actual)
	case reflect.Bool:
		return boolFormatter.Format(actual)
	default:
		if rv.Elem().Type() == reflect.TypeOf(time.Time{}) {
			return timeFormatter.Format(actual)
		}
		return stringFormatter.Format(actual)
	}
}
