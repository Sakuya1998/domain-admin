package repository

import (
	"domain-admin/model"
	"domain-admin/pkg/pagination"
	"errors"

	"gorm.io/gorm"
)

// PermissionRepository 权限仓储接口
type PermissionRepository interface {
	Create(permission *model.Permission) error
	GetByID(id uint) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	Update(permission *model.Permission) error
	Delete(id uint) error
	List(page pagination.Pagination) ([]*model.Permission, int64, error)
	UpdateStatus(id uint, status int) error
	GetPermissionsByRole(roleID uint) ([]*model.Permission, error)
}

type permissionRepository struct {
	db *gorm.DB
}

// NewPermissionRepository 创建权限仓储实例
func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

// Create 创建权限
func (r *permissionRepository) Create(permission *model.Permission) error {
	return r.db.Create(permission).Error
}

// GetByID 根据ID获取权限
func (r *permissionRepository) GetByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("id = ?", id).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("权限不存在")
		}
		return nil, err
	}
	return &permission, nil
}

// GetByName 根据名称获取权限
func (r *permissionRepository) GetByName(name string) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.Where("name = ?", name).First(&permission).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("权限不存在")
		}
		return nil, err
	}
	return &permission, nil
}

// Update 更新权限
func (r *permissionRepository) Update(permission *model.Permission) error {
	return r.db.Save(permission).Error
}

// Delete 删除权限
func (r *permissionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Permission{}, id).Error
}

// List 获取权限列表
func (r *permissionRepository) List(page pagination.Pagination) ([]*model.Permission, int64, error) {
	var permissions []*model.Permission
	var total int64

	query := r.db.Model(&model.Permission{})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(page.Offset).Limit(page.Limit).Order(page.GetOrderClause()).Find(&permissions).Error; err != nil {
		return nil, 0, err
	}

	return permissions, total, nil
}

// UpdateStatus 更新权限状态
func (r *permissionRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Permission{}).Where("id = ?", id).Update("status", status).Error
}

// GetPermissionsByRole 根据角色ID获取权限列表
func (r *permissionRepository) GetPermissionsByRole(roleID uint) ([]*model.Permission, error) {
	var permissions []*model.Permission

	err := r.db.Table("domain_permission").
		Joins("JOIN domain_role_permissions ON domain_permission.id = domain_role_permissions.permission_id").
		Where("domain_role_permissions.role_id = ? AND domain_permission.status = ?", roleID, 1).
		Find(&permissions).Error

	if err != nil {
		return nil, err
	}

	return permissions, nil
}