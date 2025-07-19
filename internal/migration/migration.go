package migration

import (
	"domain-admin/model"
	"domain-admin/pkg/logger"

	"gorm.io/gorm"
)

// AutoMigrate 自动迁移数据库表结构
func AutoMigrate(db *gorm.DB) error {
	logger.Info("开始数据库迁移...")

	// 迁移用户表
	if err := db.AutoMigrate(&model.User{}); err != nil {
		logger.Errorf("用户表迁移失败: %v", err)
		return err
	}

	// 迁移角色表
	if err := db.AutoMigrate(&model.Role{}); err != nil {
		logger.Errorf("角色表迁移失败: %v", err)
		return err
	}

	// 迁移权限表
	if err := db.AutoMigrate(&model.Permission{}); err != nil {
		logger.Errorf("权限表迁移失败: %v", err)
		return err
	}



	logger.Info("数据库迁移完成")
	return nil
}

// CreateDefaultAdmin 创建默认管理员账户
func CreateDefaultAdmin(db *gorm.DB) error {
	// 检查是否已存在管理员账户
	var count int64
	db.Model(&model.User{}).Where("role = ?", "admin").Count(&count)
	if count > 0 {
		logger.Info("管理员账户已存在，跳过创建")
		return nil
	}

	// 创建默认管理员
	admin := &model.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // password
		Nickname: "系统管理员",
		Role:     "admin",
		Status:   1,
	}

	if err := db.Create(admin).Error; err != nil {
		logger.Errorf("创建默认管理员失败: %v", err)
		return err
	}

	logger.Info("默认管理员创建成功 - 用户名: admin, 密码: password")
	return nil
}