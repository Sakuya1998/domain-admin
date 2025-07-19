package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"domain-admin/api"
	"domain-admin/pkg/config"
	"domain-admin/pkg/db"
	"domain-admin/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var (
	router *gin.Engine
	adminToken string
	userToken string
	testUserPrefix string
)

// TestMain 测试主函数，初始化测试环境
func TestMain(m *testing.M) {
	// 切换到主项目目录以访问配置文件
	originalDir, _ := os.Getwd()
	os.Chdir("../")
	defer os.Chdir(originalDir)
	
	// 生成唯一的测试用户前缀
	rand.Seed(time.Now().UnixNano())
	testUserPrefix = fmt.Sprintf("testuser_%d_", rand.Intn(100000))
	
	// 初始化配置
	cfg := config.InitConfig()
	
	// 初始化日志
	logger.InitLogger(cfg.Log)
	
	// 初始化数据库
	db.InitDB(cfg.Database)
	
	// 设置Gin为测试模式
	gin.SetMode(gin.TestMode)
	
	// 初始化路由
	router = gin.Default()
	api.RegisterRoutes(router)
	
	// 运行测试
	m.Run()
}

// makeRequest 发送HTTP请求的辅助函数
func makeRequest(method, url string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody []byte
	if body != nil {
		reqBody, _ = json.Marshal(body)
	}
	
	req := httptest.NewRequest(method, url, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// TestAdminLogin 测试管理员登录
func TestAdminLogin(t *testing.T) {
	body := map[string]string{
		"username": "admin",
		"password": "password",
	}
	
	w := makeRequest("POST", "/api/auth/login", body, "")
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
	
	// 保存token供后续测试使用
	data := response["data"].(map[string]interface{})
	adminToken = data["token"].(string)
	assert.NotEmpty(t, adminToken)
}

// TestUserRegisterWithoutRole 测试用户注册（无Role字段）
func TestUserRegisterWithoutRole(t *testing.T) {
	body := map[string]string{
		"username": testUserPrefix + "001",
		"password": "password123",
		"email":    testUserPrefix + "001@test.com",
		"nickname": "Test User 001",
	}
	
	w := makeRequest("POST", "/api/auth/register", body, "")
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
	
	// 验证默认role为user
	if data, ok := response["data"].(map[string]interface{}); ok {
		assert.Equal(t, "user", data["role"])
	}
}

// TestUserRegisterWithRole 测试用户注册（有Role字段）
func TestUserRegisterWithRole(t *testing.T) {
	body := map[string]string{
		"username": testUserPrefix + "002",
		"password": "password123",
		"email":    testUserPrefix + "002@test.com",
		"nickname": "Test User 002",
		"role":     "user",
	}
	
	w := makeRequest("POST", "/api/auth/register", body, "")
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
}

// TestUserLogin 测试普通用户登录
func TestUserLogin(t *testing.T) {
	body := map[string]string{
		"username": testUserPrefix + "001",
		"password": "password123",
	}
	
	w := makeRequest("POST", "/api/auth/login", body, "")
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
	
	// 保存token供后续测试使用
	if data, ok := response["data"].(map[string]interface{}); ok {
		userToken = data["token"].(string)
		assert.NotEmpty(t, userToken)
	}
}

// TestUnauthorizedAccess 测试无Authorization头访问
func TestUnauthorizedAccess(t *testing.T) {
	w := makeRequest("GET", "/api/user/profile", nil, "")
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(401), response["code"])
}

// TestInvalidTokenAccess 测试无效token访问
func TestInvalidTokenAccess(t *testing.T) {
	w := makeRequest("GET", "/api/user/profile", nil, "invalid_token")
	
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(401), response["code"])
}

// TestUserProfile 测试普通用户获取资料
func TestUserProfile(t *testing.T) {
	w := makeRequest("GET", "/api/user/profile", nil, userToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
}

// TestUpdateUserProfile 测试普通用户更新资料
func TestUpdateUserProfile(t *testing.T) {
	body := map[string]string{
		"nickname": "Updated Test User 001",
	}
	
	w := makeRequest("PUT", "/api/user/profile", body, userToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
}

// TestGetUserList 测试获取用户列表（管理员功能）
func TestGetUserList(t *testing.T) {
	w := makeRequest("GET", "/api/admin/users", nil, adminToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
}

// TestCreateUser 测试创建用户（管理员功能）
func TestCreateUser(t *testing.T) {
	body := map[string]string{
		"username": testUserPrefix + "admin001",
		"password": "password123",
		"email":    testUserPrefix + "admin001@test.com",
		"nickname": "Admin User 001",
		"role":     "admin",
	}
	
	w := makeRequest("POST", "/api/admin/users", body, adminToken)
	
	assert.Equal(t, http.StatusOK, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(200), response["code"])
	assert.Equal(t, "success", response["message"])
}

// TestUpdateNonExistentUser 测试更新不存在的用户
func TestUpdateNonExistentUser(t *testing.T) {
	body := map[string]string{
		"nickname": "test",
	}
	
	w := makeRequest("PUT", "/api/admin/users/99999", body, adminToken)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(404), response["code"])
}

// TestDeleteNonExistentUser 测试删除不存在的用户
func TestDeleteNonExistentUser(t *testing.T) {
	w := makeRequest("DELETE", "/api/admin/users/99999", nil, adminToken)
	
	assert.Equal(t, http.StatusNotFound, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(404), response["code"])
}

// TestForbiddenAccess 测试权限不足访问（普通用户访问管理员接口）
func TestForbiddenAccess(t *testing.T) {
	w := makeRequest("GET", "/api/admin/users", nil, userToken)
	
	assert.Equal(t, http.StatusForbidden, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(403), response["code"])
}

// TestValidationError 测试参数验证错误
func TestValidationError(t *testing.T) {
	body := map[string]string{
		"username": "",
		"password": "123",
		"email":    "invalid-email",
	}
	
	w := makeRequest("POST", "/api/auth/register", body, "")
	
	assert.Equal(t, http.StatusBadRequest, w.Code)
	
	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(t, err)
	assert.Equal(t, float64(400), response["code"])
}

// BenchmarkAdminLogin 管理员登录性能测试
func BenchmarkAdminLogin(b *testing.B) {
	body := map[string]string{
		"username": "admin",
		"password": "password",
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeRequest("POST", "/api/auth/login", body, "")
	}
}

// BenchmarkUserProfile 用户资料获取性能测试
func BenchmarkUserProfile(b *testing.B) {
	if userToken == "" {
		b.Skip("需要先运行用户登录测试")
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		makeRequest("GET", "/api/user/profile", nil, userToken)
	}
}