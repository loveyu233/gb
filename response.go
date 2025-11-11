package gb

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// AppError 自定义错误类型
type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("错误码: %d, 错误信息: %s", e.Code, e.Message)
}

func (e *AppError) WithMessage(format string, args ...any) *AppError {
	if len(args) > 0 {
		format = fmt.Sprintf(format, args...)
	}
	if format == "" {
		format = e.Message
	}
	newErr := NewAppError(e.Code, format)
	return newErr
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// 预定义错误 http状态码 + 业务错误码
var (
	// 100xxx 请求外部服务失败
	ErrRequestExternalService = NewAppError(100000, "请求外部服务失败")
	ErrRequestWechat          = NewAppError(100001, "请求wechat服务失败")
	ErrRequestWechatPay       = NewAppError(100002, "请求wechat支付服务失败")
	ErrRequestAli             = NewAppError(100003, "请求zfb服务失败")
	ErrRequestAliPay          = NewAppError(100004, "请求zfb支付服务失败")

	// 400xxx 客户端错误
	ErrBadRequest   = NewAppError(400000, "请求错误")
	ErrInvalidParam = NewAppError(400001, "请求参数错误")
	ErrTokenInvalid = NewAppError(400002, "token验证失败")

	// 401xxx 未授权
	ErrUnauthorized = NewAppError(401000, "用户未登录或token已失效")

	// 403xxx 禁止操作
	ErrForbiddenAuth = NewAppError(403000, "权限不足")
	ErrUserDisabled  = NewAppError(403001, "用户不存在或已被禁用")

	// 404xxx 数据不存在
	ErrNotFound = NewAppError(404000, "数据不存在")

	// 409xxx 数据已存在
	ErrDataExists          = NewAppError(409000, "数据已存在")
	ErrUniqueIndexConflict = NewAppError(409001, "索引冲突")

	// 5xxxxx 服务器错误
	ErrServerBusy = NewAppError(500000, "服务器繁忙")
	ErrDatabase   = NewAppError(500001, "数据库错误")
	ErrRedis      = NewAppError(500002, "redis错误")

	EncryptErr = NewAppError(600000, "加密错误")
	// ... 可以继续添加其他预定义错误
)

func ReturnErrDatabase(err error, msg string, notfoundMsg ...string) *AppError {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		if len(notfoundMsg) == 0 {
			notfoundMsg = append(notfoundMsg, ErrNotFound.Message)
		}
		return ErrNotFound.WithMessage(notfoundMsg[0])
	}
	return ErrDatabase.WithMessage(msg)
}

func ReturnErrSimpleDatabase(err error) *AppError {
	return ErrDatabase.WithMessage(err.Error())
}

func ReturnErrInvalidParam(msg string) *AppError {
	return ErrInvalidParam.WithMessage(msg)
}

// ConvertToAppError 将各种错误转换为业务相关的AppError
func ConvertToAppError(err error) *AppError {
	if err == nil {
		return ErrServerBusy.WithMessage("未知错误")
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Code == ErrInvalidParam.Code {
			appErr.Message = TranslateError(errors.New(appErr.Message)).Error()
			return appErr
		}
		return appErr
	}

	// 映射特定的错误到业务错误
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return ErrNotFound.WithMessage("数据不存在")
	case errors.Is(err, gorm.ErrDuplicatedKey):
		return ErrDataExists.WithMessage("数据冲突")
	case errors.Is(err, gorm.ErrInvalidField):
		return ErrDatabase.WithMessage("字段无效")
	case errors.Is(err, gorm.ErrInvalidTransaction):
		return ErrDatabase.WithMessage("数据库事务错误")
	}

	// 处理mysql特定错误
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		return ErrDatabase.WithMessage(mysqlErr.Message)
	}

	return ErrServerBusy.WithMessage(err.Error())
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	TraceID string      `json:"trace_id,omitempty"`
}

func ResponseError(c *gin.Context, err error) {
	appErr := ConvertToAppError(err)
	c.Set("resp-status", appErr.Code)
	c.Set("resp-msg", appErr.Message)
	c.Header("X-Request-Id", c.GetString("trace_id"))
	c.JSON(http.StatusOK, &Response{
		Code:    appErr.Code,
		Message: appErr.Message,
	})
}

func ResponseParamError(c *gin.Context, err error) {
	te := TranslateError(err).Error()
	c.Set("resp-status", ErrInvalidParam.Code)
	c.Set("resp-msg", te)
	c.Header("X-Request-Id", c.GetString("trace_id"))
	if te == "" {
		te = ErrInvalidParam.Message
	}
	c.JSON(http.StatusOK, &Response{
		Code:    ErrInvalidParam.Code,
		Message: te,
	})
}

func ResponseSuccess(c *gin.Context, data interface{}) {
	c.Set("resp-status", http.StatusOK)
	c.Set("resp-msg", "请求成功")
	c.Header("X-Request-Id", c.GetString("trace_id"))
	c.JSON(http.StatusOK, &Response{
		Code:    http.StatusOK,
		Message: "请求成功",
		Data:    data,
	})
}
func ResponseSuccessEncryptData(c *gin.Context, data interface{}, custom func(now int64) (key, nonce string)) {
	c.Set("resp-status", http.StatusOK)
	c.Set("resp-msg", "请求成功")
	response, err := EncryptData(data, custom)
	if err != nil {
		c.JSON(http.StatusOK, &Response{
			Code:    EncryptErr.Code,
			Message: EncryptErr.Message,
		})
		return
	}
	c.JSON(http.StatusOK, &Response{
		Code:    http.StatusOK,
		Message: "请求成功",
		Data:    response,
	})
}

func ResponseThirdPartyHTTPBody(c *gin.Context, body any, code ...int) {
	if len(code) == 0 {
		code = append(code, 200)
	}
	c.JSON(code[0], body)
}

func FuncErr(fun func() error) error {
	if err := fun(); err != nil {
		return ReturnErrSimpleDatabase(err)
	}
	return nil
}
