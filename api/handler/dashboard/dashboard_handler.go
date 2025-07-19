package dashboard

import (
	"domain-admin/internal/repository"
	"domain-admin/pkg/db"
	"domain-admin/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	userRepo       repository.UserRepository
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewDashboardHandler 创建仪表盘处理器
func NewDashboardHandler() *DashboardHandler {
	db := db.GetDB("default")
	return &DashboardHandler{
		userRepo:       repository.NewUserRepository(db),
		roleRepo:       repository.NewRoleRepository(db),
		permissionRepo: repository.NewPermissionRepository(db),
	}
}

// DashboardStats 仪表盘统计数据结构
type DashboardStats struct {
	UserCount       int64 `json:"userCount"`
	RoleCount       int64 `json:"roleCount"`
	PermissionCount int64 `json:"permissionCount"`
	OnlineCount     int64 `json:"onlineCount"`
}

// GetStats 获取仪表盘统计数据
func (h *DashboardHandler) GetStats(c *gin.Context) {
	stats := &DashboardStats{}

	// 获取用户总数
	userCount, err := h.userRepo.Count()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取用户统计失败")
		return
	}
	stats.UserCount = userCount

	// 获取角色总数
	roleCount, err := h.roleRepo.Count()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取角色统计失败")
		return
	}
	stats.RoleCount = roleCount

	// 获取权限总数
	permissionCount, err := h.permissionRepo.Count()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "获取权限统计失败")
		return
	}
	stats.PermissionCount = permissionCount

	// 在线用户数（暂时设为固定值，后续可以通过Redis或其他方式实现）
	stats.OnlineCount = 1 // 当前登录用户至少为1

	response.Success(c, stats)
}
