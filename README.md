## 功能特性

- **WebSocket 连接管理**: 实时统计在线用户数量
- **页面访问量统计**: 提供 API 接口统计页面访问次数
- **静态文件服务**: 支持 favicon.ico 和自定义 404 页面
- **CORS 支持**: 跨域请求支持
- **错误处理**: 完善的错误处理机制

## 项目结构

```
ismismcube-backend/
├── main.go                    # 主程序入口
├── go.mod                     # Go 模块文件
├── go.sum                     # 依赖版本锁定文件
├── internal/                  # 内部包
│   ├── api/                  # API 路由声明
│   │   ├── router.go         # 所有 API 路由定义
│   │   └── middleware.go     # API 相关中间件
│   ├── config/               # 配置管理
│   │   └── config.go
│   ├── handlers/             # 请求处理器
│   │   ├── pageview.go       # 页面访问量处理
│   │   ├── static.go         # 静态文件处理
│   │   └── routes.go         # 已废弃，路由已移至 api 包
│   ├── middleware/           # 中间件
│   │   └── cors.go           # CORS 中间件
│   └── websocket/            # WebSocket 管理
│       └── manager.go        # WebSocket 连接管理器
├── resources/                # 静态资源
│   ├── favicon.ico          # 网站图标
│   └── 404.html             # 404 错误页面
└── README.md                # 项目说明
```

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
- `GIN_MODE`: Gin 运行模式（debug/release，默认: debug）
- `PAGE_VIEW_FILE`: 页面访问量存储文件（默认: page-view.txt）

## 开发

### 项目结构说明

- `internal/api`: API 路由声明，包含所有接口的路由定义
- `internal/config`: 配置管理，支持环境变量配置
- `internal/handlers`: 请求处理器，包含各种 API 处理逻辑
- `internal/middleware`: 中间件，如 CORS、缓存控制等
- `internal/websocket`: WebSocket 连接管理，支持多客户端连接

### 添加新功能

1. 在 `internal/handlers` 中创建新的处理器
2. 在 `internal/api/routes.go` 中注册路由
3. 如需要，在 `internal/middleware` 中添加中间件

## 许可证

MIT License
