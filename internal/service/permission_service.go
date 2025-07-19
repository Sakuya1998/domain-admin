package service

import (
	"domain-admin/internal/repository"
	"domain-admin/model"
	"domain-admin/pkg/pagination"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// PermissionService 权限服务接口
type PermissionService interface {
	Create(permission *model.Permission) error
	GetByID(id uint) (*model.Permission, error)
	GetByName(name string) (*model.Permission, error)
	Update(permission *model.Permission) error
	Delete(id uint) error
	List(page pagination.Pagination) ([]*model.Permission, int64, error)
	UpdateStatus(id uint, status int) error
	GetPermissionsByRole(roleID uint) ([]*model.Permission, error)
}

type permissionService struct {
	permissionRepo repository.PermissionRepository
	roleRepo       repository.RoleRepository
}

// NewPermissionService 创建权限服务实例
func NewPermissionService(permissionRepo repository.PermissionRepository, roleRepo repository.RoleRepository) PermissionService {
	return &permissionService{
		permissionRepo: permissionRepo,
		roleRepo:       roleRepo,
	}
}

// Create 创建权限
func (s *permissionService) Create(permission *model.Permission) error {
	if permission.Name == "" {
		return errors.New("权限名称不能为空")
	}

	if permission.DisplayName == "" {
		return errors.New("权限显示名称不能为空")
	}

	if permission.Resource == "" {
		return errors.New("资源路径不能为空")
	}

	if permission.Action == "" {
		return errors.New("操作类型不能为空")
	}

	// 验证操作类型
	validActions := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"*":      true,
	}
	if !validActions[permission.Action] {
		return errors.New("无效的操作类型")
	}

	// 检查权限名称是否已存在
	existingPermission, err := s.permissionRepo.GetByName(permission.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("检查权限名称失败: %w", err)
	}
	if existingPermission != nil {
		return errors.New("权限名称已存在")
	}

	return s.permissionRepo.Create(permission)
}

// GetByID 根据ID获取权限
func (s *permissionService) GetByID(id uint) (*model.Permission, error) {
	if id == 0 {
		return nil, errors.New("权限ID不能为空")
	}
	return s.permissionRepo.GetByID(id)
}

// GetByName 根据名称获取权限
func (s *permissionService) GetByName(name string) (*model.Permission, error) {
	if name == "" {
		return nil, errors.New("权限名称不能为空")
	}
	return s.permissionRepo.GetByName(name)
}

// Update 更新权限
func (s *permissionService) Update(permission *model.Permission) error {
	if permission.ID == 0 {
		return errors.New("权限ID不能为空")
	}

	if permission.Name == "" {
		return errors.New("权限名称不能为空")
	}

	if permission.DisplayName == "" {
		return errors.New("权限显示名称不能为空")
	}

	if permission.Resource == "" {
		return errors.New("资源路径不能为空")
	}

	if permission.Action == "" {
		return errors.New("操作类型不能为空")
	}

	// 验证操作类型
	validActions := map[string]bool{
		"GET":    true,
		"POST":   true,
		"PUT":    true,
		"DELETE": true,
		"*":      true,
	}
	if !validActions[permission.Action] {
		return errors.New("无效的操作类型")
	}

	// 检查权限是否存在
	existingPermission, err := s.permissionRepo.GetByID(permission.ID)
	if err != nil {
		return fmt.Errorf("权限不存在: %w", err)
	}

	// 如果权限名称发生变化，检查新名称是否已存在
	if existingPermission.Name != permission.Name {
		conflictPermission, err := s.permissionRepo.GetByName(permission.Name)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("检查权限名称失败: %w", err)
		}
		if conflictPermission != nil {
			return errors.New("权限名称已存在")
		}
	}

	return s.permissionRepo.Update(permission)
}

// Delete 删除权限
func (s *permissionService) Delete(id uint) error {
	if id == 0 {
		return errors.New("权限ID不能为空")
	}

	// 检查权限是否存在
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("权限不存在: %w", err)
	}

	// 检查是否为系统内置权限，不允许删除
	if permission.Name == "system.all" {
		return errors.New("系统内置权限不允许删除")
	}

	return s.permissionRepo.Delete(id)
}

// List 获取权限列表
func (s *permissionService) List(page pagination.Pagination) ([]*model.Permission, int64, error) {
	return s.permissionRepo.List(page)
}

// UpdateStatus 更新权限状态
func (s *permissionService) UpdateStatus(id uint, status int) error {
	if id == 0 {
		return errors.New("权限ID不能为空")
	}

	if status != 0 && status != 1 {
		return errors.New("状态值无效")
	}

	// 检查权限是否存在
	_, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("权限不存在: %w", err)
	}

	return s.permissionRepo.UpdateStatus(id, status)
}

// GetPermissionsByRole 根据角色ID获取权限列表
func (s *permissionService) GetPermissionsByRole(roleID uint) ([]*model.Permission, error) {
	if roleID == 0 {
		return nil, errors.New("角色ID不能为空")
	}

	// 检查角色是否存在
	_, err := s.roleRepo.GetByID(roleID)
	if err != nil {
		return nil, fmt.Errorf("角色不存在: %w", err)
	}

	return s.permissionRepo.GetPermissionsByRole(roleID)
}