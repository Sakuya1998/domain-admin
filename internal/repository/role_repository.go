package repository

import (
	"domain-admin/model"
	"domain-admin/pkg/pagination"
	"errors"

	"gorm.io/gorm"
)

// RoleRepository 角色仓储接口
type RoleRepository interface {
	Create(role *model.Role) error
	GetByID(id uint) (*model.Role, error)
	GetByName(name string) (*model.Role, error)
	Update(role *model.Role) error
	Delete(id uint) error
	List(page pagination.Pagination) ([]*model.Role, int64, error)
	UpdateStatus(id uint, status int) error
	AssignPermissions(roleID uint, permissionIDs []uint) error
	GetRolePermissions(roleID uint) ([]*model.Permission, error)
}

type roleRepository struct {
	db *gorm.DB
}

// NewRoleRepository 创建角色仓储实例
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

// Create 创建角色
func (r *roleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

// GetByID 根据ID获取角色
func (r *roleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("id = ?", id).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// GetByName 根据名称获取角色
func (r *roleRepository) GetByName(name string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("角色不存在")
		}
		return nil, err
	}
	return &role, nil
}

// Update 更新角色
func (r *roleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

// Delete 删除角色
func (r *roleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

// List 获取角色列表
func (r *roleRepository) List(page pagination.Pagination) ([]*model.Role, int64, error) {
	var roles []*model.Role
	var total int64

	query := r.db.Model(&model.Role{})

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	if err := query.Offset(page.Offset).Limit(page.Limit).Order(page.GetOrderClause()).Find(&roles).Error; err != nil {
		return nil, 0, err
	}

	return roles, total, nil
}

// UpdateStatus 更新角色状态
func (r *roleRepository) UpdateStatus(id uint, status int) error {
	return r.db.Model(&model.Role{}).Where("id = ?", id).Update("status", status).Error
}

// AssignPermissions 为角色分配权限
func (r *roleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 获取角色
		var role model.Role
		if err := tx.First(&role, roleID).Error; err != nil {
			return err
		}

		// 获取权限列表
		var permissions []model.Permission
		if len(permissionIDs) > 0 {
			if err := tx.Find(&permissions, permissionIDs).Error; err != nil {
				return err
			}
		}

		// 使用GORM的Association替换权限
		if err := tx.Model(&role).Association("Permissions").Replace(&permissions); err != nil {
			return err
		}

		return nil
	})
}

// GetRolePermissions 获取角色的权限列表
func (r *roleRepository) GetRolePermissions(roleID uint) ([]*model.Permission, error) {
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