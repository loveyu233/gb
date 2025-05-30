package gb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"runtime/debug"
	"strings"
	"time"
)

func AddTraceID() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.Request.Header.Get("trace_id")
		if traceID == "" {
			traceID = uuid.NewString()
		}
		c.Set("trace_id", traceID)
	}
}

// AddRequestTime 注入请求时间
func AddRequestTime() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("request_time", GetCurrentTime())
	}
}

// GinRecovery recover掉项目可能出现的panic，并使用zap记录相关日志
func GinRecovery(stack bool, log ...GBLog) gin.HandlerFunc {
	var loclLog GBLog

	if len(log) > 0 {
		loclLog = log[0]
	} else {
		loclLog = new(GbDefaultlogger)
	}
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						if strings.Contains(strings.ToLower(se.Error()), "broken pipe") || strings.Contains(strings.ToLower(se.Error()), "connection reset by peer") {
							brokenPipe = true
						}
					}
				}

				httpRequest, _ := httputil.DumpRequest(c.Request, false)
				if brokenPipe {
					loclLog.Errorf("path: %s\nerror: %v\nreq:%s\n", c.Request.URL.Path, err, string(httpRequest))
					c.Error(err.(error)) // nolint: errcheck
					c.Abort()
					return
				}

				if stack {
					loclLog.Errorf("[Recovery from panic]\nerr:%s\nreq:%s\nstack:%s\n", err, string(httpRequest), string(debug.Stack()))
				} else {
					loclLog.Errorf("[Recovery from panic]\nerr:%s\nreq:%s\n", err, string(httpRequest))

				}
				ResponseError(c, ErrServerBusy.WithMessage(fmt.Sprintf("[Recovery from panic]: %v", err)))
				c.AbortWithStatus(http.StatusOK)
			}
		}()
		c.Next()
	}
}

// ResponseWriter 是对 gin.ResponseWriter 的包装，用于捕获写入的响应
type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写 Write 方法以捕获响应内容
func (w ResponseWriter) Write(b []byte) (int, error) {
	// 写入到缓冲区
	w.body.Write(b)
	// 继续原始的写入操作
	return w.ResponseWriter.Write(b)
}

// WriteString 重写 WriteString 方法以捕获响应内容
func (w ResponseWriter) WriteString(s string) (int, error) {
	// 写入到缓冲区
	w.body.WriteString(s)
	// 继续原始的写入操作
	return w.ResponseWriter.WriteString(s)
}

type MiddlewareLogConfig struct {
	HeaderKeys []string
	SaveLog    func(ReqLog)
	IsSaveLog  bool
}

