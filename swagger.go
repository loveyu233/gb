package gb

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
)

// APIDoc 表示API文档根结构
type APIDoc struct {
	Swagger     string                `json:"swagger"`
	Info        Info                  `json:"info"`
	Host        string                `json:"host,omitempty"`
	BasePath    string                `json:"basePath,omitempty"`
	Schemes     []string              `json:"schemes,omitempty"`
	Consumes    []string              `json:"consumes,omitempty"`
	Produces    []string              `json:"produces,omitempty"`
	Paths       map[string]PathItem   `json:"paths"`
	Definitions map[string]Definition `json:"definitions"`
}

// Info 表示API文档信息
type Info struct {
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Version     string `json:"version"`
}

// PathItem 表示路径项
type PathItem map[string]Operation

// Operation 表示操作
type Operation struct {
	Summary     string          `json:"summary,omitempty"`
	Description string          `json:"description,omitempty"`
	OperationID string          `json:"operationId,omitempty"`
	Parameters  []Parameter     `json:"parameters,omitempty"`
	Responses   SwaggerResponse `json:"responses"`
	Tags        []string        `json:"tags,omitempty"`
}

// Parameter 表示参数
type Parameter struct {
	Name        string     `json:"name"`
	In          string     `json:"in"` // query, path, body, header, formData
	Description string     `json:"description,omitempty"`
	Required    bool       `json:"required"`
	Type        string     `json:"type,omitempty"`
	Schema      *SchemaRef `json:"schema,omitempty"`
}

// SwaggerResponse 表示响应
type SwaggerResponse map[string]ResponseObject

// ResponseObject 表示响应对象
type ResponseObject struct {
	Description string     `json:"description"`
	Schema      *SchemaRef `json:"schema,omitempty"`
}

// SchemaRef 表示Schema引用
type SchemaRef struct {
	Ref string `json:"$ref,omitempty"`
}

// Definition 表示模型定义
type Definition struct {
	Type       string              `json:"type,omitempty"`
	Properties map[string]Property `json:"properties,omitempty"`
	Required   []string            `json:"required,omitempty"`
}

// Property 表示属性
type Property struct {
	Type        string   `json:"type,omitempty"`
	Description string   `json:"description,omitempty"`
	Format      string   `json:"format,omitempty"`
	Items       *Items   `json:"items,omitempty"`
	Enum        []string `json:"enum,omitempty"`
}

// Items 表示数组项
type Items struct {
	Type string `json:"type,omitempty"`
	Ref  string `json:"$ref,omitempty"`
}

// SwaggerGlobalConfig 表示工具配置
type SwaggerGlobalConfig struct {
	Title       string
	Description string
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	OutputPath  string // 不设置默认 swagger/swagger.json
}

// GlobalParams 表示全局参数配置
type GlobalParams struct {
	PathParams   []SwaggerParamDescription // 全局路径参数
	QueryParams  []SwaggerParamDescription // 全局查询参数
	HeaderParams []SwaggerParamDescription // 全局头部参数
}

var (
	ParamTypeString  = "string"
	ParamTypeInteger = "integer"
	ParamTypeNumber  = "number"
	ParamTypeBoolean = "boolean"
)

// SwaggerParamDescription 表示参数描述
type SwaggerParamDescription struct {
	Name        string // 参数名称
	Description string // 参数描述
	Type        string // 参数类型，默认为string
	Required    bool   // 是否必传
}

// Generator 表示Swagger生成器
type Generator struct {
	Doc          APIDoc
	Config       SwaggerGlobalConfig
	GlobalParams GlobalParams // 新增：全局参数配置
}

// NewSwaggerGenerator 创建一个新的Swagger生成器
func NewSwaggerGenerator(config SwaggerGlobalConfig) *Generator {
	if config.Schemes == nil {
		config.Schemes = []string{"http", "https"}
	}
	if config.OutputPath == "" {
		config.OutputPath = "swagger/swagger.json"
	}

	return &Generator{
		Config: config,
		Doc: APIDoc{
			Swagger: "2.0",
			Info: Info{
				Title:       config.Title,
				Description: config.Description,
				Version:     config.Version,
			},
			Host:        config.Host,
			BasePath:    config.BasePath,
			Schemes:     config.Schemes,
			Consumes:    []string{"application/json"},
			Produces:    []string{"application/json"},
			Paths:       make(map[string]PathItem),
			Definitions: make(map[string]Definition),
		},
		GlobalParams: GlobalParams{}, // 初始化全局参数
	}
}

