package user

import (
	"domain-admin/internal/service"
	"domain-admin/model"
	"domain-admin/pkg/db"
	"domain-admin/pkg/logger"
	"domain-admin/pkg/pagination"
	"domain-admin/pkg/response"
	"domain-admin/pkg/validator"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler 创建用户处理器
func NewUserHandler() *UserHandler {
	userService := service.NewUserService(db.GetDB("default"))
	return &UserHandler{
		userService: userService,
	}
}

// GetUserList 获取用户列表（管理员功能）
// @Summary 获取用户列表
// @Description 获取用户列表，仅管理员可访问
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param offset query int false "偏移量" default(0)
// @Param limit query int false "每页数量" default(10)
// @Param orderBy query string false "排序字段" default(id)
// @Param sort query string false "排序方式" default(asc)
// @Success 200 {object} response.Response{data=pagination.PageResult}
// @Failure 403 {object} response.Response
// @Router /api/users [get]
func (h *UserHandler) GetUserList(c *gin.Context) {
	page := pagination.New(c)

	result, err := h.userService.GetUserList(page)
	if err != nil {
		logger.Errorf("获取用户列表失败: %v", err)
		response.Error(c, 500, "获取用户列表失败")
		return
	}

	response.Success(c, result)
}

// GetUserByID 根据ID获取用户（管理员功能）
// @Summary 根据ID获取用户
// @Description 根据ID获取用户信息，仅管理员可访问
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 404 {object} response.Response
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "用户ID格式错误")
		return
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		logger.Errorf("获取用户失败: %v", err)
		response.Error(c, 404, err.Error())
		return
	}

	response.Success(c, user)
}

// CreateUser 创建用户（管理员功能）
// @Summary 创建用户
// @Description 创建新用户，仅管理员可访问
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body model.UserCreateRequest true "用户信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
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

	user, err := h.userService.CreateUser(&req)
	if err != nil {
		logger.Errorf("创建用户失败: %v", err)
		response.Error(c, 400, err.Error())
		return
	}

	response.Success(c, user)
}

// UpdateUser 更新用户（管理员功能）
// @Summary 更新用户
// @Description 更新用户信息，仅管理员可访问
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body model.UserUpdateRequest true "更新信息"
// @Success 200 {object} response.Response{data=model.UserResponse}
// @Failure 400 {object} response.Response
// @Router /api/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "用户ID格式错误")
		return
	}

	var req model.UserUpdateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		logger.Warnf("参数绑定失败: %v", err)
		response.Error(c, 400, "参数格式错误")
		return
	}

	// 参数验证
	err = validator.ValidateStruct(&req)
	if err != nil {
		logger.Warnf("参数验证失败: %v", err)
		response.Error(c, 400, err.Error())
		return
	}

	user, err := h.userService.UpdateUser(uint(id), &req)
	if err != nil {
		logger.Errorf("更新用户失败: %v", err)
		if strings.Contains(err.Error(), "用户不存在") {
			response.Error(c, 404, err.Error())
		} else {
			response.Error(c, 400, err.Error())
		}
		return
	}

	response.Success(c, user)
}

// DeleteUser 删除用户（管理员功能）
// @Summary 删除用户
// @Description 删除用户，仅管理员可访问
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "用户ID格式错误")
		return
	}

	if err := h.userService.DeleteUser(uint(id)); err != nil {
		logger.Errorf("删除用户失败: %v", err)
		if strings.Contains(err.Error(), "用户不存在") {
			response.Error(c, 404, err.Error())
		} else {
			response.Error(c, 400, err.Error())
		}
		return
	}

	response.Success(c, gin.H{"message": "删除成功"})
}

// UpdateUserStatus 更新用户状态（管理员功能）
// @Summary 更新用户状态
// @Description 启用或禁用用户，仅管理员可访问
// @Tags 用户管理
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "用户ID"
// @Param request body map[string]int true "状态信息" example({"status": 1})
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/users/{id}/status [put]
func (h *UserHandler) UpdateUserStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, 400, "用户ID格式错误")
		return
	}

	var req struct {
		Status int `json:"status" validate:"oneof=0 1"`
	}
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

	if err := h.userService.UpdateUserStatus(uint(id), req.Status); err != nil {
		logger.Errorf("更新用户状态失败: %v", err)
		if strings.Contains(err.Error(), "用户不存在") {
			response.Error(c, 404, err.Error())
		} else {
			response.Error(c, 400, err.Error())
		}
		return
	}

	statusText := "启用"
	if req.Status == 0 {
		statusText = "禁用"
	}
	response.Success(c, gin.H{"message": statusText + "成功"})
}