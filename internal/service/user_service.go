package service

import (
	"context"
	"domain-admin/internal/repository"
	"domain-admin/model"
	"domain-admin/pkg/cache"
	"domain-admin/pkg/jwt"
	"domain-admin/pkg/logger"
	"domain-admin/pkg/pagination"
	"errors"
	"fmt"
	"strconv"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserService 用户服务接口
type UserService interface {
	Register(req *model.UserCreateRequest) (*model.UserResponse, error)
	Login(req *model.UserLoginRequest) (string, *model.UserResponse, error)
	Logout(userID uint) error
	GetProfile(userID uint) (*model.UserResponse, error)
	UpdateProfile(userID uint, req *model.UserUpdateRequest) (*model.UserResponse, error)
	GetUserList(page pagination.Pagination) (*pagination.PageResult, error)
	GetUserByID(id uint) (*model.UserResponse, error)
	CreateUser(req *model.UserCreateRequest) (*model.UserResponse, error)
	UpdateUser(id uint, req *model.UserUpdateRequest) (*model.UserResponse, error)
	DeleteUser(id uint) error
	UpdateUserStatus(id uint, status int) error
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) UserService {
	return &userService{
		userRepo: repository.NewUserRepository(db),
	}
}

// Register 用户注册
func (s *userService) Register(req *model.UserCreateRequest) (*model.UserResponse, error) {
	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("密码加密失败: %v", err)
		return nil, errors.New("密码加密失败")
	}

	// 设置默认角色
	role := req.Role
	if role == "" {
		role = "user" // 注册用户默认为普通用户
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Role:     role,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	logger.Infof("用户注册成功: %s", user.Username)
	return user.ToResponse(), nil
}

// Login 用户登录
func (s *userService) Login(req *model.UserLoginRequest) (string, *model.UserResponse, error) {
	ctx := context.Background()

	// 获取用户信息
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status == 0 {
		return "", nil, errors.New("用户已被禁用")
	}

	// 验证密码
	if PasswordErr := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); PasswordErr != nil {
		return "", nil, errors.New("用户名或密码错误")
	}

	// 生成JWT token
	token, tokenErr := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if tokenErr != nil {
		logger.Errorf("生成token失败: %v", tokenErr)
		return "", nil, errors.New("登录失败")
	}

	// 缓存用户会话信息
	sessionInfo := map[string]interface{}{
		"user_id":  user.ID,
		"username": user.Username,
		"role":     user.Role,
		"status":   user.Status,
	}
	sessionKey := strconv.FormatUint(uint64(user.ID), 10)
	if err := cache.SetSessionCache(ctx, sessionKey, sessionInfo); err != nil {
		logger.Warnf("缓存会话信息失败: %v", err)
	}

	// 缓存用户信息
	userResponse := user.ToResponse()
	if err := cache.SetUserCache(ctx, user.ID, userResponse); err != nil {
		logger.Warnf("缓存用户信息失败: %v", err)
	}

	// 更新最后登录时间
	if err := s.userRepo.UpdateLastLogin(user.ID); err != nil {
		logger.Warnf("更新最后登录时间失败: %v", err)
	}

	logger.Infof("用户登录成功: %s", user.Username)
	return token, userResponse, nil
}

// Logout 用户登出
func (s *userService) Logout(userID uint) error {
	ctx := context.Background()

	// 删除会话缓存
	sessionKey := strconv.FormatUint(uint64(userID), 10)
	if err := cache.DelSessionCache(ctx, sessionKey); err != nil {
		logger.Warnf("删除会话缓存失败: %v", err)
		return errors.New("登出失败")
	}

	logger.Infof("用户登出成功，用户ID: %d", userID)
	return nil
}

// GetProfile 获取用户资料
func (s *userService) GetProfile(userID uint) (*model.UserResponse, error) {
	ctx := context.Background()

	// 先尝试从缓存获取
	var cachedUser model.UserResponse
	if err := cache.GetUserCache(ctx, userID, &cachedUser); err == nil {
		logger.Debugf("从缓存获取用户资料: %d", userID)
		return &cachedUser, nil
	} else if !errors.Is(err, redis.Nil) {
		logger.Warnf("获取用户缓存失败: %v", err)
	}

	// 缓存未命中，从数据库获取
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	userResponse := user.ToResponse()

	// 更新缓存
	if err := cache.SetUserCache(ctx, userID, userResponse); err != nil {
		logger.Warnf("缓存用户信息失败: %v", err)
	}

	return userResponse, nil
}

// UpdateProfile 更新用户资料
func (s *userService) UpdateProfile(userID uint, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	ctx := context.Background()

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if err := s.userRepo.Update(user); err != nil {
		logger.Errorf("更新用户资料失败: %v", err)
		return nil, errors.New("更新用户资料失败")
	}

	userResponse := user.ToResponse()

	// 更新缓存
	if err := cache.SetUserCache(ctx, userID, userResponse); err != nil {
		logger.Warnf("更新用户缓存失败: %v", err)
	}

	// 清除用户列表缓存
	if err := cache.DelUserListCache(ctx, "*"); err != nil {
		logger.Warnf("清除用户列表缓存失败: %v", err)
	}

	logger.Infof("用户资料更新成功: %s", user.Username)
	return userResponse, nil
}

