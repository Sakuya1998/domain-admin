package data

import (
	"github.com/Sakuya1998/domain-admin/app/user/internal/conf"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ProviderSet is data providers.
var ProviderSet = wire.NewSet(NewData, NewUserRepo)

// Data .
type Data struct {
	db     *gorm.DB
	logger log.Logger
}

// NewData .
func NewData(c *conf.Data, l log.Logger) (*Data, func(), error) {
	helper := log.NewHelper(l)

	// 连接数据库
	db, err := gorm.Open(mysql.Open(c.Database.Source), &gorm.Config{})
	if err != nil {
		helper.Errorf("failed opening connection to mysql: %v", err)
		return nil, nil, err
	}

	// 自动迁移User模型
	if err := db.AutoMigrate(&User{}); err != nil {
		helper.Errorf("failed to migrate database: %v", err)
		return nil, nil, err
	}

	cleanup := func() {
		helper.Info("closing the data resources")
		if sqlDB, err := db.DB(); err == nil {
			sqlDB.Close()
		}
	}

	return &Data{db: db, logger: l}, cleanup, nil
}
