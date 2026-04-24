package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库连接
	initDB()

	// 创建Gin引擎
	r := gin.Default()

	// 注册路由
	registerRoutes(r)

	// 启动服务器
	serverAddr := ":8080"
	fmt.Printf("Server running at http://localhost%s\n", serverAddr)
	if err := r.Run(serverAddr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
