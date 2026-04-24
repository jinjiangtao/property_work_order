package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func registerRoutes(r *gin.Engine) {
	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	// API路由组
	api := r.Group("/api")
	{
		// 认证相关
		auth := api.Group("/auth")
		{
			auth.POST("/login", login)
		}

		// 保修相关
		repairs := api.Group("/repairs")
		{
			repairs.POST("/", createRepair)
			repairs.GET("/", getRepairs)
			repairs.GET("/:id", getRepair)
			repairs.PUT("/:id/status", updateRepairStatus)
		}

		// 用户相关
		users := api.Group("/users")
		{
			users.POST("/", createUser)
			users.GET("/", getUsers)
		}
	}

	// 文件上传
	r.POST("/upload", uploadImage)
}
