package auth

import (
	"domain-admin/internal/service"
	"domain-admin/model"
	"domain-admin/pkg/db"
	"domain-admin/pkg/logger"
	"domain-admin/pkg/response"
	"domain-admin/pkg/validator"

	"github.com/gin-gonic/gin"
)

// AuthHandler 认证处理器
type AuthHandler struct {
	userService service.UserService
}

// NewAuthHandler 创建认证处理器
func NewAuthHandler() *AuthHandler {
	userService := service.NewUserService(db.GetDB("default"))
	return &AuthHandler{
		userService: userService,
	}
}

// Register 用户注册
// @Summary 用户注册
// @Description 用户注册接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.UserCreateRequest true "注册信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req model.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("参数绑定失败: %v", err)
		response.Error(c, 400, "参数格式错误")
		return
	}

	// 参数验证
	if err := validator.ValidateStruct(&req); err != nil {
		logger.Warnf("参数验证失败: %v", err)
		response.Error(c, 400, err.Error())
		return
	}

	user, err := h.userService.Register(&req)
	if err != nil {
		logger.Errorf("用户注册失败: %v", err)
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, user)
}

// Login 用户登录
// @Summary 用户登录
// @Description 用户登录接口
// @Tags 认证
// @Accept json
// @Produce json
// @Param request body model.UserLoginRequest true "登录信息"
// @Success 200 {object} response.Response{data=map[string]interface{}}
// @Failure 400 {object} response.Response
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req model.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("参数绑定失败: %v", err)
		response.Error(c, 400, "参数格式错误")
		return
	}

	// 参数验证
	if err := validator.ValidateStruct(&req); err != nil {
		logger.Warnf("参数验证失败: %v", err)
		response.Error(c, 400, err.Error())
		return
	}

	token, user, err := h.userService.Login(&req)
	if err != nil {
		logger.Errorf("用户登录失败: %v", err)
		response.Error(c, 401, err.Error())
		return
	}

	response.Success(c, gin.H{
		"token": token,
		"user":  user,
	})
}

// Logout 用户登出
// @Summary 用户登出
// @Description 用户登出，清除会话缓存
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, 401, "用户未登录")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Error(c, 401, "用户ID格式错误")
		return
	}

	if err := h.userService.Logout(uid); err != nil {
		logger.Errorf("用户登出失败: %v", err)
		response.Error(c, 500, err.Error())
		return
	}

	response.Success(c, "登出成功")
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的资料
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 401 {object} response.Response
// @Router /api/auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, 401, "用户未登录")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Error(c, 401, "用户ID格式错误")
		return
	}

	user, err := h.userService.GetProfile(uid)
	if err != nil {
		logger.Errorf("获取用户资料失败: %v", err)
		response.Error(c, 404, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateProfile 更新用户资料
// @Summary 更新用户资料
// @Description 更新当前登录用户的资料
// @Tags 认证
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UserUpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Router /api/auth/profile [put]
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, 401, "用户未登录")
		return
	}

	uid, ok := userID.(uint)
	if !ok {
		response.Error(c, 401, "用户ID格式错误")
		return
	}

	var req model.UserUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Warnf("参数绑定失败: %v", err)
		response.Error(c, 400, "参数格式错误")
		return
	}

	// 参数验证
	if err := validator.ValidateStruct(&req); err != nil {
		logger.Warnf("参数验证失败: %v", err)
		response.Error(c, 400, err.Error())
		return
	}

	// 普通用户不能修改角色和状态
	req.Role = ""
	req.Status = nil

	user, err := h.userService.UpdateProfile(uid, &req)
	if err != nil {
		logger.Errorf("更新用户资料失败: %v", err)
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, user)
}