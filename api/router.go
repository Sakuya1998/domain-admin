package api

import (
	"domain-admin/api/handler/auth"
	"domain-admin/api/handler/permission"
	"domain-admin/api/handler/role"
	"domain-admin/api/handler/user"
	"domain-admin/pkg/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// 创建处理器
	authHandler := auth.NewAuthHandler()
	userHandler := user.NewUserHandler()
	roleHandler := role.NewRoleHandler()
	permissionHandler := permission.NewPermissionHandler()

	// API 路由组
	api := r.Group("/api")
	{
		// 认证相关路由（无需登录）
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			// 登出需要认证
			auth.POST("/logout", middleware.JWTAuth(), authHandler.Logout)
			// 个人资料相关（需要认证）
			auth.GET("/profile", middleware.JWTAuth(), authHandler.GetProfile)
			auth.PUT("/profile", middleware.JWTAuth(), authHandler.UpdateProfile)
			auth.PUT("/password", middleware.JWTAuth(), authHandler.UpdateProfile) // 修改密码接口
		}

		// 用户管理路由（需要认证和权限）
		users := api.Group("/users")
		users.Use(middleware.JWTAuth(), middleware.RBACMiddleware())
		{
			users.GET("", userHandler.GetUserList)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.PUT("/:id/status", userHandler.UpdateUserStatus)
		}

		// 角色管理路由（需要认证和权限）
		roles := api.Group("/roles")
		roles.Use(middleware.JWTAuth(), middleware.RBACMiddleware())
		{
			roles.GET("", roleHandler.ListRoles)
			roles.GET("/:id", roleHandler.GetRole)
			roles.POST("", roleHandler.CreateRole)
			roles.PUT("/:id", roleHandler.UpdateRole)
			roles.DELETE("/:id", roleHandler.DeleteRole)
			roles.PUT("/:id/status", roleHandler.UpdateRoleStatus)
			roles.GET("/:id/permissions", roleHandler.GetRolePermissions)
			roles.POST("/:id/permissions", roleHandler.AssignPermissions)
		}

		// 权限管理路由（需要认证和权限）
		permissions := api.Group("/permissions")
		permissions.Use(middleware.JWTAuth(), middleware.RBACMiddleware())
		{
			permissions.GET("", permissionHandler.ListPermissions)
			permissions.GET("/:id", permissionHandler.GetPermission)
			permissions.POST("", permissionHandler.CreatePermission)
			permissions.PUT("/:id", permissionHandler.UpdatePermission)
			permissions.DELETE("/:id", permissionHandler.DeletePermission)
			permissions.PUT("/:id/status", permissionHandler.UpdatePermissionStatus)
		}
	}
}
