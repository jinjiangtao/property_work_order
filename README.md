# 物业保修系统

## 项目结构

```
property_work_order/
├── server/          # API服务端
│   ├── main.go      # 主入口文件
│   ├── database.go  # 数据库初始化
│   ├── routes.go    # 路由注册
│   ├── handlers.go  # 处理函数
│   └── go.mod       # 依赖管理
├── web/             # 前端页面
│   ├── admin/       # 后台管理页面
│   │   └── index.html
│   └── h5/          # 业主提交页面
│       └── index.html
└── README.md        # 项目说明
```

## 技术栈

- **服务端**：Go + Gin + MySQL
- **前端**：HTML + Vue + Axios

## 功能特点

1. **业主功能**：
   - 登录系统
   - 提交保修申请（位置描述、问题描述、图片上传）

2. **物业管理人员功能**：
   - 登录系统
   - 查看所有保修单
   - 更新保修单状态（待处理、处理中、已完成）

## 快速开始

### 服务端启动

1. 进入server目录：
   ```bash
   cd server
   ```

2. 安装依赖：
   ```bash
   go mod tidy
   ```

3. 启动服务：
   ```bash
   go run main.go
   ```

   服务将在 `http://localhost:8080` 运行

### 前端访问

- **业主页面**：`web/h5/index.html`
- **后台管理页面**：`web/admin/index.html`

## 数据库配置

默认数据库配置：
- 用户名：root
- 密码：root
- 数据库名：property_work_order
- 端口：3306

## 注意事项

1. 请确保MySQL服务已启动
2. 首次运行时会自动创建数据库表结构
3. 图片上传功能会在server目录下创建uploads文件夹

物业工单管理系统
