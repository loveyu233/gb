# gb · Go Web 多端服务工具箱

`gb` 是一个面向业务团队的 Go Web 工具箱，提供 **统一的登录 / 支付多渠道聚合、Web 中间件、数据访问、任务调度、日志与常见工具**。通过约定式初始化 (`InitXxx`) 与单例调用 (`InsXxx`) 的方式，开发者可以在最短时间内搭建具备生产特性的 API 服务。

## 目录
- [核心特性](#核心特性)
- [快速开始](#快速开始)
- [多渠道登录 / 支付服务](#多渠道登录--支付服务)
- [Web & 中间件](#web--中间件)
- [数据与工具](#数据与工具)
- [模块一览](#模块一览)
- [开发建议](#开发建议)
- [贡献指南](#贡献指南)

## 核心特性
| 维度 | 能力 |
| --- | --- |
| 多渠道登录 / 支付 | 内置微信 / 支付宝 Provider，可拓展自定义渠道，统一路由与日志沉淀 |
| Web 中台 | Gin 引擎封装、TraceID、链路日志、CORS、恢复、Swagger、参数校验 |
| 数据访问 | GORM 包装、参数 Scope 工具、Redis/RedSync、分布式锁、Snowflake ID |
| 工程效率 | Zerolog 日志、Resty HTTP 客户端、Gocron 调度、Excel 导入导出、Lua 扩展 |
| 安全工具 | JWT 中间件、密码哈希、掩码、AES 加密、请求重试及限流脚本 |

## 快速开始
```go
package main

import (
    "log"

    "github.com/loveyu233/gb"
    "github.com/loveyu233/gb/channel"
    "github.com/loveyu233/gb/login"
    "github.com/loveyu233/gb/pay"
)

func main() {
    // 初始化数据库、Redis、HTTP 服务等
    if err := gb.InitGormDB(gb.GormConnConfig{Username: "root", Password: "123456", Host: "127.0.0.1", Port: 3306, Database: "demo"}, gb.GormDefaultLogger()); err != nil {
        log.Fatal(err)
    }
    if err := gb.InitRedis(gb.WithRedisAddressOption([]string{"127.0.0.1:6379"})); err != nil {
        log.Fatal(err)
    }

    // 构建渠道 Provider（以微信/支付宝示例，业务方实现对应接口即可）
    wxLoginSvc, _ := login.InitWXMiniProgramService(/* 配置 */)
    wxPaySvc, _ := pay.InitWXWXPaymentApp(/* 配置 */)
    pay.InitAliClient(/* 配置 */, true, &CustomZFBImpl{})

    unifiedSvc, _ := channel.NewService(
        channel.Provider{Login: wxLoginSvc, Payment: wxPaySvc},
        channel.Provider{Login: pay.InsZFB, Payment: pay.InsZFB},
    )

    router := gb.InitHTTPServerAndStart(":8080",
        gb.WithGinRouterPrefix("/api"),
        gb.WithGinRouterGlobalMiddleware(gb.GinLogSetModuleName("demo")),
    )
    unifiedSvc.RegisterRoutes(router.Group("/channel"))

    select {}
}
```
路由示例：
- `POST /api/channel/wechat/login`
- `POST /api/channel/alipay/pay`
- `POST /api/channel/wechat/pay/notify`
- `POST /api/channel/alipay/refund`

## 多渠道登录 / 支付服务
`channel.Service` 通过注册 `LoginProvider` 与 `PaymentProvider`，实现：
- **统一路由**：自动根据 `:channel` 参数选择实现。
- **日志与追踪**：中间件自动写入 `module/option`，可结合 `MiddlewareLogger` 持久化。
- **业务解耦**：微信 / 支付宝 / 自定义渠道共享同一套 Handler，业务方只需实现接口。

### 接口约定
```go
type LoginProvider interface {
    ChannelName() string
    HandleLogin(*gin.Context)
    SaveLoginLog() bool
}

type PaymentProvider interface {
    ChannelName() string
    HandlePay(*gin.Context)
    HandlePayNotify(*gin.Context)
    HandleRefund(*gin.Context)
    HandleRefundNotify(*gin.Context)
    SavePaymentLog() bool
}
```
官方实现：
- 微信：`WXMini`（登录）、`WXPay`（支付）
- 支付宝：`ZFBClient`（登录 + 支付）

## Web & 中间件
- `InitHTTPServerAndStart`：统一封装 Gin Engine、路由前缀、读写超时等配置。
- 中间件：TraceID、请求日志（支持敏感头屏蔽、Body 截断）、CORS、异常恢复、请求耗时统计。
- 参数工具：`params_verfiy`、`params_time`、`gin_param` 等结构体 Scope，可直接挂载到 GORM 查询。
- Swagger：`swagger.Generator` 支持结构体、全局参数、模型自动生成。

## 数据与工具
- **GORM**：`InitGormDB` / `InsDB` 统一连接；常用 Scope（分页、时间、关键字）内置安全校验。
- **Redis**：`InitRedis`、RedSync 分布式锁、限流 / 计数 / 队列 / Bloom / HLL 等 Lua 脚本封装。
- **工具集**：Resty HTTP 客户端、Gocron 定时任务、ants 协程池、Excel 导入导出、AES/密码/掩码、Diff 比较、模板替换、日志适配等。

## 模块一览
| 模块 | 说明 |
| --- | --- |
| `channel` | 多渠道登录 / 支付统一调度与路由 |
| `login` | 微信小程序登录（可拓展自定义实现） |
| `pay` | 微信 / 支付宝支付能力、通知、退款 |
| `middleware_*` | 日志、TraceID、CORS、恢复、耗时统计等中间件 |
| `redis.go` | Redis 客户端、Lua 脚本（锁、限流、计数、布隆等） |
| `excel_*` | Excel 导入/导出、Mapper、格式化工具 |
| `sql_type.go` | 自定义日期/时间类型、JSON Slice、防止注入的 Scope |
| `swagger.go` | Swagger 文档生成器与全局参数管理 |
| 其他包 | `auth_jwt`、`params_*`、`mask`、`encrypt`、`snowflake` 等通用能力 |

## 开发建议
1. **模块初始化一次即可**：例如 `InitGormDB`、`InitRedis`、`InitWXWXPaymentApp`，使用单例 `InsXxx`，避免重复连接。
2. **统一错误响应**：`ResponseSuccess` / `ResponseError` / `ConvertToAppError` 帮助保持错误码一致。
3. **安全使用 Scope**：`common.ScopeOrderDesc`、`ScopeFilterKeyword` 已加列名校验，鼓励二次封装业务查询。
4. **日志与审计**：`MiddlewareLogger` 支持敏感头掩码与 Body 截断；支付、登录通过 `Save*Log` 控制审计需求。
5. **扩展渠道**：实现 `LoginProvider` / `PaymentProvider`，即可接入任意第三方；建议在 Provider 内聚合底层 SDK。

## 贡献指南
- **Issue**：欢迎提交需求、Bug 或安全问题。
- **Pull Request**：请遵循 Go 官方格式 (`gofmt`) 并附带必要的单元测试。
- **代码规范**：保持接口文档注释、中文错误提示与统一日志字段。

如果 `gb` 对你有所帮助，欢迎 ⭐️ 支持，让更多团队受益。EOF
