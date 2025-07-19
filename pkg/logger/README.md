# 日志系统使用说明

## 功能特性

- **多级别日志**: 支持 debug、info、warn、error、fatal 级别
- **多种输出**: 支持控制台、文件、或同时输出
- **日志轮转**: 基于 lumberjack 的自动日志轮转
- **灵活配置**: 支持 YAML 配置、环境变量、代码配置
- **开发友好**: 开发模式下彩色输出
- **结构化日志**: JSON 格式便于日志收集和分析

## 快速开始

### 1. 基础使用

```go
package main

import (
    "github.com/your-project/pkg/logger"
)

func main() {
    // 使用默认配置初始化
    logger.InitLogger()
    defer logger.Sync()
    
    logger.Info("应用启动成功")
    logger.Error("发生错误", zap.Error(err))
}
```

### 2. 自定义配置

```go
import (
    "github.com/your-project/pkg/logger"
    "github.com/spf13/viper"
)

func main() {
    // 从配置文件加载
    viper.SetConfigFile("configs/logger.yaml")
    viper.ReadInConfig()
    
    var config logger.Config
    viper.UnmarshalKey("logger", &config)
    
    logger.InitLoggerWithConfig(&config)
    defer logger.Sync()
}
```

### 3. 代码配置

```go
config := &logger.Config{
    Level:       "debug",
    Format:      "console",
    Output:      "console",
    Development: true,
}

logger.InitLoggerWithConfig(config)
```

## 配置说明

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| level | string | "info" | 日志级别: debug, info, warn, error |
| format | string | "json" | 输出格式: json, console |
| output | string | "both" | 输出方式: file, console, both |
| file_path | string | "logs/app.log" | 日志文件路径 |
| max_size | int | 100 | 单个文件最大大小(MB) |
| max_backups | int | 3 | 保留的最大备份文件数 |
| max_age | int | 28 | 保留的最大天数 |
| compress | bool | true | 是否压缩旧日志 |
| development | bool | false | 开发模式(彩色输出) |

## 日志级别

```go
// Debug 级别 - 调试信息
logger.Debug("调试信息", zap.String("key", "value"))

// Info 级别 - 一般信息
logger.Info("用户登录", zap.String("username", "admin"))

// Warn 级别 - 警告信息
logger.Warn("磁盘空间不足", zap.Int("free_space", 1024))

// Error 级别 - 错误信息
logger.Error("数据库连接失败", zap.Error(err))

// Fatal 级别 - 致命错误
logger.Fatal("无法启动服务", zap.Error(err))
```

## 上下文日志

```go
// 添加固定字段
logger := logger.With(
    zap.String("service", "user-service"),
    zap.String("version", "1.0.0"),
)

logger.Info("处理请求", zap.String("path", "/api/users"))
```

## 环境变量支持

可以通过环境变量覆盖配置：

```bash
# 设置日志级别
export LOGGER_LEVEL=debug

# 设置输出格式
export LOGGER_FORMAT=console

# 设置开发模式
export LOGGER_DEVELOPMENT=true
```

## 最佳实践

1. **初始化时机**: 在应用启动时尽早初始化
2. **延迟同步**: 使用 `defer logger.Sync()` 确保日志写入
3. **错误处理**: 总是记录错误堆栈
4. **敏感信息**: 避免在日志中记录密码等敏感信息
5. **日志轮转**: 生产环境务必开启日志轮转

## 示例配置

### 开发环境
```yaml
logger:
  level: "debug"
  format: "console"
  output: "console"
  development: true
```

### 生产环境
```yaml
logger:
  level: "info"
  format: "json"
  output: "both"
  file_path: "/var/log/app/app.log"
  max_size: 500
  max_backups: 7
  max_age: 30
  compress: true
  development: false
```

## 常见问题

1. **日志文件权限问题**: 确保应用有写入日志目录的权限
2. **磁盘空间**: 合理设置 max_size 和 max_backups 参数
3. **性能影响**: 在高并发场景下，考虑使用异步日志
4. **日志分析**: JSON 格式便于 ELK、Grafana 等工具分析