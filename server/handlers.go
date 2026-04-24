package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 登录请求结构
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// 保修请求结构
type RepairRequest struct {
	Location    string `json:"location" binding:"required"`
	Description string `json:"description" binding:"required"`
	ImageURL    string `json:"image_url"`
}

// 状态更新请求结构
type StatusUpdateRequest struct {
	Status string `json:"status" binding:"required,oneof=pending processing completed"`
}

// 用户创建请求结构
type UserRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=admin owner"`
}

// 登录处理
func login(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var user struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
		Role     string `json:"role"`
	}

	err := db.QueryRow("SELECT id, username, role FROM users WHERE username = ? AND password = ?", req.Username, req.Password).Scan(&user.ID, &user.Username, &user.Role)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"user":    user,
	})
}

// 创建保修单
func createRepair(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	var req RepairRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 这里简化处理，假设用户ID为1
	userID := 1

	result, err := db.Exec("INSERT INTO repairs (user_id, location, description, image_url) VALUES (?, ?, ?, ?)", userID, req.Location, req.Description, req.ImageURL)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create repair request"})
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repair ID"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Repair request created successfully",
		"id":      id,
	})
}

// 获取保修单列表
func getRepairs(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	rows, err := db.Query("SELECT id, user_id, location, description, image_url, status, created_at FROM repairs")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get repair requests"})
		return
	}
	defer rows.Close()

	var repairs []map[string]interface{}
	for rows.Next() {
		var id, userID int
		var location, description, imageURL, status string
		var createdAt time.Time

		if err := rows.Scan(&id, &userID, &location, &description, &imageURL, &status, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan repair request"})
			return
		}

		repairs = append(repairs, map[string]interface{}{
			"id":          id,
			"user_id":     userID,
			"location":    location,
			"description": description,
			"image_url":   imageURL,
			"status":      status,
			"created_at":  createdAt,
		})
	}

	c.JSON(http.StatusOK, repairs)
}

// 获取单个保修单
func getRepair(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repair ID"})
		return
	}

	var repair struct {
		ID          int       `json:"id"`
		UserID      int       `json:"user_id"`
		Location    string    `json:"location"`
		Description string    `json:"description"`
		ImageURL    string    `json:"image_url"`
		Status      string    `json:"status"`
		CreatedAt   time.Time `json:"created_at"`
	}

	err = db.QueryRow("SELECT id, user_id, location, description, image_url, status, created_at FROM repairs WHERE id = ?", id).Scan(&repair.ID, &repair.UserID, &repair.Location, &repair.Description, &repair.ImageURL, &repair.Status, &repair.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Repair request not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, repair)
}

// 更新保修单状态
func updateRepairStatus(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid repair ID"})
		return
	}

	var req StatusUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err = db.Exec("UPDATE repairs SET status = ? WHERE id = ?", req.Status, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update repair status"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Repair status updated successfully"})
}

// 创建用户
func createUser(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	var req UserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := db.Exec("INSERT INTO users (username, password, role) VALUES (?, ?, ?)", req.Username, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// 获取用户列表
func getUsers(c *gin.Context) {
	if db == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection not available"})
		return
	}

	rows, err := db.Query("SELECT id, username, role, created_at FROM users")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get users"})
		return
	}
	defer rows.Close()

	var users []map[string]interface{}
	for rows.Next() {
		var id int
		var username, role string
		var createdAt time.Time

		if err := rows.Scan(&id, &username, &role, &createdAt); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to scan user"})
			return
		}

		users = append(users, map[string]interface{}{
			"id":         id,
			"username":   username,
			"role":       role,
			"created_at": createdAt,
		})
	}

	c.JSON(http.StatusOK, users)
}

// 上传图片
func uploadImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No image file provided"})
		return
	}

	// 创建上传目录
	uploadDir := "./uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// 生成文件名
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%d%s", time.Now().Unix(), ext)
	filePath := filepath.Join(uploadDir, filename)

	// 保存文件
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save image"})
		return
	}

	// 返回文件URL
	fileURL := fmt.Sprintf("/uploads/%s", filename)
	c.JSON(http.StatusOK, gin.H{"url": fileURL})
}
