# GB - Go 后端开发框架

GB 是一个基于 Gin 框架的 Go 语言后端开发框架，提供了丰富的功能和工具，帮助开发者快速构建高性能的 Web 应用。

## 主要特性

- 🚀 基于 Gin 框架的 HTTP 服务器
- 🔐 JWT 认证支持（支持 Redis 黑名单和 Token 校验）
- 📝 自动生成 Swagger API 文档
- 🎯 优雅的路由管理（支持公共路由和私有路由）
- ⏱️ 请求追踪和日志记录
- 🔄 优雅关闭支持
- 🗄️ Redis 客户端支持（支持单节点、哨兵模式、Cluster集群模式）
- 📊 GORM 数据库集成
- ⚡ 高性能中间件
- 🔍 参数验证（支持中文错误提示）
- 🕒 时间处理工具（支持中国时区）
- 💰 Decimal 处理工具
- 🌐 sql日期,时间,数组类型

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
	// 初始化路由
	gb.InitRouter("debug", "/api",
		[]gin.HandlerFunc{gb.GinAuth(&User{}, defaultGinConfig)},
		gb.AddTraceID(),
		gb.AddRequestTime(),
		gb.ResponseLogger(),
		gb.GinRecovery(true),
	)

	// 创建并启动服务器
	server := gb.CreateHTTPServer(":8080")
	go gb.StartHTTPServer(server)
	gb.SetupGracefulShutdown(server)
}
```

## 核心功能

### Swagger 文档生成

自动生成 API 文档：

```go
config := gb.SwaggerGlobalConfig{
	Title:       "API文档标题",
	Description: "API文档描述",
	Version:     "1.0.0",
	Host:        "localhost:8080",
	BasePath:    "/api/v1",
	OutputPath:  "./swagger.json",
}

generator := gb.NewSwaggerGenerator(config)
generator.AddAPI(gb.SwaggerAPIInfo{
	Path:        "/auth/login",
	Method:      "POST",
	Summary:     "用户登录",
	Description: "用户登录接口",
	Tags:        []string{"认证"},
})

// 在指定位置[OutputPath]生成swagger.json文件
generator.Generate()
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
```

## 中间件

框架提供了多个实用的中间件：

- `AddTraceID`: 请求追踪
- `AddRequestTime`: 请求时间记录
- `ResponseLogger`: 响应日志记录
- `GinRecovery`: 错误恢复
- `GinAuth`: 认证中间件

## 工具函数

- 参数验证和转换
- 时间处理（支持中国时区）
- 错误处理
- Decimal 处理
- Redis 操作
- 数据库操作

## 依赖

- Gin v1.10.0
- GORM v1.26.1
- Redis v9.8.0
- JWT v5.2.2
- 其他依赖见 go.mod

## 贡献

欢迎提交 Issue 和 Pull Request！
