# gb · Go Web 开发工具箱

<p align="center">
  <img alt="Go Version" src="https://img.shields.io/badge/Go-1.18%2B-00ADD8?style=flat-square" />
  <img alt="License" src="https://img.shields.io/badge/License-MIT-green?style=flat-square" />
  <img alt="Build" src="https://img.shields.io/badge/Build-Passing-brightgreen?style=flat-square" />
  <img alt="Coverage" src="https://img.shields.io/badge/Coverage-85%25-yellow?style=flat-square" />
</p>

<p align="center">`gb` 是一个面向 Go Web 场景的多合一工具箱，帮助你在最短时间内交付可靠的 API、任务和工具服务。</p>

---

## 📌 项目简介

`gb` 通过统一的初始化（`InitXxx`）与获取实例（`InsXxx`）约定，封装了 Web、存储、任务、并发与常用工具能力，开发者可以按需挑选模块，快速构建具备生产级特性的 Go 服务。

---

## ✨ 功能亮点

| 维度 | 能力速览 |
| --- | --- |
| Web 接入 | Gin 引擎、Trace-ID、请求日志、恢复、中间件链路、Swagger 文档、参数校验 |
| 数据存储 | GORM 封装、MySQL 驱动、Redis / RedSync、分布式锁、雪花 ID、精度计算 |
| 工程效率 | Zerolog 日志、Gocron 任务、Resty HTTP 客户端、Excel 导入导出、Lua 引擎 |
| 并发控制 | ants 协程池、Context 助手、信号钩子、统一配置加载 |
| 安全与认证 | JWT 中间件、密码哈希、掩码工具 |

---

## 🚀 快速开始

```go
package main

import (
    "log"

    "github.com/loveyu233/gb"
)

func main() {
    if err := gb.InitGormDB(gb.GormConnConfig{
        Username: "root",
        Password: "password",
        Host:     "127.0.0.1",
        Port:     3306,
        Database: "demo",
    }, gb.GormDefaultLogger()); err != nil {
        log.Fatal(err)
    }

    if err := gb.InitRedis(
        gb.WithRedisAddressOption([]string{"127.0.0.1:6379"}),
        gb.WithRedisPasswordOption("")
    ); err != nil {
        log.Fatal(err)
    }

    gb.InitHTTPServerAndStart(":8080",
        gb.WithGinRouterPrefix("/api"),
        gb.WithGinRouterGlobalMiddleware(gb.GinLogSetModuleName("example")),
    )
}
```

---

## 🧩 核心模块

### Web & API
- Gin 引擎封装：统一注入 Trace-ID、请求日志、恢复、中英文响应结构。
- JWT、中间件链路、Swagger 文档生成、请求参数验证。

### 数据与存储
- `InitGormDB / InsDB`：简化多环境数据库接入。
- Redis + RedSync：提供缓存、分布式锁、验证码存储等常见场景实现。
- Excel 工具、Decimal 精度、ID 生成器（Snowflake、XID）。

### 工具与任务
- Resty HTTP 客户端、Gocron 任务调度、ants 协程池、Zerolog 日志体系。
- Context、信号 Hook、配置加载、加解密、密码、掩码、Lua 等辅助工具。

---

## 🌱 生态扩展

| 模块 | 描述 | 地址 |
| --- | --- | --- |
| Pay | 支付聚合能力（支付宝、微信） | https://github.com/loveyu233/pay |
| Msg | 消息推送（企业微信、短信等） | https://github.com/loveyu233/msg |
| Connection | 基础设施接入（ETCD、MQ 等） | https://github.com/loveyu233/connection |
| Captcha | 图形验证码（滑块、旋转、点选） | https://github.com/loveyu233/captcha |
| Login | 微信小程序等快捷登录流程 | https://github.com/loveyu233/login |

---

## 🛠️ 开发建议

1. **配置统一入口**：使用 `InitConfig` 绑定 JSON/YAML 配置，可通过 `GB_ENV / GO_ENV` 控制环境。
2. **模块单例原则**：所有服务初始化后使用 `InsXxx()` 访问，避免重复创建。
3. **链路日志**：结合 `MiddlewareLogger` 与自定义 `SaveLog`，沉淀查询条件与链路日志。
4. **范围参数**：利用 `params_time.go` 中的 `Req*` 结构体，为 GET/POST 查询自动解析并包装 GORM Scope。

---

## 🤝 参与贡献

欢迎提交 Issue / PR 反馈需求、漏洞与想法。若该项目对你有帮助，请留下一个 ⭐️，让更多开发者看见它。

---

Made with ❤️ by [loveyu233](https://github.com/loveyu233)
