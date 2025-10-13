# gb: A Go Web Development Toolkit
# Go Web 开发工具包

<div align="center">

![Go Version](https://img.shields.io/badge/Go-1.18+-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Coverage](https://img.shields.io/badge/coverage-85%25-yellow)

**一个为 Go Web 开发量身打造的、开箱即用的工具库，旨在简化开发流程，提升效率。**

</div>

## 核心设计理念

`gb` 遵循统一的初始化与调用模式，以简化资源管理和保证单例实例。

- **初始化**: 所有模块的初始化函数均以 `Init` 为前缀，例如 `gb.InitGormDB()`。
- **获取实例**: 初始化后，通过 `Ins` 前缀的函数获取该模块的唯一实例，例如 `gb.InsDB()`。

```go
package main

import (
	"fmt"
	"github.com/loveyu233/gb"
)

func main() {
	// 示例：初始化并使用 GORM
	err := gb.InitGormDB(gb.GormConf{
		// ... 数据库配置
	})
	if err != nil {
		panic("Failed to connect to database")
	}

	// 在应用的任何地方，通过 InsDB() 获取 GORM 实例
	db := gb.InsDB()
	fmt.Println(db)
}
```

---

## 🚀 功能特性

### 🌐 Web 与 API 框架 (基于 Gin)
- **Gin 快速启动**: 内置 Engine，并集成了一系列即用型中间件。
- **中间件套件**: 提供 `CORS`, `Trace-ID`, `请求日志`, `Recovery` 等常用中间件。
- **JWT 认证**: 基于 `jwt/v5` 的完整解决方案，支持 Token 生成、解析与刷新。
- **Swagger 文档**: 自动生成 API 文档，方便接口调试与测试。
- **参数验证**: 集成 `validator/v10`，支持自定义规则和多语言错误提示。
- **标准化响应**: 提供统一的 JSON 响应封装，简化 API 输出。

### 🗃️ 数据与存储
- **ORM 数据库**: 深度集成 `GORM`，提供便捷的数据库操作。
- **Redis**: 内置 `go-redis` 客户端，简化缓存、队列等操作。
- **RedSync 分布式锁**: 基于 Redis 实现的分布式锁，确保并发安全。
- **Excel 操作**: 使用 `excelize/v2`，轻松实现 Excel 导入导出及与结构体的映射。
- **高精度计算**: 内置 `Decimal` 处理，避免浮点数计算的精度问题。

### 🛠️ 核心工具
- **结构化日志**: 基于 `zerolog` 的高性能结构化日志，支持多种输出格式。
- **Cron 定时任务**: 集成 `gocron`，提供流畅的 API 来安排定时作业。
- **HTTP 客户端**: 集成 `Resty`，提供链式调用的 HTTP 请求体验。
- **Goroutine 池**: 内置 `ants` 协程池，有效管理和复用 Goroutine，防止大规模泄露。
- **ID 生成器**: 提供 `Snowflake` 雪花算法和 `XID`，适用于分布式ID生成。
- **密码安全**: 内置 `Bcrypt` 密码哈希与验证功能。
- **Lua 脚本**: 集成 Lua 引擎，方便执行动态脚本。
- **辅助函数库**: 预集成 `lo` (Lodash for Go), `gjson`, `copier`, `cast` 等热门工具库，极大提升开发效率。

---

## 🌱 生态系统

为了保持 `gb` 核心库的轻量和专注，部分功能已被拆分到独立的模块中，您可以按需引入：

- **支付集成 (Pay)**: 支付宝、微信支付等功能。
  - [https://github.com/loveyu233/pay](https://github.com/loveyu233/pay)
- **消息推送 (Msg)**: 企业微信机器人、短信发送等。
  - [https://github.com/loveyu233/msg](https://github.com/loveyu233/msg)
- **服务连接 (Connection)**: ETCD、RocketMQ 等。
  - [https://github.com/loveyu233/connection](https://github.com/loveyu233/connection)

---

**⭐ 如果这个项目对你有帮助，请给我们一个 Star！⭐**

Made with ❤️ by [loveyu233](https://github.com/loveyu233)