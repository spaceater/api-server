## 安装和运行

### 前置要求

- Go 1.21 或更高版本

### 安装依赖

```bash
go mod tidy
```

### 运行服务器

```bash
go run main.go
```

服务器将在 `http://127.0.0.1:2998` 启动。

## API 接口

### 页面访问量统计

- **URL**: `GET /api/page_view`
- **功能**: 获取并增加页面访问量计数
- **响应**: 
  ```json
  {
    "page_view": 123
  }
  ```

### WebSocket 连接

- **URL**: `ws://127.0.0.1:2998/ws/ismismcube_online`
- **功能**: 建立 WebSocket 连接，接收在线用户数量更新
- **消息格式**:
  ```json
  {
    "online_count": 5
  }
  ```

### 静态文件

- **Favicon**: `GET /favicon.ico`
- **404 页面**: 访问不存在的路径时返回自定义 404 页面

## 配置

可以通过环境变量进行配置：

- `PORT`: 服务器端口（默认: 2998）
- `PAGE_VIEW_FILE`: 页面访问量存储文件（默认: resources/page-view.txt）

## 项目结构

### 核心文件
- `main.go` - 应用入口，启动服务器
- `go.mod` - Go 模块依赖管理

### 内部模块
- `internal/api/router.go` - API 路由注册和初始化
- `internal/config/config.go` - 配置管理，支持环境变量
- `internal/handlers/` - 请求处理器
  - `pageview.go` - 页面访问量统计处理
  - `ping.go` - 健康检查接口
- `internal/middleware/cors.go` - CORS 和缓存控制中间件
- `internal/router/router.go` - 自定义路由系统实现
- `internal/websocket/ismismcube.go` - WebSocket 连接管理

### 资源文件
- `resources/page-view.txt` - 页面访问量数据存储
- `test/` - 测试相关文件

## 许可证

MIT License
