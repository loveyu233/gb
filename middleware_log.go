package gb

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/textproto"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// LogEntry 表示一个日志条目
type LogEntry struct {
	Level   zerolog.Level
	Message string
	Fields  map[string]any
	Time    time.Time
}

// RequestLogger 存储请求链路中的所有日志
type RequestLogger struct {
	entries []LogEntry
	mu      sync.RWMutex
	ctx     context.Context
	logger  zerolog.Logger
}

// NewRequestLogger 函数用于处理NewRequestLogger相关逻辑。
func NewRequestLogger(ctx context.Context, logger zerolog.Logger) *RequestLogger {
	return &RequestLogger{
		entries: make([]LogEntry, 0),
		ctx:     ctx,
		logger:  logger,
	}
}

// AddEntry 方法用于处理AddEntry相关逻辑。
func (rl *RequestLogger) AddEntry(level zerolog.Level, message string, fields map[string]any) {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	entry := LogEntry{
		Level:   level,
		Message: message,
		Fields:  make(map[string]any),
		Time:    Now(),
	}

	// 复制字段避免并发问题
	for k, v := range fields {
		entry.Fields[k] = v
	}

	rl.entries = append(rl.entries, entry)
}

// Flush 方法用于处理Flush相关逻辑。
func (rl *RequestLogger) Flush() {
	rl.mu.RLock()
	defer rl.mu.RUnlock()

	if len(rl.entries) == 0 {
		return
	}

	// 构建合并的日志事件
	event := rl.logger.Info()

	// 添加所有收集的日志条目
	logEntries := make([]map[string]any, 0, len(rl.entries))
	for _, entry := range rl.entries {
		logEntry := map[string]any{
			"level":     entry.Level.String(),
			"message":   entry.Message,
			"timestamp": entry.Time.Format(time.RFC3339Nano),
		}

		// 添加字段
		if len(entry.Fields) > 0 {
			logEntry["fields"] = entry.Fields
		}

		logEntries = append(logEntries, logEntry)
	}

	// 输出合并的日志
	event.Interface("request_logs", logEntries).
		Int("log_count", len(rl.entries)).
		Msg("Request completed with collected logs")
}

// ContextLogger 提供链路日志记录功能
type ContextLogger struct {
	requestLogger *RequestLogger
}

// Info 方法用于处理Info相关逻辑。
func (cl *ContextLogger) Info() *ContextLogEvent {
	return &ContextLogEvent{
		level:         zerolog.InfoLevel,
		requestLogger: cl.requestLogger,
		fields:        make(map[string]any),
	}
}

// Error 方法用于处理Error相关逻辑。
func (cl *ContextLogger) Error() *ContextLogEvent {
	return &ContextLogEvent{
		level:         zerolog.ErrorLevel,
		requestLogger: cl.requestLogger,
		fields:        make(map[string]any),
	}
}

// Warn 方法用于处理Warn相关逻辑。
func (cl *ContextLogger) Warn() *ContextLogEvent {
	return &ContextLogEvent{
		level:         zerolog.WarnLevel,
		requestLogger: cl.requestLogger,
		fields:        make(map[string]any),
	}
}

// Debug 方法用于处理Debug相关逻辑。
func (cl *ContextLogger) Debug() *ContextLogEvent {
	return &ContextLogEvent{
		level:         zerolog.DebugLevel,
		requestLogger: cl.requestLogger,
		fields:        make(map[string]any),
	}
}

// ContextLogEvent 链路日志事件
type ContextLogEvent struct {
	level         zerolog.Level
	requestLogger *RequestLogger
	fields        map[string]any
}

// Str 方法用于处理Str相关逻辑。
func (e *ContextLogEvent) Str(key, val string) *ContextLogEvent {
	e.fields[key] = val
	return e
}

// Int 方法用于处理Int相关逻辑。
func (e *ContextLogEvent) Int(key string, val int) *ContextLogEvent {
	e.fields[key] = val
	return e
}