// SetGlobalParams 设置全局参数
func (g *Generator) SetGlobalParams(params GlobalParams) {
	g.GlobalParams = params
}

// AddGlobalPathParams 添加全局路径参数
func (g *Generator) AddGlobalPathParams(params []SwaggerParamDescription) {
	g.GlobalParams.PathParams = append(g.GlobalParams.PathParams, params...)
}

// AddGlobalQueryParams 添加全局查询参数
func (g *Generator) AddGlobalQueryParams(params []SwaggerParamDescription) {
	g.GlobalParams.QueryParams = append(g.GlobalParams.QueryParams, params...)
}

// AddGlobalHeaderParams 添加全局头部参数
func (g *Generator) AddGlobalHeaderParams(params []SwaggerParamDescription) {
	g.GlobalParams.HeaderParams = append(g.GlobalParams.HeaderParams, params...)
}

// SwaggerAPIInfo 表示API信息 - 修改支持多种类型
type SwaggerAPIInfo struct {
	Path           string
	Method         string
	Summary        string
	Description    string
	Tags           []string
	Request        any               // 支持结构体和[]SwaggerParamDescription
	Response       any               // 支持结构体和[]SwaggerParamDescription
	PathParams     any               // 修改：支持结构体和[]SwaggerParamDescription
	QueryParams    any               // 修改：查询参数，支持结构体和[]SwaggerParamDescription
	HeaderParams   any               // 修改：头部参数，支持结构体和[]SwaggerParamDescription
	ResponseStatus map[string]string // 状态码描述，如 {"200": "成功", "404": "未找到"}
	IgnoreGlobal   bool              // 新增：是否忽略全局参数
}

// AddAPI 添加API - 增强版
func (g *Generator) AddAPI(info SwaggerAPIInfo) {
	method := strings.ToLower(info.Method)

	pathItem, exists := g.Doc.Paths[info.Path]
	if !exists {
		pathItem = make(PathItem)
		g.Doc.Paths[info.Path] = pathItem
	}

	operation := Operation{
		Summary:     info.Summary,
		Description: info.Description,
		Tags:        info.Tags,
		OperationID: generateOperationID(info.Method, info.Path),
		Parameters:  []Parameter{},
		Responses:   make(SwaggerResponse),
	}

	// 0. 添加全局参数（如果未忽略）
	if !info.IgnoreGlobal {
		g.addGlobalParams(&operation)
	}

	// 1. 处理路径参数
	if info.PathParams != nil {
		g.processParams(info.PathParams, "path", &operation)
	} else {
		// 如果没有提供路径参数描述，则自动提取路径中的参数
		pathParams := extractPathParams(info.Path)
		for _, paramName := range pathParams {
			operation.Parameters = append(operation.Parameters, Parameter{
				Name:        paramName,
				In:          "path",
				Description: "",
				Required:    true,
				Type:        "string",
			})
		}
	}

	// 2. 处理查询参数
	if info.QueryParams != nil {
		g.processParams(info.QueryParams, "query", &operation)
	}

	// 3. 处理头部参数
	if info.HeaderParams != nil {
		g.processParams(info.HeaderParams, "header", &operation)
	}

	// 4. 处理请求参数
	if info.Request != nil {
		g.processRequestParams(info.Request, &operation)
	}

	// 5. 处理响应参数
	if info.Response != nil {
		g.processResponseParams(info.Response, &operation)
	}

	// 6. 处理响应状态码
	if len(info.ResponseStatus) > 0 {
		for code, desc := range info.ResponseStatus {
			operation.Responses[code] = ResponseObject{
				Description: desc,
			}
		}
	} else if len(operation.Responses) == 0 {
		// 默认响应
		operation.Responses["200"] = ResponseObject{
			Description: "成功",
		}
	}

	pathItem[method] = operation
	g.Doc.Paths[info.Path] = pathItem
}

// addGlobalParams 添加全局参数到操作中
func (g *Generator) addGlobalParams(operation *Operation) {
	// 添加全局路径参数
	for _, param := range g.GlobalParams.PathParams {
		paramType := param.Type
		if paramType == "" {
			paramType = "string"
		}

		operation.Parameters = append(operation.Parameters, Parameter{
			Name:        param.Name,
			In:          "path",
			Description: param.Description,
			Required:    param.Required,
			Type:        paramType,
		})
	}

	// 添加全局查询参数
	for _, param := range g.GlobalParams.QueryParams {
		paramType := param.Type
		if paramType == "" {
			paramType = "string"
		}

		operation.Parameters = append(operation.Parameters, Parameter{
			Name:        param.Name,
			In:          "query",
			Description: param.Description,
			Required:    param.Required,
			Type:        paramType,
		})
	}

	// 添加全局头部参数
	for _, param := range g.GlobalParams.HeaderParams {
		paramType := param.Type
		if paramType == "" {
			paramType = "string"
		}

		operation.Parameters = append(operation.Parameters, Parameter{
			Name:        param.Name,
			In:          "header",
			Description: param.Description,
			Required:    param.Required,
			Type:        paramType,
		})
	}
}