// GetUserList 获取用户列表（管理员功能）
func (s *userService) GetUserList(page pagination.Pagination) (*pagination.PageResult, error) {
	ctx := context.Background()

	// 生成缓存键
	cacheKey := fmt.Sprintf("offset_%d_limit_%d_order_%s", page.Offset, page.Limit, page.GetOrderClause())

	// 先尝试从缓存获取
	var cachedResult pagination.PageResult
	if err := cache.GetUserListCache(ctx, cacheKey, &cachedResult); err == nil {
		logger.Debugf("从缓存获取用户列表: %s", cacheKey)
		return &cachedResult, nil
	} else if !errors.Is(err, redis.Nil) {
		logger.Warnf("获取用户列表缓存失败: %v", err)
	}

	// 缓存未命中，从数据库获取
	users, total, err := s.userRepo.List(page)
	if err != nil {
		return nil, err
	}

	// 转换为响应格式
	userResponses := make([]*model.UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = user.ToResponse()
	}

	result := &pagination.PageResult{
		Total: total,
		Items: userResponses,
	}

	// 缓存结果
	if err := cache.SetUserListCache(ctx, cacheKey, result); err != nil {
		logger.Warnf("缓存用户列表失败: %v", err)
	}

	return result, nil
}

// GetUserByID 根据ID获取用户（管理员功能）
func (s *userService) GetUserByID(id uint) (*model.UserResponse, error) {
	ctx := context.Background()

	// 先尝试从缓存获取
	var cachedUser model.UserResponse
	if err := cache.GetUserCache(ctx, id, &cachedUser); err == nil {
		logger.Debugf("从缓存获取用户信息: %d", id)
		return &cachedUser, nil
	} else if !errors.Is(err, redis.Nil) {
		logger.Warnf("获取用户缓存失败: %v", err)
	}

	// 缓存未命中，从数据库获取
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	userResponse := user.ToResponse()

	// 更新缓存
	if err := cache.SetUserCache(ctx, id, userResponse); err != nil {
		logger.Warnf("缓存用户信息失败: %v", err)
	}

	return userResponse, nil
}

// CreateUser 创建用户（管理员功能）
func (s *userService) CreateUser(req *model.UserCreateRequest) (*model.UserResponse, error) {
	ctx := context.Background()

	// 检查用户名是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Errorf("密码加密失败: %v", err)
		return nil, errors.New("密码加密失败")
	}

	// 设置默认角色
	role := req.Role
	if role == "" {
		role = "user" // 默认为普通用户
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Phone:    req.Phone,
		Role:     role,
		Status:   1,
	}

	if err := s.userRepo.Create(user); err != nil {
		logger.Errorf("创建用户失败: %v", err)
		return nil, errors.New("创建用户失败")
	}

	userResponse := user.ToResponse()

	// 缓存新用户信息
	if err := cache.SetUserCache(ctx, user.ID, userResponse); err != nil {
		logger.Warnf("缓存新用户信息失败: %v", err)
	}

	// 清除用户列表缓存
	if err := cache.DelUserListCache(ctx, "*"); err != nil {
		logger.Warnf("清除用户列表缓存失败: %v", err)
	}

	logger.Infof("管理员创建用户成功: %s", user.Username)
	return userResponse, nil
}

// UpdateUser 更新用户（管理员功能）
func (s *userService) UpdateUser(id uint, req *model.UserUpdateRequest) (*model.UserResponse, error) {
	ctx := context.Background()

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// 更新字段
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.userRepo.Update(user); err != nil {
		logger.Errorf("更新用户失败: %v", err)
		return nil, errors.New("更新用户失败")
	}

	userResponse := user.ToResponse()

	// 更新用户缓存
	if err := cache.SetUserCache(ctx, id, userResponse); err != nil {
		logger.Warnf("更新用户缓存失败: %v", err)
	}

	// 清除用户列表缓存
	if err := cache.DelUserListCache(ctx, "*"); err != nil {
		logger.Warnf("清除用户列表缓存失败: %v", err)
	}

	logger.Infof("管理员更新用户成功: %s", user.Username)
	return userResponse, nil
}

// DeleteUser 删除用户（管理员功能）
func (s *userService) DeleteUser(id uint) error {
	ctx := context.Background()

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if err := s.userRepo.Delete(id); err != nil {
		logger.Errorf("删除用户失败: %v", err)
		return errors.New("删除用户失败")
	}

	// 删除用户缓存
	if err := cache.DelUserCache(ctx, id); err != nil {
		logger.Warnf("删除用户缓存失败: %v", err)
	}

	// 清除用户列表缓存
	if err := cache.DelUserListCache(ctx, "*"); err != nil {
		logger.Warnf("清除用户列表缓存失败: %v", err)
	}

	logger.Infof("管理员删除用户成功: %s", user.Username)
	return nil
}

// UpdateUserStatus 更新用户状态（管理员功能）
func (s *userService) UpdateUserStatus(id uint, status int) error {
	ctx := context.Background()

	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if UpdateStatusErr := s.userRepo.UpdateStatus(id, status); UpdateStatusErr != nil {
		logger.Errorf("更新用户状态失败: %v", UpdateStatusErr)
		return errors.New("更新用户状态失败")
	}

	// 更新用户状态后，需要重新获取用户信息并更新缓存
	updatedUser, err := s.userRepo.GetByID(id)
	if err == nil {
		userResponse := updatedUser.ToResponse()
		if err := cache.SetUserCache(ctx, id, userResponse); err != nil {
			logger.Warnf("更新用户状态缓存失败: %v", err)
		}
	}

	// 清除用户列表缓存
	if err := cache.DelUserListCache(ctx, "*"); err != nil {
		logger.Warnf("清除用户列表缓存失败: %v", err)
	}

	statusText := "启用"
	if status == 0 {
		statusText = "禁用"
	}
	logger.Infof("管理员%s用户成功: %s", statusText, user.Username)
	return nil
}