// Float64 方法用于处理Float64相关逻辑。
func (e *ContextLogEvent) Float64(key string, val float64) *ContextLogEvent {
	e.fields[key] = val
	return e
}

// Bool 方法用于处理Bool相关逻辑。
func (e *ContextLogEvent) Bool(key string, val bool) *ContextLogEvent {
	e.fields[key] = val
	return e
}

// Err 方法用于处理Err相关逻辑。
func (e *ContextLogEvent) Err(err error) *ContextLogEvent {
	if err != nil {
		e.fields["error"] = err.Error()
	}
	return e
}

// Interface 方法用于处理Interface相关逻辑。
func (e *ContextLogEvent) Interface(key string, val any) *ContextLogEvent {
	e.fields[key] = val
	return e
}

// Dur 方法用于处理Dur相关逻辑。
func (e *ContextLogEvent) Dur(key string, d time.Duration) *ContextLogEvent {
	e.fields[key] = d.String()
	return e
}

// Msg 方法用于处理Msg相关逻辑。
func (e *ContextLogEvent) Msg(msg string) {
	e.requestLogger.AddEntry(e.level, msg, e.fields)
}

// Msgf 方法用于处理Msgf相关逻辑。
func (e *ContextLogEvent) Msgf(format string, v ...any) {
	e.requestLogger.AddEntry(e.level, fmt.Sprintf(format, v...), e.fields)
}

// 上下文键
type contextKey string

const (
	RequestLoggerKey contextKey = "request_logger"
)

const maxLoggedBodyBytes = 1 << 20

// ResponseWriter 是对 gin.ResponseWriter 的包装，用于捕获写入的响应
type ResponseWriter struct {
	gin.ResponseWriter
	body *limitedBuffer
}

type limitedBuffer struct {
	limit     int
	truncated bool
	buf       bytes.Buffer
}

// newLimitedBuffer 函数用于处理newLimitedBuffer相关逻辑。
func newLimitedBuffer(limit int) *limitedBuffer {
	return &limitedBuffer{limit: limit}
}

// Write 方法用于处理Write相关逻辑。
func (b *limitedBuffer) Write(p []byte) (int, error) {
	if b.limit <= 0 {
		return b.buf.Write(p)
	}
	remaining := b.limit - b.buf.Len()
	if remaining > 0 {
		toWrite := remaining
		if len(p) < toWrite {
			toWrite = len(p)
		}
		b.buf.Write(p[:toWrite])
		if len(p) > toWrite {
			b.truncated = true
		}
	} else if len(p) > 0 {
		b.truncated = true
	}
	return len(p), nil
}

// WriteString 方法用于处理WriteString相关逻辑。
func (b *limitedBuffer) WriteString(s string) (int, error) {
	return b.Write([]byte(s))
}

// Bytes 方法用于处理Bytes相关逻辑。
func (b *limitedBuffer) Bytes() []byte {
	return b.buf.Bytes()
}

// Len 方法用于处理Len相关逻辑。
func (b *limitedBuffer) Len() int {
	return b.buf.Len()
}

// Truncated 方法用于处理Truncated相关逻辑。
func (b *limitedBuffer) Truncated() bool {
	return b.truncated
}

// Write 方法用于处理Write相关逻辑。
func (w ResponseWriter) Write(b []byte) (int, error) {
	// 写入到缓冲区
	w.body.Write(b)
	// 继续原始的写入操作
	return w.ResponseWriter.Write(b)
}

// WriteString 方法用于处理WriteString相关逻辑。
func (w ResponseWriter) WriteString(s string) (int, error) {
	// 写入到缓冲区
	w.body.WriteString(s)
	// 继续原始的写入操作
	return w.ResponseWriter.WriteString(s)
}

var zlog zerolog.Logger

// init 函数用于处理init相关逻辑。
func init() {
	//zerolog.TimeFieldFormat = CSTLayout
	zerolog.TimestampFunc = func() time.Time {
		return Now()
	}
	zlog = zerolog.New(os.Stdout).With().
		Timestamp().
		Logger()
}

