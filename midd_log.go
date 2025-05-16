package gb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"time"
)

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

// ResponseLogger 中间件用于记录响应数据
func ResponseLogger(saveLog ...func(ReqLog)) gin.HandlerFunc {
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

		// 获取请求体
		var requestBody []byte
		if c.Request.Body != nil && c.Request.Body != http.NoBody {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体以便后续处理
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

			// 尝试解析 JSON 请求体
			if c.ContentType() == "application/json" && len(requestBody) > 0 {
				bodyParams := make(map[string]interface{})
				if err := json.Unmarshal(requestBody, &bodyParams); err == nil {
					for k, v := range bodyParams {
						params[k] = v
					}
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

		// 获取 Authorization 头
		authorization := c.GetHeader("YXB-TOKEN")

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
		fmt.Printf("Authorization: %s\n", authorization)
		fmt.Printf("参数: %v\n", params)
		fmt.Printf("状态码: %d\n", c.Writer.Status())
		fmt.Printf("耗时: %v\n", latency)
		fmt.Printf("响应数据: %s\n", responseBody)
		fmt.Printf("========================\n")

		if c.GetBool("record") && len(saveLog) > 0 {
			saveLog[0](ReqLog{
				ReqTime: startTime,
				User:    value,
				Module:  c.GetString("module"),
				Option:  c.GetString("option"),
				Method:  method,
				URL:     fullURL,
				IP:      clientIP,
				Token:   authorization,
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
	Token   string
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
