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

// WithPaginationMinPage 函数用于处理WithPaginationMinPage相关逻辑。
func WithPaginationMinPage(minPage int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.minPage = minPage
	}
}

// WithPaginationMinSize 函数用于处理WithPaginationMinSize相关逻辑。
func WithPaginationMinSize(minSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.minSize = minSize
	}
}

// WithPaginationMaxSize 函数用于处理WithPaginationMaxSize相关逻辑。
func WithPaginationMaxSize(maxSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.maxSize = maxSize
	}
}

// WithPaginationDefaultPage 函数用于处理WithPaginationDefaultPage相关逻辑。
func WithPaginationDefaultPage(defaultPage int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.defaultPage = defaultPage
	}
}

// WithPaginationDefaultSize 函数用于处理WithPaginationDefaultSize相关逻辑。
func WithPaginationDefaultSize(defaultSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.defaultSize = defaultSize
	}
}

// WithPaginationPageFieldName 函数用于处理WithPaginationPageFieldName相关逻辑。
func WithPaginationPageFieldName(pageFieldName string) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.pageFieldName = pageFieldName
	}
}

// WithPaginationSizeFieldName 函数用于处理WithPaginationSizeFieldName相关逻辑。
func WithPaginationSizeFieldName(sizeFieldName string) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.sizeFieldName = sizeFieldName
	}
}

// ParsePaginationParams 函数用于处理ParsePaginationParams相关逻辑。
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

// GetGinQueryDefault 函数用于处理GetGinQueryDefault相关逻辑。
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

// GetGinQueryRequired 函数用于处理GetGinQueryRequired相关逻辑。
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

// GetGinPathRequired 函数用于处理GetGinPathRequired相关逻辑。
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