type ReqLog struct {
	ReqTime     time.Time         `json:"req_time"`
	Module      string            `json:"module,omitempty"`
	Option      string            `json:"option,omitempty"`
	Method      string            `json:"method,omitempty"`
	Path        string            `json:"path,omitempty"`
	URL         string            `json:"url,omitempty"`
	IP          string            `json:"ip,omitempty"`
	Content     map[string]any    `json:"content,omitempty"`
	Headers     map[string]string `json:"headers,omitempty"`
	Params      map[string]any    `json:"params,omitempty"`
	Status      int               `json:"status,omitempty"`
	Latency     time.Duration     `json:"latency,omitempty"`
	Body        map[string]any    `json:"body,omitempty"`
	RespStatus  int               `json:"resp_status"`  // 响应数据中的状态码
	RespMessage string            `json:"resp_message"` // 响应数据中的message
}

type MiddlewareLogConfig struct {
	HeaderKeys  []string
	ContentKeys []string
	SaveLog     func(ReqLog)
}

type FileInfo struct {
	Filename string               `json:"filename"`
	Size     int64                `json:"size"`
	Header   textproto.MIMEHeader `json:"header"`
}

// MiddlewareLogger 函数用于处理MiddlewareLogger相关逻辑。
func MiddlewareLogger(mc MiddlewareLogConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		startTime := Now()
		// 创建请求日志器
		requestLogger := NewRequestLogger(c.Request.Context(), zlog)
		c.Set(string(RequestLoggerKey), requestLogger)

		// 创建自定义 ResponseWriter
		bodyBuffer := newLimitedBuffer(maxLoggedBodyBytes)
		responseWriter := &ResponseWriter{
			ResponseWriter: c.Writer,
			body:           bodyBuffer,
		}
		c.Writer = responseWriter

		// 获取请求参数，分类存储
		params := make(map[string]any)

		// 1. 处理URL查询参数 (query parameters)
		queryParams := make(map[string]any)
		for k, v := range c.Request.URL.Query() {
			if len(v) == 1 {
				queryParams[k] = v[0]
			} else {
				queryParams[k] = v
			}
		}
		if len(queryParams) > 0 {
			params["query"] = queryParams
		}

		// 2. 处理路径参数 (path parameters)
		pathParams := make(map[string]any)
		for _, param := range c.Params {
			pathParams[param.Key] = param.Value
		}
		if len(pathParams) > 0 {
			params["path"] = pathParams
		}

		// 3. 获取请求体和处理不同类型的参数
		contentType := c.ContentType()

		if strings.Contains(contentType, "multipart/form-data") {
			// 处理 multipart/form-data（包含文件上传）
			err := c.Request.ParseMultipartForm(32 << 20) // 32MB 最大内存
			if err == nil && c.Request.MultipartForm != nil {
				// 处理普通表单字段
				formData := make(map[string]any)
				for key, values := range c.Request.MultipartForm.Value {
					if len(values) == 1 {
						formData[key] = values[0]
					} else {
						formData[key] = values
					}
				}
				if len(formData) > 0 {
					params["form"] = formData
				}

				// 处理文件字段
				fileParams := make(map[string][]map[string]FileInfo)
				for key, files := range c.Request.MultipartForm.File {
					fileInfos := make([]map[string]FileInfo, len(files))
					for i, file := range files {
						fileInfos[i] = map[string]FileInfo{
							key: {
								Filename: file.Filename,
								Size:     file.Size,
								Header:   file.Header,
							},
						}
					}
					fileParams[key] = fileInfos
				}
				if len(fileParams) > 0 {
					params["files"] = fileParams
				}
			}
		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") {
			// 处理表单编码数据
			err := c.Request.ParseForm()
			if err == nil {
				formData := make(map[string]any)
				for key, values := range c.Request.PostForm {
					if len(values) == 1 {
						formData[key] = values[0]
					} else {
						formData[key] = values
					}
				}
				if len(formData) > 0 {
					params["form"] = formData
				}
			}
		} else if strings.Contains(contentType, "application/json") {
			if ok, reason := shouldCaptureRequestBody(c.Request); ok {
				if requestBody, err := io.ReadAll(c.Request.Body); err == nil {
					c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
					if len(requestBody) > 0 {
						var bodyParams any
						if err := json.Unmarshal(requestBody, &bodyParams); err == nil {
							params["json"] = bodyParams
						} else {
							params["json_raw"] = string(requestBody)
						}
					}
				} else {
					recordBodySkip(params, fmt.Sprintf("failed to read request body: %v", err))
				}
			} else {
				recordBodySkip(params, reason)
			}
		} else if strings.Contains(contentType, "application/xml") || strings.Contains(contentType, "text/xml") {
			if ok, reason := shouldCaptureRequestBody(c.Request); ok {
				if requestBody, err := io.ReadAll(c.Request.Body); err == nil {
					c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
					if len(requestBody) > 0 {
						params["xml"] = string(requestBody)
					}
				} else {
					recordBodySkip(params, fmt.Sprintf("failed to read request body: %v", err))
				}
			} else {
				recordBodySkip(params, reason)
			}
		} else {
			if ok, reason := shouldCaptureRequestBody(c.Request); ok {
				if requestBody, err := io.ReadAll(c.Request.Body); err == nil {
					c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
					if len(requestBody) > 0 {
						params["raw"] = map[string]any{
							"content_type": contentType,
							"body":         string(requestBody),
							"size":         len(requestBody),
						}
					}
				} else {
					recordBodySkip(params, fmt.Sprintf("failed to read request body: %v", err))
				}
			} else {
				recordBodySkip(params, reason)
			}
		}

		headerMap := make(map[string]string)
		for _, item := range mc.HeaderKeys {
			headerMap[item] = c.GetHeader(item)
		}

		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		fullURL := scheme + "://" + c.Request.Host + c.Request.RequestURI

		c.Next()
		if c.GetBool("skip") {
			c.Next()
			return
		}

		var contentKV = make(map[string]any)
		for _, key := range mc.ContentKeys {
			value, exists := c.Get(key)
			if exists {
				contentKV[key] = value
			}
		}
		// 记录请求开始信息
		requestLogger.AddEntry(zerolog.InfoLevel, "request", map[string]any{
			"req_time":   startTime.Format(CSTLayout),
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"full_url":   fullURL,
			"req_body":   params,
			"user_agent": c.Request.UserAgent(),
			"client_ip":  c.ClientIP(),
			"header":     headerMap,
			"module":     c.GetString("module"),
			"option":     c.GetString("option"),
			"content_kv": contentKV,
		})

		if c.GetBool("only-req") {
			// 输出所有收集的日志
			requestLogger.Flush()
			c.Next()
			return
		}

		duration := time.Since(startTime)

		bodyMap := make(map[string]any)
		if !c.GetBool("brief") {
			if data := bodyBuffer.Bytes(); len(data) > 0 {
				if err := json.Unmarshal(data, &bodyMap); err != nil {
					bodyMap["raw"] = string(data)
				}
			}
			if bodyBuffer.Truncated() {
				bodyMap["resp_truncated"] = fmt.Sprintf("response body exceeded %d bytes and was truncated", maxLoggedBodyBytes)
			}
		}
		bodyMap["resp-status"] = c.GetInt("resp-status")
		bodyMap["resp-message"] = c.GetString("resp-msg")

		requestLogger.AddEntry(zerolog.InfoLevel, "response", map[string]any{
			"status_code": c.Writer.Status(),
			"duration":    duration.String(),
			"resp_body":   bodyMap,
		})

		// 输出所有收集的日志
		requestLogger.Flush()

		if mc.SaveLog != nil && !c.GetBool("no_record") {
			mc.SaveLog(ReqLog{
				ReqTime:     startTime,
				Module:      c.GetString("module"),
				Option:      c.GetString("option"),
				Method:      c.Request.Method,
				Path:        c.Request.URL.Path,
				URL:         fullURL,
				IP:          c.ClientIP(),
				Content:     contentKV,
				Headers:     headerMap,
				Params:      params,
				Status:      c.Writer.Status(),
				Latency:     duration,
				Body:        bodyMap,
				RespStatus:  c.GetInt("resp-status"),
				RespMessage: c.GetString("resp-msg"),
			})
		}
	}
}

