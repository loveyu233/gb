package gb

import (
	"fmt"

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

// GetGinQuery 如果为空,有默认值返回默认值,没有默认值返回errStr错误,errStr为空则使用%s不能为空错误
func GetGinQuery[T any](c *gin.Context, key, errStr string, defaultValue ...T) (T, error) {
	var zero T
	if errStr == "" {
		errStr = fmt.Sprintf("%s不能为空", key)
	}
	paramStr := c.Query(key)
	if paramStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return zero, ErrInvalidParam.WithMessage(errStr)
	}
	return convertToType[T](paramStr)
}

// GetGinPath 返回路径参数key的值
func GetGinPath[T any](c *gin.Context, key string) (T, error) {
	paramStr := c.Param(key)
	if paramStr == "" {
		var zero T
		return zero, ErrInvalidParam.WithMessage(fmt.Sprintf("%s不存在", key))
	}
	return convertToType[T](paramStr)
}
