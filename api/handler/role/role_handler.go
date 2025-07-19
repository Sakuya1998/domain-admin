package role

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

type RoleHandler struct {
	roleService service.RoleService
}

// NewRoleHandler 创建角色处理器
func NewRoleHandler() *RoleHandler {
	roleRepo := repository.NewRoleRepository(db.GetDB("default"))
	permissionRepo := repository.NewPermissionRepository(db.GetDB("default"))
	roleService := service.NewRoleService(roleRepo, permissionRepo)

	return &RoleHandler{
		roleService: roleService,
	}
}

// CreateRole 创建角色
// @Summary 创建角色
// @Description 创建新的角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param role body model.Role true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles [post]
func (h *RoleHandler) CreateRole(c *gin.Context) {
	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.roleService.Create(&role); err != nil {
		response.Error(c, http.StatusBadRequest, "创建角色失败")
		return
	}

	response.Success(c, role)
}

// GetRole 获取角色详情
// @Summary 获取角色详情
// @Description 根据ID获取角色详情
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	role, err := h.roleService.GetByID(uint(id))
	if err != nil {
		response.Error(c, http.StatusNotFound, "角色不存在")
		return
	}

	response.Success(c, role)
}

// UpdateRole 更新角色
// @Summary 更新角色
// @Description 更新角色信息
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param role body model.Role true "角色信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id} [put]
func (h *RoleHandler) UpdateRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var role model.Role
	if err := c.ShouldBindJSON(&role); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	role.ID = uint(id)
	if err := h.roleService.Update(&role); err != nil {
		response.Error(c, http.StatusBadRequest, "更新角色失败")
		return
	}

	response.Success(c, role)
}

// DeleteRole 删除角色
// @Summary 删除角色
// @Description 删除角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id} [delete]
func (h *RoleHandler) DeleteRole(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	if err := h.roleService.Delete(uint(id)); err != nil {
		response.Error(c, http.StatusBadRequest, "删除角色失败")
		return
	}

	response.Success(c, nil)
}

// ListRoles 获取角色列表
// @Summary 获取角色列表
// @Description 分页获取角色列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(10)
// @Param search query string false "搜索关键词"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	page := pagination.New(c)

	roles, total, err := h.roleService.List(page)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色列表失败")
		return
	}

	response.SuccessWithPagination(c, roles, total, page)
}

// UpdateRoleStatus 更新角色状态
// @Summary 更新角色状态
// @Description 启用或禁用角色
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param status body map[string]int true "状态信息"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id}/status [put]
func (h *RoleHandler) UpdateRoleStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		Status int `json:"status" binding:"required,oneof=0 1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.roleService.UpdateStatus(uint(id), req.Status); err != nil {
		response.Error(c, http.StatusBadRequest, "更新角色状态失败")
		return
	}

	response.Success(c, nil)
}

// AssignPermissions 为角色分配权限
// @Summary 为角色分配权限
// @Description 为指定角色分配权限
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Param permissions body map[string][]uint true "权限ID列表"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id}/permissions [post]
func (h *RoleHandler) AssignPermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	var req struct {
		PermissionIDs []uint `json:"permission_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "参数错误")
		return
	}

	if err := h.roleService.AssignPermissions(uint(id), req.PermissionIDs); err != nil {
		response.Error(c, http.StatusBadRequest, "分配权限失败")
		return
	}

	response.Success(c, nil)
}

// GetRolePermissions 获取角色权限列表
// @Summary 获取角色权限列表
// @Description 获取指定角色的权限列表
// @Tags 角色管理
// @Accept json
// @Produce json
// @Param id path int true "角色ID"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /api/roles/{id}/permissions [get]
func (h *RoleHandler) GetRolePermissions(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "无效的角色ID")
		return
	}

	permissions, err := h.roleService.GetRolePermissions(uint(id))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "获取角色权限失败")
		return
	}

	response.Success(c, permissions)
}