// GetContextLogger 函数用于处理GetContextLogger相关逻辑。
func GetContextLogger(c *gin.Context) *ContextLogger {
	if requestLogger, exists := c.Get(string(RequestLoggerKey)); exists {
		if rl, ok := requestLogger.(*RequestLogger); ok {
			return &ContextLogger{requestLogger: rl}
		}
	}

	// 如果获取失败，返回一个空的日志器避免 panic
	return &ContextLogger{
		requestLogger: NewRequestLogger(context.Background(), log.Logger),
	}
}

// WriteGinInfoLog 函数用于处理WriteGinInfoLog相关逻辑。
func WriteGinInfoLog(c *gin.Context, format string, args ...any) {
	GetContextLogger(c).Info().Msgf(format, args)
}

// WriteGinDebugLog 函数用于处理WriteGinDebugLog相关逻辑。
func WriteGinDebugLog(c *gin.Context, format string, args ...any) {
	GetContextLogger(c).Debug().Msgf(format, args)
}

// WriteGinWarnLog 函数用于处理WriteGinWarnLog相关逻辑。
func WriteGinWarnLog(c *gin.Context, format string, args ...any) {
	GetContextLogger(c).Warn().Msgf(format, args)
}

// WriteGinErrLog 函数用于处理WriteGinErrLog相关逻辑。
func WriteGinErrLog(c *gin.Context, format string, args ...any) {
	GetContextLogger(c).Error().Msgf(format, args)
}

