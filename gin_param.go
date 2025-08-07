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

func WithPaginationMinPage(minPage int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.minPage = minPage
	}
}
func WithPaginationMinSize(minSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.minSize = minSize
	}
}

func WithPaginationMaxSize(maxSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.maxSize = maxSize
	}
}
func WithPaginationDefaultPage(defaultPage int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.defaultPage = defaultPage
	}
}
func WithPaginationDefaultSize(defaultSize int) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.defaultSize = defaultSize
	}
}

func WithPaginationPageFieldName(pageFieldName string) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.pageFieldName = pageFieldName
	}
}
func WithPaginationSizeFieldName(sizeFieldName string) PaginationParamsOption {
	return func(p *PaginationParams) {
		p.sizeFieldName = sizeFieldName
	}
}

// ParsePaginationParams 从请求中解析分页参数
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

// GetGinQueryDefault 带默认值的可选参数方法
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

// GetGinQueryRequired 必需的查询参数，不能为空
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

// GetGinPathRequired 必需的查询参数，不能为空
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
