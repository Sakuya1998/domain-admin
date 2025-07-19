package middleware

import (
	"context"
	"domain-admin/pkg/cache"
	"domain-admin/pkg/jwt"
	"domain-admin/pkg/logger"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 使用defer捕获panic，防止无效token导致500错误
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("JWT middleware panic: %v", r)
				c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
					"code":    401,
					"message": "Invalid or expired token",
				})
			}
		}()

		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warnf("Authorization header missing, path: %s", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Authorization header missing",
			})
			return
		}

		// 检查 Bearer 前缀
		const prefix = "Bearer "
		if !strings.HasPrefix(authHeader, prefix) {
			logger.Warnf("Invalid authorization header format, path: %s", c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid authorization header format",
			})
			return
		}

		token := strings.TrimPrefix(authHeader, prefix)
		if token == "" {
			logger.Warn("Empty token provided")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Token is empty",
			})
			return
		}

		claims, err := jwt.ParseToken(token)
		if err != nil {
			logger.Warnf("Token validation failed: %v, path: %s", err, c.Request.URL.Path)
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "Invalid or expired token",
			})
			return
		}

		// 验证会话缓存（如果Redis可用）
		ctx := context.Background()
		sessionKey := strconv.FormatUint(uint64(claims.UserID), 10)
		var sessionData map[string]interface{}
		userID := claims.UserID
		username := claims.Username
		role := claims.Role
		
		if err := cache.GetSessionCache(ctx, sessionKey, &sessionData); err != nil {
			if !errors.Is(err, redis.Nil) {
				logger.Warnf("获取会话缓存失败: %v", err)
				// Redis连接错误时，跳过缓存验证，继续使用JWT信息
			} else {
				// 会话缓存不存在，可能是Redis未启动或会话过期
				// 在测试环境下允许继续，生产环境可以要求重新登录
				logger.Warnf("会话已过期或不存在，用户ID: %d", claims.UserID)
			}
		} else {
			// 从会话缓存中获取最新的用户信息
			if cachedUsername, ok := sessionData["username"].(string); ok && cachedUsername != "" {
				username = cachedUsername
			}
			if cachedRole, ok := sessionData["role"].(string); ok && cachedRole != "" {
				role = cachedRole
			}
		}

		// 记录用户信息到日志
		logger.Debugf("userID: %d, role: %s, username: %s", userID, role, username)

		c.Set("userID", userID)
		c.Set("user_id", userID) // 兼容性
		c.Set("role", role)
		c.Set("username", username)
		c.Next()
	}
}

// AdminAuth 管理员权限中间件
func AdminAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未认证",
			})
			c.Abort()
			return
		}

		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"code":    403,
				"message": "权限不足，需要管理员权限",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}