// ResponseLogger 中间件用于记录响应数据
func ResponseLogger(config MiddlewareLogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := time.Now()

		// 创建自定义 ResponseWriter
		bodyBuffer := &bytes.Buffer{}
		responseWriter := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:           bodyBuffer,
		}
		c.Writer = responseWriter

		// 获取请求参数
		params := make(map[string]interface{})

		// 处理URL查询参数
		for k, v := range c.Request.URL.Query() {
			if len(v) == 1 {
				params[k] = v[0]
			} else {
				params[k] = v
			}
		}

		// 处理路径参数
		for _, param := range c.Params {
			params[param.Key] = param.Value
		}

		// 获取请求体和处理不同类型的参数
		contentType := c.ContentType()

		if strings.Contains(contentType, "multipart/form-data") {
			// 处理 multipart/form-data（包含文件上传）
			err := c.Request.ParseMultipartForm(32 << 20) // 32MB 最大内存
			if err == nil && c.Request.MultipartForm != nil {
				// 处理普通表单字段
				for key, values := range c.Request.MultipartForm.Value {
					if len(values) == 1 {
						params[key] = values[0]
					} else {
						params[key] = values
					}
				}

				// 处理文件字段
				for key, files := range c.Request.MultipartForm.File {
					if len(files) == 1 {
						params[key] = fmt.Sprintf("[文件: %s]", files[0].Filename)
					} else {
						fileNames := make([]string, len(files))
						for i, file := range files {
							fileNames[i] = fmt.Sprintf("[文件: %s]", file.Filename)
						}
						params[key] = fileNames
					}
				}
			}
		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			// 处理表单编码数据
			err := c.Request.ParseForm()
			if err == nil {
				for key, values := range c.Request.PostForm {
					if len(values) == 1 {
						params[key] = values[0]
					} else {
						params[key] = values
					}
				}
			}
		} else if strings.Contains(contentType, "application/json") {
			// 处理 JSON 数据
			var requestBody []byte
			if c.Request.Body != nil && c.Request.Body != http.NoBody {
				requestBody, _ = io.ReadAll(c.Request.Body)
				// 重新设置请求体以便后续处理
				c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

				if len(requestBody) > 0 {
					bodyParams := make(map[string]interface{})
					if err := json.Unmarshal(requestBody, &bodyParams); err == nil {
						for k, v := range bodyParams {
							params[k] = v
						}
					}
				}
			}
		} else {
			// 其他类型，尝试读取原始请求体
			var requestBody []byte
			if c.Request.Body != nil && c.Request.Body != http.NoBody {
				requestBody, _ = io.ReadAll(c.Request.Body)
				// 重新设置请求体以便后续处理
				c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

				if len(requestBody) > 0 {
					params["raw_body"] = string(requestBody)
				}
			}
		}

		// 处理请求
		c.Next()

		// 计算耗时
		latency := time.Since(startTime)

		// 获取请求信息
		method := c.Request.Method
		// 完整URL（包含协议、域名和路径）
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		fullURL := scheme + "://" + c.Request.Host + c.Request.RequestURI

		// 获取客户端 IP
		clientIP := c.ClientIP()

		headerMap := make(map[string]string)
		for _, item := range config.HeaderKeys {
			headerMap[item] = c.GetHeader(item)
		}
		headerBytes, _ := json.Marshal(headerMap)

		// 获取响应数据
		responseBody := bodyBuffer.String()

		// 记录日志 - 使用结构化格式便于阅读
		value, _ := c.Get("user")
		fmt.Printf("===== 请求响应日志 =====\n")
		fmt.Printf("请求时间: %s\n", startTime.Format("2006-01-02 - 15:04:05.000"))
		fmt.Printf("user: %+v\n", value)
		fmt.Printf("模块: %s\n", c.GetString("module"))
		fmt.Printf("操作: %s\n", c.GetString("option"))
		fmt.Printf("方法: %s\n", method)
		fmt.Printf("URL: %s\n", fullURL)
		fmt.Printf("IP: %s\n", clientIP)
		fmt.Printf("请求头: %s\n", string(headerBytes))
		fmt.Printf("参数: %v\n", params)
		fmt.Printf("状态码: %d\n", c.Writer.Status())
		fmt.Printf("耗时: %v\n", latency)
		fmt.Printf("响应数据: %s\n", responseBody)
		fmt.Printf("========================\n")

		if c.GetBool("record") && config.IsSaveLog {
			config.SaveLog(ReqLog{
				ReqTime: startTime,
				User:    value,
				Module:  c.GetString("module"),
				Option:  c.GetString("option"),
				Method:  method,
				URL:     fullURL,
				IP:      clientIP,
				Headers: string(headerBytes),
				Params:  params,
				Status:  c.Writer.Status(),
				Latency: latency,
				Body:    responseBody,
			})
		}
	}
}

type ReqLog struct {
	ReqTime time.Time
	User    any
	Module  string
	Option  string
	Method  string
	URL     string
	IP      string
	Headers string
	Params  map[string]any
	Status  int
	Latency time.Duration
	Body    string
}

func SetModuleName(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("module", name)
		c.Next()
	}
}

func SetOptionName(name string, noRecord ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("option", name)
		if len(noRecord) > 0 && noRecord[0] {
			c.Set("no_record", true)
		}
		c.Next()
	}
}
