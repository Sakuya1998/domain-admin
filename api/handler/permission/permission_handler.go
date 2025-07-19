package permission

import (
	"domain-admin/internal/repository"
	"domain-admin/internal/service"
	"domain-admin/model"
	"domain-admin/pkg/db"
	"domain-admin/pkg/pagination"
	"domain-admin/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PermissionHandler struct {
	permissionService service.PermissionService
}

// NewPermissionHandler 创建权限处理器
func NewPermissionHandler() *PermissionHandler {
	permissionRepo := repository.NewPermissionRepository(db.GetDB("default"))
	roleRepo := repository.NewRoleRepository(db.GetDB("default"))
	permissionService := service.NewPermissionService(permissionRepo, roleRepo)

	return &PermissionHandler{
		permissionService: permissionService,
	}
}

// CreatePermission 创建权限
// @Summary 创建权限
// @Description 创建新的权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param permission body model.Permission true "权限信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/permissions [post]
func (h *PermissionHandler) CreatePermission(c *gin.Context) {
	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.permissionService.Create(&permission); err != nil {
		response.Error(c, http.StatusBadRequest, "创建权限失败")
		return
	}

	response.Success(c, permission)
}

// GetPermission 获取权限详情
// @Summary 获取权限详情
// @Description 根据ID获取权限详情
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/permissions/{id} [get]
func (h *PermissionHandler) GetPermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	permission, err := h.permissionService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "权限不存在")
		return
	}

	response.Success(c, permission)
}

// UpdatePermission 更新权限
// @Summary 更新权限
// @Description 更新权限信息
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Param permission body model.Permission true "权限信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/permissions/{id} [put]
func (h *PermissionHandler) UpdatePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	var permission model.Permission
	if err := c.ShouldBindJSON(&permission); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	permission.ID = uint(id)
	if err := h.permissionService.Update(&permission); err != nil {
		response.Error(c, http.StatusBadRequest, "更新权限失败")
		return
	}

	response.Success(c, permission)
}

// DeletePermission 删除权限
// @Summary 删除权限
// @Description 删除权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/permissions/{id} [delete]
func (h *PermissionHandler) DeletePermission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	if err := h.permissionService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, "删除权限失败")
		return
	}

	response.Success(c, nil)
}

// ListPermissions 获取权限列表
// @Summary 获取权限列表
// @Description 分页获取权限列表
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} response.Response{data=pagination.PageResult}
// @Failure 400 {object} response.Response
// @Router /api/permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	page := pagination.New(c)

	permissions, total, err := h.permissionService.List(page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取权限列表失败")
		return
	}

	result := pagination.NewPageResult(total, permissions)
	response.Success(c, result)
}

// UpdatePermissionStatus 更新权限状态
// @Summary 更新权限状态
// @Description 启用或禁用权限
// @Tags 权限管理
// @Accept json
// @Produce json
// @Param id path int true "权限ID"
// @Param status body map[string]int true "状态信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/permissions/{id}/status [put]
func (h *PermissionHandler) UpdatePermissionStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的权限ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.permissionService.UpdateStatus(uint(id), req.Status); err != nil {
		response.Error(c, http.StatusBadRequest, "更新权限状态失败")
		return
	}

	response.Success(c, nil)
}