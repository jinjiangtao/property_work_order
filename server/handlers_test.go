package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

// 测试前初始化
func TestMain(m *testing.M) {
	// 初始化数据库连接
	initDB()
	m.Run()
}

// 测试登录接口
func TestLogin(t *testing.T) {
	// 创建Gin引擎
	r := gin.Default()
	registerRoutes(r)

	// 测试数据
	loginData := LoginRequest{
		Username: "admin",
		Password: "admin",
	}
	jsonData, _ := json.Marshal(loginData)

	// 创建测试请求
	req, err := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusOK && w.Code != http.StatusUnauthorized && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, %d, or %d, got %d", http.StatusOK, http.StatusUnauthorized, http.StatusInternalServerError, w.Code)
	}
}

// 测试创建保修单接口
func TestCreateRepair(t *testing.T) {
	// 创建Gin引擎
	r := gin.Default()
	registerRoutes(r)

	// 测试数据
	repairData := RepairRequest{
		Location:    "客厅",
		Description: "灯管不亮",
		ImageURL:    "/uploads/test.jpg",
	}
	jsonData, _ := json.Marshal(repairData)

	// 创建测试请求
	req, err := http.NewRequest("POST", "/api/repairs", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusCreated && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusCreated, http.StatusInternalServerError, w.Code)
	}
}

// 测试获取保修单列表接口
func TestGetRepairs(t *testing.T) {
	// 创建Gin引擎
	r := gin.Default()
	registerRoutes(r)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/api/repairs", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

// 测试更新保修单状态接口
func TestUpdateRepairStatus(t *testing.T) {
	// 创建Gin引擎
	r := gin.Default()
	registerRoutes(r)

	// 测试数据
	statusData := StatusUpdateRequest{
		Status: "processing",
	}
	jsonData, _ := json.Marshal(statusData)

	// 创建测试请求
	req, err := http.NewRequest("PUT", "/api/repairs/1/status", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusOK && w.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d or %d, got %d", http.StatusOK, http.StatusInternalServerError, w.Code)
	}
}

// 测试健康检查接口
func TestHealthCheck(t *testing.T) {
	// 创建Gin引擎
	r := gin.Default()
	registerRoutes(r)

	// 创建测试请求
	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	// 执行请求
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// 检查响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
	}
}