// processParams 处理参数 - 统一的参数处理方法
func (g *Generator) processParams(params interface{}, paramIn string, operation *Operation) {
	if params == nil {
		return
	}

	// 使用反射检查参数类型
	paramValue := reflect.ValueOf(params)
	paramType := reflect.TypeOf(params)

	// 如果是指针，获取元素
	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
		paramValue = paramValue.Elem()
	}

	// 判断是否为切片类型（[]SwaggerParamDescription）
	if paramType.Kind() == reflect.Slice {
		// 检查切片元素类型
		elemType := paramType.Elem()
		if elemType == reflect.TypeOf(SwaggerParamDescription{}) {
			// 处理 []SwaggerParamDescription 类型
			if paramSlice, ok := params.([]SwaggerParamDescription); ok {
				for _, param := range paramSlice {
					paramType := param.Type
					if paramType == "" {
						paramType = "string"
					}

					// 路径参数总是必需的
					required := param.Required
					if paramIn == "path" {
						required = true
					}

					operation.Parameters = append(operation.Parameters, Parameter{
						Name:        param.Name,
						In:          paramIn,
						Description: param.Description,
						Required:    required,
						Type:        paramType,
					})
				}
			}
		}
	} else if paramType.Kind() == reflect.Struct {
		// 处理结构体类型
		g.processStructAsParams(params, paramIn, operation)
	}
}

// processStructAsParams 将结构体字段处理为参数 - 通用方法
func (g *Generator) processStructAsParams(paramStruct interface{}, paramIn string, operation *Operation) {
	paramType := reflect.TypeOf(paramStruct)
	if paramType.Kind() == reflect.Ptr {
		paramType = paramType.Elem()
	}

	if paramType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < paramType.NumField(); i++ {
		field := paramType.Field(i)

		// 忽略非导出字段
		if field.PkgPath != "" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		// 获取参数名
		paramName := jsonTag
		if paramName == "" {
			paramName = field.Name
		} else {
			paramName = strings.Split(paramName, ",")[0]
		}

		// 获取参数描述
		description := field.Tag.Get("desc")
		if description == "" {
			description = field.Tag.Get("description")
		}

		// 判断是否必需
		required := field.Tag.Get("binding") == "required" || strings.Contains(field.Tag.Get("binding"), "required")

		// 路径参数总是必需的
		if paramIn == "path" {
			required = true
		}

		// 获取参数类型
		paramType := getSwaggerTypeFromReflectType(field.Type)

		operation.Parameters = append(operation.Parameters, Parameter{
			Name:        paramName,
			In:          paramIn,
			Description: description,
			Required:    required,
			Type:        paramType,
		})
	}
}

// processRequestParams 处理请求参数 - 支持多种类型
func (g *Generator) processRequestParams(req interface{}, operation *Operation) {
	reqType := reflect.TypeOf(req)
	if reqType.Kind() == reflect.Ptr {
		reqType = reqType.Elem()
	}

	// 判断是否为切片类型（[]SwaggerParamDescription）
	if reqType.Kind() == reflect.Slice {
		elemType := reqType.Elem()
		if elemType == reflect.TypeOf(SwaggerParamDescription{}) {
			// 处理 []SwaggerParamDescription 类型，作为 formData 参数
			if paramSlice, ok := req.([]SwaggerParamDescription); ok {
				for _, param := range paramSlice {
					paramType := param.Type
					if paramType == "" {
						paramType = "string"
					}

					operation.Parameters = append(operation.Parameters, Parameter{
						Name:        param.Name,
						In:          "formData",
						Description: param.Description,
						Required:    param.Required,
						Type:        paramType,
					})
				}
			}
			return
		}
	}

	// 处理结构体类型
	if reqType.Kind() == reflect.Struct {
		// 为请求定义添加模型
		modelName := reqType.Name()
		g.addModelDefinition(modelName, req)

		// 添加body参数
		operation.Parameters = append(operation.Parameters, Parameter{
			Name:        "body",
			In:          "body",
			Description: "请求参数",
			Required:    true,
			Schema: &SchemaRef{
				Ref: "#/definitions/" + modelName,
			},
		})
	}
}

