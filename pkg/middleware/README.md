# 中间件使用指南

本目录包含三个核心中间件：认证、追踪和权限控制。

## 1. 认证中间件 (auth.go)

### 功能
- JWT Token 验证
- Bearer Token 格式检查
- 用户角色提取
- 详细错误信息返回

### 使用方法
```go
// 在路由中使用
r.Use(middleware.JWTAuth())
```

### 配置要求
- 需要在请求头中添加：`Authorization: Bearer <your-jwt-token>`
- JWT Token 需要包含：userID、role、username 字段

### 错误代码
- `MISSING_AUTH_HEADER`: 缺少认证头
- `INVALID_AUTH_FORMAT`: 认证头格式错误
- `EMPTY_TOKEN`: Token为空
- `INVALID_TOKEN`: Token无效或已过期

## 2. OpenTelemetry 追踪中间件 (otel.go)

### 功能
- OpenTelemetry 追踪初始化
- OTLP HTTP 导出器支持
- 可配置采样率
- 健康检查端点过滤
- 优雅关闭

### 使用方法
```go
// 初始化
shutdown := middleware.InitTracer("domain-admin", "localhost:4318")
defer shutdown()

// 在路由中使用
r.Use(middleware.OTLPMiddleware())
```

### 环境变量配置
```bash
# OTLP端点
OTEL_EXPORTER_OTLP_ENDPOINT=localhost:4318

# 采样率 (0.0-1.0)
OTEL_SAMPLE_RATE=0.5

# 是否使用不安全连接
OTEL_INSECURE=true

# 服务信息
ENVIRONMENT=production
SERVICE_VERSION=1.0.0
```

### 过滤的端点
- `/health*` - 健康检查
- `/metrics*` - 监控指标
- `/swagger*` - API文档

## 3. RBAC 权限控制中间件 (rbac.go)

### 功能
- 基于角色的访问控制
- Casbin 权限模型支持
- 动态策略重载
- 详细的权限检查日志
- 错误处理和状态码

### 使用方法
```go
// 初始化RBAC系统
err := middleware.InitRBAC("", "") // 使用默认配置
if err != nil {
    log.Fatal("Failed to initialize RBAC:", err)
}

// 在路由中使用（需要在JWTAuth之后）
r.Use(middleware.JWTAuth())
r.Use(middleware.RBACMiddleware())
```

### 配置文件

#### 权限模型 (configs/rbac_model.conf)
```ini
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && keyMatch(r.obj, p.obj) && (r.act == p.act || p.act == "*")
```

#### 权限策略 (configs/rbac_policy.csv)
```csv
p, admin, /api/*, *
p, user, /api/user/*, GET
p, user, /api/profile, *
g, alice, admin
g, bob, user
```

### API
- `InitRBAC(modelPath, policyPath)` - 初始化RBAC系统
- `IsRBACInitialized()` - 检查是否已初始化
- `ReloadPolicy()` - 重新加载权限策略

### 错误代码
- `RBAC_NOT_INITIALIZED`: RBAC系统未初始化
- `MISSING_USER_ROLE`: 用户角色缺失
- `INVALID_ROLE_FORMAT`: 角色格式无效
- `RBAC_ENFORCEMENT_ERROR`: 权限检查错误
- `ACCESS_DENIED`: 访问被拒绝

## 使用顺序

正确的中间件使用顺序：

```go
// 1. 初始化追踪
shutdown := middleware.InitTracer("domain-admin", "localhost:4318")
defer shutdown()

// 2. 初始化RBAC
err := middleware.InitRBAC("", "")
if err != nil {
    log.Fatal(err)
}

// 3. 设置路由中间件
r.Use(middleware.OTLPMiddleware())  // 追踪
r.Use(middleware.JWTAuth())        // 认证
r.Use(middleware.RBACMiddleware()) // 权限控制
```

## 最佳实践

1. **错误处理**: 所有中间件都返回结构化的错误信息
2. **日志记录**: 使用统一的日志格式，包含关键信息
3. **配置管理**: 使用环境变量进行配置
4. **性能优化**: 追踪中间件过滤不必要的端点
5. **安全性**: JWT认证和RBAC权限控制配合使用

## 故障排查

### 认证失败
- 检查Authorization头格式是否正确
- 验证JWT Token是否有效
- 查看日志中的具体错误代码

### 追踪不工作
- 确认OTLP端点配置正确
- 检查网络连接
- 查看日志中的初始化信息

### 权限被拒绝
- 检查用户的角色设置
- 验证权限策略配置
- 使用日志中的调试信息