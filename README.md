# Go Web 开发工具包

<div align="center">


![Go Version](https://img.shields.io/badge/Go-1.18+-blue)
![License](https://img.shields.io/badge/license-MIT-green)
![Build Status](https://img.shields.io/badge/build-passing-brightgreen)
![Coverage](https://img.shields.io/badge/coverage-85%25-yellow)

**Go Web 开发工具包简称: gb**

> 所有的初始化均使用      `Init` 为前缀.
> 所有初始化后调用均使用   `Ins` 为前缀.
> 
> 例如:gb.InitGormDB() 使用则: gb.InsDB 其他均如此

</div>

## 🚀 特性一览

### 核心功能

- 🏗️ **Gin快速中间件** - 提供快速启动gin服务包含各种中间件
- 🔐 **JWT 认证** - 内置完整的 JWT 认证解决方案，支持 Token 刷新和黑名单机制
- 📝 **Swagger 文档** - 自动生成 API 文档，支持在线调试和接口测试
- 🎯 **智能路由** - 优雅的路由管理，支持路由分组、中间件链式调用
- ⏱️ **请求追踪** - 全链路请求追踪，便于性能监控和问题排查

### 第三方集成

- 💰 **支付宝集成** - 开箱即用的支付宝支付功能，支持 APP、Web、扫码支付、阿里的短信发送
- 💳 **微信支付** - 完整的微信支付解决方案，包含小程序支付
- 🔑 **小程序登录** - 支持微信小程序一键登录和用户信息获取
- 📚 **企业微信机器人** - 快速启动一个企业微信机器人，实现消息推送 

### 数据处理

- 📊 **多数据库支持** - 集成 GORM、Redis、ETCD，支持多种数据存储方案
- 🔍 **参数验证** - 强大的参数验证功能，支持中文错误提示
- 🕒 **时间处理** - 专为中国时区优化的时间处理工具
- 💰 **精确计算** - 内置 Decimal 处理，避免浮点数精度问题
- 📚 **excel处理** - 支持快速的excel和结构体互相转化

### 开发工具

- 🧰 **常用库集成** - 预集成 ants、gjson、copier、lo 等热门工具库
- 🚨 **错误处理** - 统一的错误处理机制，标准化的错误响应格式
- 📝 **结构化日志** - 高性能的结构化日志记录，支持多种输出格式
- 🔄 **优雅关闭** - 支持服务优雅关闭，确保请求完整处理
- 🔑 **lua脚本** - 集成了常用的lua脚本,快速使用

**⭐ 如果这个项目对你有帮助，请给我们一个 Star！⭐**

Made with ❤️ by [loveyu233](https://github.com/loveyu233)
