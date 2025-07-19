# Domain Admin 前端项目

基于 Vue 3 + TypeScript + Element Plus 的域名管理系统前端。

## 环境配置

### 后端API地址配置

项目支持通过环境变量配置后端API地址：

#### 开发环境
编辑 `.env` 文件：
```bash
# 后端API基础地址
VITE_API_BASE_URL=http://localhost:8000
# 应用标题
VITE_APP_TITLE=Domain Admin
```

#### 生产环境
编辑 `.env.production` 文件：
```bash
# 后端API基础地址
VITE_API_BASE_URL=https://your-domain.com
# 应用标题
VITE_APP_TITLE=Domain Admin
```

### 配置说明

1. **VITE_API_BASE_URL**: 后端API的基础地址
   - 开发环境默认: `http://localhost:8000`
   - 生产环境需要修改为实际的后端服务地址
   - 不需要包含 `/api` 路径，系统会自动添加

2. **VITE_APP_TITLE**: 应用标题，显示在浏览器标签页

### 代理配置

开发环境下，Vite会自动将 `/api/*` 的请求代理到配置的后端地址。

生产环境下，需要确保前端和后端部署在同一域名下，或者配置CORS。

## 开发

```bash
# 安装依赖
yarn install

# 启动开发服务器
yarn dev

# 构建生产版本
yarn build
```

## 部署

1. 修改 `.env.production` 中的 `VITE_API_BASE_URL` 为实际的后端地址
2. 运行 `yarn build` 构建项目
3. 将 `dist` 目录部署到Web服务器

## 注意事项

- 确保后端服务已启动并可访问
- 生产环境需要配置CORS或使用同域部署
- 修改环境变量后需要重启开发服务器