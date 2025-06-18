# GB 后端开发框架

GB 是一个基于 Gin 框架的 Go 语言后端开发框架，提供了丰富的功能和工具，帮助开发者快速构建高性能的 Web 应用。

## 主要特性

- 🚀 基于 Gin 框架的 HTTP 服务器
- 🔐 JWT 认证支持（支持 Redis 黑名单和 Token 校验）
- 📝 自动生成 Swagger API 文档（支持 OpenAPI 2.0）
- 🎯 优雅的路由管理（支持公共路由和私有路由）
- ⏱️ 请求追踪和日志记录
- 🔄 优雅关闭支持
- 🗄️ Redis 客户端支持（支持单节点、哨兵模式、Cluster集群模式）
- 📊 GORM 数据库集成
- ⚡ 高性能中间件
- 🔍 参数验证（支持中文错误提示）
- 🕒 时间处理工具（支持中国时区）
- 💰 Decimal 处理工具
- 🚨 统一的错误处理机制
- 📝 结构化的日志记录

## 快速开始

### 安装

```bash
go get github.com/loveyu233/gb
```

### 基本使用

```go
package main

import "github.com/loveyu233/gb"

func main() {
	gb.InitHTTPServerAndStart(":8080", gb.WithGinRouterModel("release"))
}
```

## 核心功能

### 错误处理

统一的错误处理机制，支持业务错误码和错误信息：

```go
// 预定义错误
var (
    ErrBadRequest   = gb.NewAppError(400000, "请求错误")
    ErrInvalidParam = gb.NewAppError(400001, "请求参数错误")
    ErrUnauthorized = gb.NewAppError(401000, "用户未登录或token已失效")
    ErrForbidden    = gb.NewAppError(403000, "权限不足")
    ErrNotFound     = gb.NewAppError(404000, "数据不存在")
    ErrServerBusy   = gb.NewAppError(500000, "服务器繁忙")
)

// 统一响应格式
type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
    TraceID string      `json:"trace_id,omitempty"`
}

// 使用示例
func handler(c *gin.Context) {
    if err := someOperation(); err != nil {
        gb.ResponseError(c, err)
        return
    }
    gb.ResponseSuccess(c, data)
}
```

### 路由管理

支持公共路由和私有路由的注册，并提供模块化路由管理：

```go
func registerPublicRoutes(r *gin.RouterGroup) {
    routes := r.Group("/api", gb.SetModuleName("API模块"))
    {
        routes.GET("/hello", gb.SetOptionName("hello"), func(c *gin.Context) {
            gb.ResponseSuccess(c, "hello")
        })
    }
}
```

### JWT 认证

支持多种认证方式和配置选项：

```go
// 创建 JWT 服务
tokenService := gb.NewJWTTokenService("your-secret-key",
    gb.WithRedisClient(redisClient),
    gb.WithRedisBlacklist(true),
    gb.WithRedisTokenCheck(true),
)

// 生成 Token
token, err := tokenService.Generate(user, 24*time.Hour)
```

### Redis 客户端

支持多种 Redis 部署模式：

```go
// 创建 Redis 客户端
client, err := gb.NewRedisClient("localhost:6379",
    gb.WithDB(0),
    gb.WithPassword("password"),
    gb.WithDialTimeout(5*time.Second),
)

// 使用示例
client.Set("key", "value", time.Hour)
client.Get("key")
client.HGetAll("hash")
client.ZAdd("sorted_set", score, member)
```

### 数据库集成

基于 GORM 的数据库操作：

```go
// 初始化数据库连接
db, err := gb.InitGormDB("dsn", gb.GormDefaultLogger())
```

### 参数验证

支持结构体验证和中文错误提示：

```go
type User struct {
    Name     string `json:"name" validate:"required"`
    Age      int    `json:"age" validate:"required,min=18"`
    Phone    string `json:"phone" validate:"required,phone"`
}

// 验证错误会自动转换为中文提示
```

### 时间处理

内置中国时区支持：

```go
// 获取当前时间
now := gb.GetCurrentTime()

// 格式化时间
formatted := gb.DateTimeToString(now)

// 解析时间字符串
t, err := gb.StringToDateTime("2024-01-01 12:00:00")

// 生成目录名
dirName := gb.MakeDirNameByCurrentTime() // 返回格式：2024/0101
```

### Decimal 处理

提供精确的数值计算工具：

```go
// 浮点数转 Decimal
decimal := gb.Float64ToDecimal(3.14)

// Decimal 转浮点数
float := gb.DecimalToFloat64(decimal)

// 元转分
fen := gb.DecimalYuanToFen(decimal)

// 分转元
yuan := gb.FenToDecimalYuan(100)

// 计算百分比
result := gb.DecimalPercent(value, percent)
```

### Swagger 文档生成

支持 OpenAPI 2.0 规范的 API 文档生成：

```go
// 创建 Swagger 生成器
generator := gb.NewSwaggerGenerator(gb.SwaggerGlobalConfig{
    Title:       "API文档标题",
    Description: "API文档描述",
    Version:     "1.0.0",
    Host:        "localhost:8080",
    BasePath:    "/api/v1",
    OutputPath:  "./swagger.json",
})

// 添加 API 信息
generator.AddAPI(gb.SwaggerAPIInfo{
    Path:        "/auth/login",
    Method:      "POST",
    Summary:     "用户登录",
    Description: "用户登录接口",
    Tags:        []string{"认证"},
    Request:     LoginRequest{},
    Response:    LoginResponse{},
    PathParams: []gb.SwaggerParamDescription{
        {
            Name:        "id",
            Description: "用户ID",
            Type:        "integer",
        },
    },
    QueryParams: []gb.SwaggerParamDescription{
        {
            Name:        "type",
            Description: "登录类型",
            Type:        "string",
        },
    },
    ResponseStatus: map[string]string{
        "200": "成功",
        "400": "请求参数错误",
        "401": "未授权",
    },
})

// 生成文档
err := generator.Generate()
```

### 链式调用方式

```go
// 使用链式调用方式添加 API
generator.WithOperation("/users/{id}", "GET").
    Summary("获取用户信息").
    Description("根据用户ID获取用户详细信息").
    Tags([]string{"用户"}).
    PathParam("id", "用户ID").
    QueryParam("fields", "需要返回的字段", false).
    Response("200", "成功", UserResponse{}).
    Build()
```

## 中间件

框架提供了多个实用的中间件：

- `AddTraceID`: 请求追踪
- `AddRequestTime`: 请求时间记录
- `ResponseLogger`: 响应日志记录
- `GinRecovery`: 错误恢复
- `GinAuth`: 认证中间件

## 日志记录

支持结构化的日志记录：

```go
// 请求日志记录
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
```

## 依赖

- Gin v1.10.0
- GORM v1.26.1
- Redis v9.8.0
- JWT v5.2.2
- 其他依赖见 go.mod

## 贡献

欢迎提交 Issue 和 Pull Request！