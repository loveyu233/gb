package gb

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
)

type PaginationParams struct {
	minPage       int
	minSize       int
	maxSize       int
	defaultPage   int
	defaultSize   int
	pageFieldName string
	sizeFieldName string
}

type PaginationParamsOption func(*PaginationParams)

// WithPaginationMinPage 用来设置允许的最小页码。
func WithPaginationMinPage(minPage int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.minPage = minPage
	}
}

// WithPaginationMinSize 用来设置分页的最小条数。
func WithPaginationMinSize(minSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.minSize = minSize
	}
}

// WithPaginationMaxSize 用来限制每页可请求的最大数量。
func WithPaginationMaxSize(maxSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.maxSize = maxSize
	}
}

// WithPaginationDefaultPage 用来配置默认页码。
func WithPaginationDefaultPage(defaultPage int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.defaultPage = defaultPage
	}
}

// WithPaginationDefaultSize 用来配置默认分页大小。
func WithPaginationDefaultSize(defaultSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.defaultSize = defaultSize
	}
}

// WithPaginationPageFieldName 用来自定义页码参数名。
func WithPaginationPageFieldName(pageFieldName string) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.pageFieldName = pageFieldName
	}
}

// WithPaginationSizeFieldName 用来自定义分页大小参数名。
func WithPaginationSizeFieldName(sizeFieldName string) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.sizeFieldName = sizeFieldName
	}
}

// ParsePaginationParams 用来从请求查询参数解析分页信息。
func ParsePaginationParams(c *gin.Context, options ...PaginationParamsOption) (page, size int) {
	var defaultPagination = &PaginationParams{
		defaultPage:   1,
		defaultSize:   10,
		maxSize:       30,
		minSize:       10,
		minPage:       1,
		pageFieldName: "page",
		sizeFieldName: "size",
	}
	for _, opt := range options {
		opt(defaultPagination)
	}

	page = cast.ToInt(c.Query(defaultPagination.pageFieldName))
	if page < defaultPagination.minPage {
		page = defaultPagination.defaultPage
	}

	size = cast.ToInt(c.Query(defaultPagination.sizeFieldName))
	if size < defaultPagination.minSize || size > defaultPagination.maxSize {
		size = defaultPagination.defaultSize
	}

	return page, size
}

// GetGinQueryDefault 用来获取 query 参数并在缺失时返回默认值。
func GetGinQueryDefault[T any](c *gin.Context, key string, defaultValue T) (T, error) {
	value := c.Query(key)

	// 如果参数为空，返回默认值
	if value == "" {
		return defaultValue, nil
	}

	// 转换为指定类型
	result, err := convertToType[T](value)
	if err != nil {
		// 返回可翻译的类型错误
		return defaultValue, CreateTypeError(key, value, err)
	}

	return result, nil
}

// GetGinQueryRequired 用来获取必须存在的 query 参数并转换类型。
func GetGinQueryRequired[T any](c *gin.Context, key string) (T, error) {
	var zero T
	value := c.Query(key)

	// 如果参数为空，返回 CreateRequiredError
	if value == "" {
		return zero, CreateRequiredError(key)
	}

	// 转换为指定类型
	result, err := convertToType[T](value)
	if err != nil {
		// 返回可翻译的类型错误
		return zero, CreateTypeError(key, value, err)
	}

	return result, nil
}

// GetGinPathRequired 用来读取路径参数并转换为指定类型。
func GetGinPathRequired[T any](c *gin.Context, key string) (T, error) {
	var zero T
	value := c.Param(key)

	// 如果参数为空，返回 CreateRequiredError
	if value == "" {
		return zero, CreateRequiredError(key)
	}

	// 转换为指定类型
	result, err := convertToType[T](value)
	if err != nil {
		// 返回可翻译的类型错误
		return zero, CreateTypeError(key, value, err)
	}

	return result, nil
}