// processResponseParams 处理响应参数 - 支持多种类型
func (g *Generator) processResponseParams(resp interface{}, operation *Operation) {
	respType := reflect.TypeOf(resp)
	if respType.Kind() == reflect.Ptr {
		respType = respType.Elem()
	}

	// 判断是否为切片类型（[]SwaggerParamDescription）
	if respType.Kind() == reflect.Slice {
		elemType := respType.Elem()
		if elemType == reflect.TypeOf(SwaggerParamDescription{}) {
			// 对于响应参数，[]SwaggerParamDescription 可能不太常用
			// 这里可以根据需要处理，比如作为响应头部信息
			return
		}
	}

	// 处理结构体类型
	if respType.Kind() == reflect.Struct {
		// 为响应定义添加模型
		modelName := respType.Name()
		g.addModelDefinition(modelName, resp)

		// 添加成功响应
		operation.Responses["200"] = ResponseObject{
			Description: "成功",
			Schema: &SchemaRef{
				Ref: "#/definitions/" + modelName,
			},
		}
	}
}

// getSwaggerTypeFromReflectType 从反射类型获取Swagger类型
func getSwaggerTypeFromReflectType(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	switch t.Kind() {
	case reflect.Bool:
		return "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return "integer"
	case reflect.Int64, reflect.Uint64:
		return "integer"
	case reflect.Float32, reflect.Float64:
		return "number"
	case reflect.String:
		return "string"
	case reflect.Slice, reflect.Array:
		return "array"
	default:
		return "string"
	}
}

// 提取路径中的参数名
func extractPathParams(path string) []string {
	var params []string
	re := regexp.MustCompile(`{([^}]+)}`)
	matches := re.FindAllStringSubmatch(path, -1)

	for _, match := range matches {
		if len(match) > 1 {
			params = append(params, match[1])
		}
	}

	return params
}

// addModelDefinition 添加模型定义
func (g *Generator) addModelDefinition(name string, model interface{}) {
	if _, exists := g.Doc.Definitions[name]; exists {
		return
	}

	modelType := reflect.TypeOf(model)
	if modelType.Kind() == reflect.Ptr {
		modelType = modelType.Elem()
	}

	definition := Definition{
		Type:       "object",
		Properties: make(map[string]Property),
		Required:   []string{},
	}

	if modelType.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < modelType.NumField(); i++ {
		field := modelType.Field(i)

		// 忽略非导出字段
		if field.PkgPath != "" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		jsonName := jsonTag
		if jsonName == "" {
			jsonName = field.Name
		} else {
			jsonName = strings.Split(jsonName, ",")[0]
		}

		// 处理required标签
		required := field.Tag.Get("binding") == "required" || strings.Contains(field.Tag.Get("binding"), "required")
		if required {
			definition.Required = append(definition.Required, jsonName)
		}

		// 获取字段描述
		description := field.Tag.Get("desc")

		// 添加属性
		property := Property{
			Description: description,
		}

		g.setPropertyType(&property, field.Type)
		definition.Properties[jsonName] = property

		// 如果是嵌套结构体，递归添加
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct && fieldType.Name() != "Time" {
			// 创建结构体实例
			structValue := reflect.New(fieldType).Interface()
			g.addModelDefinition(fieldType.Name(), structValue)
		}
	}

	g.Doc.Definitions[name] = definition
}