// GinLogSetModuleName 函数用于处理GinLogSetModuleName相关逻辑。
func GinLogSetModuleName(name string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("module", name)
		c.Next()
	}
}

// GinLogSetOptionName 函数用于处理GinLogSetOptionName相关逻辑。
func GinLogSetOptionName(name string, noRecord ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("option", name)
		if len(noRecord) > 0 && noRecord[0] {
			c.Set("no_record", true)
		}
		c.Next()
	}
}

// GinLogSetSkipLogFlag 函数用于处理GinLogSetSkipLogFlag相关逻辑。
func GinLogSetSkipLogFlag() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("skip", true)
		c.Next()
	}
}

// GinLogOnlyReqMsg 函数用于处理GinLogOnlyReqMsg相关逻辑。
func GinLogOnlyReqMsg() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("only-req", true)
		c.Next()
	}
}

// GinLogBriefInformation 函数用于处理GinLogBriefInformation相关逻辑。
func GinLogBriefInformation() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("brief", true)
		c.Next()
	}
}

// shouldCaptureRequestBody 函数用于处理shouldCaptureRequestBody相关逻辑。
func shouldCaptureRequestBody(r *http.Request) (bool, string) {
	if r == nil {
		return false, "request is nil"
	}
	if r.Body == nil || r.Body == http.NoBody {
		return false, "request body is empty"
	}
	if r.ContentLength == -1 {
		return false, "request body size is unknown (chunked transfer)"
	}
	if r.ContentLength > maxLoggedBodyBytes {
		return false, fmt.Sprintf("request body exceeds %d bytes", maxLoggedBodyBytes)
	}
	return true, ""
}

// recordBodySkip 函数用于处理recordBodySkip相关逻辑。
func recordBodySkip(params map[string]any, reason string) {
	if reason == "" {
		return
	}
	if _, exists := params["body_skipped"]; !exists {
		params["body_skipped"] = reason
	}
}
