package middleware

import (
	"fmt"
	"net/http"
	"sync"

	"domain-admin/pkg/logger"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var (
	Enforcer    *casbin.Enforcer
	enforcerMux sync.RWMutex
	initialized bool
)

// InitRBAC 初始化基于Gorm Adapter的Casbin RBAC
// 传入 gorm.DB 和模型配置文件路径（rbac_model.conf）
func InitRBAC(db *gorm.DB, modelPath string) error {
	enforcerMux.Lock()
	defer enforcerMux.Unlock()

	if modelPath == "" {
		modelPath = "configs/rbac_model.conf"
	}

	// 创建 Gorm 适配器，使用传入的 gorm.DB 连接
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		logger.Errorf("Failed to create Gorm adapter for Casbin", "error", err)
		return fmt.Errorf("failed to create adapter: %w", err)
	}

	// 创建 enforcer，加载模型和数据库策略
	e, err := casbin.NewEnforcer(modelPath, adapter)
	if err != nil {
		logger.Errorf("Failed to create Casbin enforcer", "error", err)
		return fmt.Errorf("failed to create enforcer: %w", err)
	}

	// 从数据库加载策略
	if loadErr := e.LoadPolicy(); loadErr != nil {
		logger.Errorf("Failed to load policy from database", "error", err)
		return fmt.Errorf("failed to load policy: %w", err)
	}

	Enforcer = e
	initialized = true
	policies, err := e.GetPolicy()
	if err != nil {
		logger.Errorf("Failed to get policies", "error", err)
		return fmt.Errorf("failed to get policies: %w", err)
	}
	policiesCount := len(policies)

	logger.Infof("RBAC system initialized with Gorm adapter, model=%s, policies_count=%d", modelPath, policiesCount)

	return nil
}

// RBACMiddleware Gin中间件：权限校验
func RBACMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !initialized || Enforcer == nil {
			logger.Error("RBAC system not initialized")
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "RBAC system not initialized",
			})
			return
		}

		roleInterface, exists := c.Get("role")
		if !exists {
			logger.Warnf("User role not found in context", "path", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "User role not found",
			})
			return
		}

		role, ok := roleInterface.(string)
		if !ok {
			logger.Error("Invalid role type in context")
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Invalid role format",
			})
			return
		}

		path := c.Request.URL.Path
		method := c.Request.Method

		enforcerMux.RLock()
		allowed, err := Enforcer.Enforce(role, path, method)
		enforcerMux.RUnlock()

		if err != nil {
			logger.Errorf("RBAC enforcement error", "error", err, "role", role, "path", path, "method", method)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": "Access control error",
			})
			return
		}

		if !allowed {
			logger.Warnf("Access denied", "role", role, "path", path, "method", method)
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"error": "Access denied",
			})
			return
		}

		c.Next()
	}
}

// SyncRBACPolicies 同步RBAC策略到Casbin
func SyncRBACPolicies(db *gorm.DB) error {
	if !initialized || Enforcer == nil {
		return fmt.Errorf("RBAC system not initialized")
	}

	logger.Info("开始同步RBAC策略...")

	// 清除现有策略
	enforcerMux.Lock()
	Enforcer.ClearPolicy()
	enforcerMux.Unlock()

	// 同步角色权限策略
	if err := syncRolePermissionPolicies(db); err != nil {
		return fmt.Errorf("同步角色权限策略失败: %w", err)
	}

	// 保存策略到数据库
	enforcerMux.Lock()
	err := Enforcer.SavePolicy()
	enforcerMux.Unlock()

	if err != nil {
		return fmt.Errorf("保存策略失败: %w", err)
	}

	logger.Info("RBAC策略同步完成")
	return nil
}

// syncRolePermissionPolicies 同步角色权限策略
func syncRolePermissionPolicies(db *gorm.DB) error {
	// 获取所有启用的角色权限关联
	var rolePermissions []struct {
		RoleName       string `gorm:"column:role_name"`
		PermissionName string `gorm:"column:permission_name"`
		Resource       string `gorm:"column:resource"`
		Action         string `gorm:"column:action"`
	}

	err := db.Table("domain_role_permissions rp").
		Joins("JOIN domain_role r ON rp.role_id = r.id").
		Joins("JOIN domain_permission p ON rp.permission_id = p.id").
		Where("r.status = ? AND p.status = ?", 1, 1).
		Select("r.name as role_name, p.name as permission_name, p.resource, p.action").
		Find(&rolePermissions).Error

	if err != nil {
		return fmt.Errorf("查询角色权限关联失败: %w", err)
	}

	// 添加策略到Casbin
	enforcerMux.Lock()
	defer enforcerMux.Unlock()

	for _, rp := range rolePermissions {
		// 添加策略: p, role, resource, action
		success, err := Enforcer.AddPolicy(rp.RoleName, rp.Resource, rp.Action)
		if err != nil {
			logger.Errorf("添加策略失败", "role", rp.RoleName, "resource", rp.Resource, "action", rp.Action, "error", err)
			continue
		}
		if success {
			logger.Debugf("添加策略成功: %s, %s, %s", rp.RoleName, rp.Resource, rp.Action)
		}
	}

	logger.Infof("同步了 %d 条角色权限策略", len(rolePermissions))
	return nil
}

// AddRolePermissionPolicy 添加角色权限策略
func AddRolePermissionPolicy(role, resource, action string) error {
	if !initialized || Enforcer == nil {
		return fmt.Errorf("RBAC system not initialized")
	}

	enforcerMux.Lock()
	defer enforcerMux.Unlock()

	success, err := Enforcer.AddPolicy(role, resource, action)
	if err != nil {
		return fmt.Errorf("添加策略失败: %w", err)
	}

	if success {
		logger.Infof("添加策略成功: %s, %s, %s", role, resource, action)
		// 保存策略
		if err := Enforcer.SavePolicy(); err != nil {
			logger.Errorf("保存策略失败", "error", err)
		}
	}

	return nil
}

// RemoveRolePermissionPolicy 移除角色权限策略
func RemoveRolePermissionPolicy(role, resource, action string) error {
	if !initialized || Enforcer == nil {
		return fmt.Errorf("RBAC system not initialized")
	}

	enforcerMux.Lock()
	defer enforcerMux.Unlock()

	success, err := Enforcer.RemovePolicy(role, resource, action)
	if err != nil {
		return fmt.Errorf("移除策略失败: %w", err)
	}

	if success {
		logger.Infof("移除策略成功: %s, %s, %s", role, resource, action)
		// 保存策略
		if err := Enforcer.SavePolicy(); err != nil {
			logger.Errorf("保存策略失败", "error", err)
		}
	}

	return nil
}

// GetAllPolicies 获取所有策略
func GetAllPolicies() ([][]string, error) {
	if !initialized || Enforcer == nil {
		return nil, fmt.Errorf("RBAC system not initialized")
	}

	enforcerMux.RLock()
	defer enforcerMux.RUnlock()

	policies, err := Enforcer.GetPolicy()
	return policies, err
}
