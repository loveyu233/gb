package gb

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"strings"
)

// ParseIDFromUrl 从请求中解析ID
func ParseIDFromUrl(c *gin.Context) (int64, error) {
	return ParseInt64FromUrl(c, "id")
}

// ParseInt64FromUrl 从请求中解析int64字段
func ParseInt64FromUrl(c *gin.Context, paramName string) (int64, error) {
	return strconv.ParseInt(c.Param(paramName), 10, 64)
}

// ParseIntFromQueryString 从请求中解析int字段
func ParseIntFromQueryString(c *gin.Context, fieldName string) (int, error) {
	return strconv.Atoi(c.DefaultQuery(fieldName, "0"))
}

// ParseInt64FromQueryString 从请求中解析int64字段
func ParseInt64FromQueryString(c *gin.Context, fieldName string) (int64, error) {
	return strconv.ParseInt(c.DefaultQuery(fieldName, "0"), 10, 64)
}

// ParseStringFromQueryString 从请求中解析string字段
func ParseStringFromQueryString(c *gin.Context, fieldName string) string {
	return strings.TrimSpace(c.Query(fieldName))
}

// ParsePaginationParams 从请求中解析分页参数
func ParsePaginationParams(c *gin.Context, defaultPagination ...map[string]int) (page, size int) {
	defaultPage := 1
	defaultSize := 10

	if len(defaultPagination) > 0 {
		defaultPage = defaultPagination[0]["page"]
		defaultSize = defaultPagination[0]["size"]
	}

	page, err := ParseIntFromQueryString(c, "page")
	if err != nil || page <= 0 {
		page = defaultPage
	}

	size, err = ParseIntFromQueryString(c, "size")
	if err != nil || size <= 0 || size >= 100 {
		size = defaultSize
	}

	return page, size
}
