package migration

import (
	"domain-admin/model"
	"domain-admin/pkg/logger"

	"gorm.io/gorm"
)

// InitRBACData 初始化RBAC基础数据
func InitRBACData(db *gorm.DB) error {
	logger.Info("开始初始化RBAC基础数据...")

	// 初始化角色
	if err := initRoles(db); err != nil {
		return err
	}

	// 初始化权限
	if err := initPermissions(db); err != nil {
		return err
	}

	// 初始化角色权限关联
	if err := initRolePermissions(db); err != nil {
		return err
	}

	logger.Info("RBAC基础数据初始化完成")
	return nil
}

// initRoles 初始化角色数据
func initRoles(db *gorm.DB) error {
	roles := []model.Role{
		{
			Name:        "admin",
			DisplayName: "系统管理员",
			Description: "拥有系统所有权限的超级管理员",
			Status:      1,
		},
		{
			Name:        "user",
			DisplayName: "普通用户",
			Description: "系统普通用户，拥有基础功能权限",
			Status:      1,
		},
		{
			Name:        "guest",
			DisplayName: "访客用户",
			Description: "只读权限的访客用户",
			Status:      1,
		},
	}

	for _, role := range roles {
		var existingRole model.Role
		result := db.Where("name = ?", role.Name).First(&existingRole)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&role).Error; err != nil {
				logger.Errorf("创建角色失败: %s, error: %v", role.Name, err)
				return err
			}
			logger.Infof("创建角色成功: %s", role.Name)
		} else if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// initPermissions 初始化权限数据
func initPermissions(db *gorm.DB) error {
	permissions := []model.Permission{
		// 用户管理权限
		{Name: "user.list", DisplayName: "查看用户列表", Description: "查看系统用户列表", Resource: "/api/users", Action: "GET", Status: 1},
		{Name: "user.create", DisplayName: "创建用户", Description: "创建新用户", Resource: "/api/users", Action: "POST", Status: 1},
		{Name: "user.update", DisplayName: "更新用户", Description: "更新用户信息", Resource: "/api/users/*", Action: "PUT", Status: 1},
		{Name: "user.delete", DisplayName: "删除用户", Description: "删除用户", Resource: "/api/users/*", Action: "DELETE", Status: 1},
		{Name: "user.detail", DisplayName: "查看用户详情", Description: "查看用户详细信息", Resource: "/api/users/*", Action: "GET", Status: 1},

		// 角色管理权限
		{Name: "role.list", DisplayName: "查看角色列表", Description: "查看系统角色列表", Resource: "/api/roles", Action: "GET", Status: 1},
		{Name: "role.create", DisplayName: "创建角色", Description: "创建新角色", Resource: "/api/roles", Action: "POST", Status: 1},
		{Name: "role.update", DisplayName: "更新角色", Description: "更新角色信息", Resource: "/api/roles/*", Action: "PUT", Status: 1},
		{Name: "role.delete", DisplayName: "删除角色", Description: "删除角色", Resource: "/api/roles/*", Action: "DELETE", Status: 1},
		{Name: "role.detail", DisplayName: "查看角色详情", Description: "查看角色详细信息", Resource: "/api/roles/*", Action: "GET", Status: 1},

		// 权限管理权限
		{Name: "permission.list", DisplayName: "查看权限列表", Description: "查看系统权限列表", Resource: "/api/permissions", Action: "GET", Status: 1},
		{Name: "permission.create", DisplayName: "创建权限", Description: "创建新权限", Resource: "/api/permissions", Action: "POST", Status: 1},
		{Name: "permission.update", DisplayName: "更新权限", Description: "更新权限信息", Resource: "/api/permissions/*", Action: "PUT", Status: 1},
		{Name: "permission.delete", DisplayName: "删除权限", Description: "删除权限", Resource: "/api/permissions/*", Action: "DELETE", Status: 1},
		{Name: "permission.detail", DisplayName: "查看权限详情", Description: "查看权限详细信息", Resource: "/api/permissions/*", Action: "GET", Status: 1},

		// 认证相关权限
		{Name: "auth.login", DisplayName: "用户登录", Description: "用户登录系统", Resource: "/api/auth/login", Action: "POST", Status: 1},
		{Name: "auth.logout", DisplayName: "用户登出", Description: "用户登出系统", Resource: "/api/auth/logout", Action: "POST", Status: 1},
		{Name: "auth.profile", DisplayName: "查看个人信息", Description: "查看个人资料信息", Resource: "/api/auth/profile", Action: "GET", Status: 1},
		{Name: "auth.update_profile", DisplayName: "更新个人信息", Description: "更新个人资料信息", Resource: "/api/auth/profile", Action: "PUT", Status: 1},
		{Name: "auth.change_password", DisplayName: "修改密码", Description: "修改用户密码", Resource: "/api/auth/password", Action: "PUT", Status: 1},

		// 系统管理权限
		{Name: "system.all", DisplayName: "系统全部权限", Description: "系统所有功能的访问权限", Resource: "/api/*", Action: "*", Status: 1},
	}

	for _, permission := range permissions {
		var existingPermission model.Permission
		result := db.Where("name = ?", permission.Name).First(&existingPermission)
		if result.Error == gorm.ErrRecordNotFound {
			if err := db.Create(&permission).Error; err != nil {
				logger.Errorf("创建权限失败: %s, error: %v", permission.Name, err)
				return err
			}
			logger.Infof("创建权限成功: %s", permission.Name)
		} else if result.Error != nil {
			return result.Error
		}
	}

	return nil
}

