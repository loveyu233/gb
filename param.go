package gb

import (
	"fmt"
	"strconv"

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

func ParserString(s string) (string, error) {
	return s, nil
}

func ParserInt(s string) (int, error) {
	return strconv.Atoi(s)
}

func ParserInt64(s string) (int64, error) {
	return strconv.ParseInt(s, 10, 64)
}

func ParserFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, 64)
}

func ParserBool(s string) (bool, error) {
	return strconv.ParseBool(s)
}

// ParseFromQuery 有默认值使用默认值而不是返回错误,errStr为空则会只用key不能为空作为错误返回
func ParseFromQuery[T any](c *gin.Context, key, errStr string, parser func(string) (T, error), defaultValue ...T) (T, error) {
	if errStr == "" {
		errStr = fmt.Sprintf("%s不能为空", key)
	}
	paramStr := c.Query(key)
	if paramStr == "" {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		var zero T
		return zero, ErrInvalidParam.WithMessage(errStr)
	}

	val, err := parser(paramStr)
	if err != nil {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		var zero T
		return zero, ErrInvalidParam.WithMessage(fmt.Sprintf("%s类型错误", key))
	}

	return val, nil
}

// ParseFromPath 有默认值使用默认值而不是返回错误,errStr为空则会只用key不能为空作为错误返回
func ParseFromPath[T any](c *gin.Context, key string, parser func(string) (T, error)) (T, error) {
	paramStr := c.Param(key)
	if paramStr == "" {
		var zero T
		return zero, ErrInvalidParam.WithMessage(fmt.Sprintf("%s不存在", key))
	}

	val, err := parser(paramStr)
	if err != nil {
		var zero T
		return zero, ErrInvalidParam.WithMessage(fmt.Sprintf("%s类型错误", key))
	}

	return val, nil
}
