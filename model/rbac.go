package model

import (
	"time"

	"gorm.io/gorm"
)

// Role 角色模型
type Role struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:50;not null;comment:角色名称"`
	DisplayName string         `json:"display_name" gorm:"size:100;not null;comment:角色显示名称"`
	Description string         `json:"description" gorm:"size:255;comment:角色描述"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1:启用 0:禁用"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Permissions []Permission `json:"permissions" gorm:"many2many:role_permissions;"`
	Users       []User       `json:"users" gorm:"foreignKey:Role;references:Name"`
}

// Permission 权限模型
type Permission struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	ParentID    uint           `json:"parent_id" gorm:"default:0;comment:父权限ID"`
	Name        string         `json:"name" gorm:"uniqueIndex;size:100;not null;comment:权限名称"`
	DisplayName string         `json:"display_name" gorm:"size:100;not null;comment:权限显示名称"`
	Description string         `json:"description" gorm:"size:255;comment:权限描述"`
	Resource    string         `json:"resource" gorm:"size:100;not null;comment:资源路径"`
	Action      string         `json:"action" gorm:"size:20;not null;comment:操作类型(GET,POST,PUT,DELETE,*)"`
	Type        string         `json:"type" gorm:"size:20;default:menu;comment:权限类型(menu,button,api)"`
	Sort        int            `json:"sort" gorm:"default:0;comment:排序"`
	Status      int            `json:"status" gorm:"default:1;comment:状态 1:启用 0:禁用"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`

	// 关联关系
	Roles []Role `json:"roles" gorm:"many2many:role_permissions;"`
}