// initRolePermissions 初始化角色权限关联
func initRolePermissions(db *gorm.DB) error {
	// 管理员角色拥有所有权限
	var adminRole model.Role
	if err := db.Preload("Permissions").Where("name = ?", "admin").First(&adminRole).Error; err != nil {
		return err
	}

	var systemAllPermission model.Permission
	if err := db.Where("name = ?", "system.all").First(&systemAllPermission).Error; err != nil {
		return err
	}

	// 检查是否已存在关联
	hasPermission := false
	for _, perm := range adminRole.Permissions {
		if perm.ID == systemAllPermission.ID {
			hasPermission = true
			break
		}
	}

	if !hasPermission {
		if err := db.Model(&adminRole).Association("Permissions").Append(&systemAllPermission); err != nil {
			logger.Errorf("创建管理员角色权限关联失败: %v", err)
			return err
		}
		logger.Info("创建管理员角色权限关联成功")
	}

	// 普通用户角色的基础权限
	var userRole model.Role
	if err := db.Preload("Permissions").Where("name = ?", "user").First(&userRole).Error; err != nil {
		return err
	}

	userPermissions := []string{"auth.login", "auth.logout", "auth.profile", "auth.update_profile", "auth.change_password"}
	var userPermsToAdd []model.Permission
	for _, permName := range userPermissions {
		var permission model.Permission
		if err := db.Where("name = ?", permName).First(&permission).Error; err != nil {
			continue
		}

		// 检查是否已存在关联
		hasPermission := false
		for _, perm := range userRole.Permissions {
			if perm.ID == permission.ID {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			userPermsToAdd = append(userPermsToAdd, permission)
		}
	}

	if len(userPermsToAdd) > 0 {
		if err := db.Model(&userRole).Association("Permissions").Append(&userPermsToAdd); err != nil {
			logger.Errorf("创建用户角色权限关联失败: %v", err)
			return err
		}
		logger.Infof("创建用户角色权限关联成功，添加了 %d 个权限", len(userPermsToAdd))
	}

	// 访客角色的只读权限
	var guestRole model.Role
	if err := db.Preload("Permissions").Where("name = ?", "guest").First(&guestRole).Error; err != nil {
		return err
	}

	guestPermissions := []string{"auth.login", "auth.logout", "auth.profile"}
	var guestPermsToAdd []model.Permission
	for _, permName := range guestPermissions {
		var permission model.Permission
		if err := db.Where("name = ?", permName).First(&permission).Error; err != nil {
			continue
		}

		// 检查是否已存在关联
		hasPermission := false
		for _, perm := range guestRole.Permissions {
			if perm.ID == permission.ID {
				hasPermission = true
				break
			}
		}

		if !hasPermission {
			guestPermsToAdd = append(guestPermsToAdd, permission)
		}
	}

	if len(guestPermsToAdd) > 0 {
		if err := db.Model(&guestRole).Association("Permissions").Append(&guestPermsToAdd); err != nil {
			logger.Errorf("创建访客角色权限关联失败: %v", err)
			return err
		}
		logger.Infof("创建访客角色权限关联成功，添加了 %d 个权限", len(guestPermsToAdd))
	}

	return nil
}