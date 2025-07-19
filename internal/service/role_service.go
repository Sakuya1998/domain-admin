package service

import (
	"domain-admin/internal/repository"
	"domain-admin/model"
	"domain-admin/pkg/pagination"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// RoleService 角色服务接口
type RoleService interface {
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

type roleService struct {
	roleRepo       repository.RoleRepository
	permissionRepo repository.PermissionRepository
}

// NewRoleService 创建角色服务实例
func NewRoleService(roleRepo repository.RoleRepository, permissionRepo repository.PermissionRepository) RoleService {
	return &roleService{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

// Create 创建角色
func (s *roleService) Create(role *model.Role) error {
	if role.Name == "" {
		return errors.New("角色名称不能为空")
	}

	if role.DisplayName == "" {
		return errors.New("角色显示名称不能为空")
	}

	// 检查角色名称是否已存在
	existingRole, err := s.roleRepo.GetByName(role.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查角色名称失败: %w", err)
	}
	if existingRole != nil {
		return errors.New("角色名称已存在")
	}

	return s.roleRepo.Create(role)
}

// GetByID 根据ID获取角色
func (s *roleService) GetByID(id uint) (*model.Role, error) {
	if id == 0 {
		return nil, errors.New("角色ID不能为空")
	}
	return s.roleRepo.GetByID(id)
}

// GetByName 根据名称获取角色
func (s *roleService) GetByName(name string) (*model.Role, error) {
	if name == "" {
		return nil, errors.New("角色名称不能为空")
	}
	return s.roleRepo.GetByName(name)
}

// Update 更新角色
func (s *roleService) Update(role *model.Role) error {
	if role.ID == 0 {
		return errors.New("角色ID不能为空")
	}

	if role.Name == "" {
		return errors.New("角色名称不能为空")
	}

	if role.DisplayName == "" {
		return errors.New("角色显示名称不能为空")
	}

	// 检查角色是否存在
	existingRole, err := s.roleRepo.GetByID(role.ID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 如果角色名称发生变化，检查新名称是否已存在
	if existingRole.Name != role.Name {
		conflictRole, err := s.roleRepo.GetByName(role.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("检查角色名称失败: %w", err)
		}
		if conflictRole != nil {
			return errors.New("角色名称已存在")
		}
	}

	return s.roleRepo.Update(role)
}

// Delete 删除角色
func (s *roleService) Delete(id uint) error {
	if id == 0 {
		return errors.New("角色ID不能为空")
	}

	// 检查角色是否存在
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 检查是否为系统内置角色，不允许删除
	if role.Name == "admin" || role.Name == "user" || role.Name == "guest" {
		return errors.New("系统内置角色不允许删除")
	}

	return s.roleRepo.Delete(id)
}

// List 获取角色列表
func (s *roleService) List(page pagination.Pagination) ([]*model.Role, int64, error) {
	return s.roleRepo.List(page)
}

// UpdateStatus 更新角色状态
func (s *roleService) UpdateStatus(id uint, status int) error {
	if id == 0 {
		return errors.New("角色ID不能为空")
	}

	if status != 0 && status != 1 {
		return errors.New("状态值无效")
	}

	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	return s.roleRepo.UpdateStatus(id, status)
}

// AssignPermissions 为角色分配权限
func (s *roleService) AssignPermissions(roleID uint, permissionIDs []uint) error {
	if roleID == 0 {
		return errors.New("角色ID不能为空")
	}

	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return fmt.Errorf("角色不存在: %w", err)
	}

	// 检查权限是否存在
	for _, permissionID := range permissionIDs {
		_, err := s.permissionRepo.GetByID(permissionID)
		if err != nil {
			return fmt.Errorf("权限ID %d 不存在: %w", permissionID, err)
		}
	}

	return s.roleRepo.AssignPermissions(roleID, permissionIDs)
}

// GetRolePermissions 获取角色的权限列表
func (s *roleService) GetRolePermissions(roleID uint) ([]*model.Permission, error) {
	if roleID == 0 {
		return nil, errors.New("角色ID不能为空")
	}

	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("角色不存在: %w", err)
	}

	return s.roleRepo.GetRolePermissions(roleID)
}