// setPropertyType 设置属性类型
func (g *Generator) setPropertyType(property *Property, t reflect.Type) {
	kind := t.Kind()

	if kind == reflect.Ptr {
		g.setPropertyType(property, t.Elem())
		return
	}

	switch kind {
	case reflect.Bool:
		property.Type = "boolean"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		property.Type = "integer"
		property.Format = "int32"
	case reflect.Int64, reflect.Uint64:
		property.Type = "integer"
		property.Format = "int64"
	case reflect.Float32:
		property.Type = "number"
		property.Format = "float"
	case reflect.Float64:
		property.Type = "number"
		property.Format = "double"
	case reflect.String:
		property.Type = "string"
	case reflect.Struct:
		if t.Name() == "Time" {
			property.Type = "string"
			property.Format = "date-time"
		} else {
			// 创建嵌套对象的引用
			property.Type = "object"
			// 由于Property结构中没有直接的Ref字段，我们应该使用Items来引用其他定义
			if property.Items == nil {
				property.Items = &Items{
					Ref: "#/definitions/" + t.Name(),
				}
			} else {
				property.Items.Ref = "#/definitions/" + t.Name()
			}
		}
	case reflect.Slice, reflect.Array:
		property.Type = "array"
		property.Items = &Items{}
		elemType := t.Elem()

		if elemType.Kind() == reflect.Struct && elemType.Name() != "Time" {
			property.Items.Ref = "#/definitions/" + elemType.Name()

			// 创建元素类型的实例并添加定义
			structValue := reflect.New(elemType).Interface()
			g.addModelDefinition(elemType.Name(), structValue)
		} else {
			switch elemType.Kind() {
			case reflect.Bool:
				property.Items.Type = "boolean"
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
				property.Items.Type = "integer"
			case reflect.Int64, reflect.Uint64:
				property.Items.Type = "integer"
			case reflect.Float32, reflect.Float64:
				property.Items.Type = "number"
			case reflect.String:
				property.Items.Type = "string"
			default:
				panic("不支持的类型")
			}
		}
	case reflect.Map:
		property.Type = "object"
		// 注意：Swagger 2.0 不完全支持Map的详细定义
	default:
		panic("不支持的类型")
	}
}

// Generate 生成Swagger文档
func (g *Generator) Generate() error {
	// 创建输出目录（如果不存在）
	dir := filepath.Dir(g.Config.OutputPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("创建输出目录失败: %v", err)
	}

	// 将API文档转换为JSON
	data, err := json.MarshalIndent(g.Doc, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化文档失败: %v", err)
	}

	// 写入文件
	if err := os.WriteFile(g.Config.OutputPath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// 生成操作ID
func generateOperationID(method, path string) string {
	// 去除路径中的参数部分
	pathParts := strings.Split(path, "/")
	for i, part := range pathParts {
		if strings.HasPrefix(part, "{") && strings.HasSuffix(part, "}") {
			pathParts[i] = "by" + strings.Title(part[1:len(part)-1])
		}
	}

	// 构建操作ID
	cleanPath := strings.Join(pathParts, "_")
	cleanPath = strings.Trim(cleanPath, "_")
	cleanPath = strings.ReplaceAll(cleanPath, "-", "_")

	return strings.ToLower(method) + "_" + cleanPath
}

// AddPath 添加单个API路径的快捷方法
func (g *Generator) AddPath(path, method, summary, description string, tags []string) *Operation {
	method = strings.ToLower(method)

	pathItem, exists := g.Doc.Paths[path]
	if !exists {
		pathItem = make(PathItem)
		g.Doc.Paths[path] = pathItem
	}

	operation := Operation{
		Summary:     summary,
		Description: description,
		Tags:        tags,
		OperationID: generateOperationID(method, path),
		Parameters:  []Parameter{},
		Responses:   make(SwaggerResponse),
	}

	// 添加全局参数
	g.addGlobalParams(&operation)

	// 处理路径参数
	pathParams := extractPathParams(path)
	for _, paramName := range pathParams {
		operation.Parameters = append(operation.Parameters, Parameter{
			Name:     paramName,
			In:       "path",
			Required: true,
			Type:     "string",
		})
	}

	// 添加默认响应
	operation.Responses["200"] = ResponseObject{
		Description: "成功",
	}

	pathItem[method] = operation
	return &operation
}

// AddPathParamDesc 为已有操作添加路径参数描述
func (g *Generator) AddPathParamDesc(path, method, paramName, description string, paramType string) error {
	method = strings.ToLower(method)

	pathItem, exists := g.Doc.Paths[path]
	if !exists {
		return fmt.Errorf("路径 %s 不存在", path)
	}

	operation, exists := pathItem[method]
	if !exists {
		return fmt.Errorf("路径 %s 的方法 %s 不存在", path, method)
	}

	// 查找并更新参数描述
	paramFound := false
	for i, param := range operation.Parameters {
		if param.In == "path" && param.Name == paramName {
			operation.Parameters[i].Description = description
			if paramType != "" {
				operation.Parameters[i].Type = paramType
			}
			paramFound = true
			break
		}
	}

	// 如果未找到参数，添加一个新参数
	if !paramFound {
		if paramType == "" {
			paramType = "string"
		}

		operation.Parameters = append(operation.Parameters, Parameter{
			Name:        paramName,
			In:          "path",
			Description: description,
			Required:    true,
			Type:        paramType,
		})
	}

	// 更新操作
	pathItem[method] = operation
	g.Doc.Paths[path] = pathItem

	return nil
}
