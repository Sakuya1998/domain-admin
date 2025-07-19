package model

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID        uint           `json:"id" gorm:"primarykey"`
	Username  string         `json:"username" gorm:"uniqueIndex;size:50;not null" validate:"required,min=3,max=50"`
	Email     string         `json:"email" gorm:"uniqueIndex;size:100;not null" validate:"required,email"`
	Password  string         `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
	Nickname  string         `json:"nickname" gorm:"size:50"`
	Avatar    string         `json:"avatar" gorm:"size:255"`
	Phone     string         `json:"phone" gorm:"size:20"`
	Role      string         `json:"role" gorm:"size:20;default:user" validate:"required,oneof=admin user"`
	Status    int            `json:"status" gorm:"default:1;comment:1-正常 0-禁用"`
	LastLogin *time.Time     `json:"last_login"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}



// UserCreateRequest 创建用户请求
type UserCreateRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Nickname string `json:"nickname" validate:"max=50"`
	Phone    string `json:"phone" validate:"max=20"`
	Role     string `json:"role" validate:"omitempty,oneof=admin user"`
}

// UserUpdateRequest 更新用户请求
type UserUpdateRequest struct {
	Nickname string `json:"nickname" validate:"max=50"`
	Avatar   string `json:"avatar" validate:"max=255"`
	Phone    string `json:"phone" validate:"max=20"`
	Role     string `json:"role" validate:"omitempty,oneof=admin user"`
	Status   *int   `json:"status" validate:"omitempty,oneof=0 1"`
}

// UserLoginRequest 用户登录请求
type UserLoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        uint       `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	Phone     string     `json:"phone"`
	Role      string     `json:"role"`
	Status    int        `json:"status"`
	LastLogin *time.Time `json:"last_login"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// ToResponse 转换为响应格式
func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		Nickname:  u.Nickname,
		Avatar:    u.Avatar,
		Phone:     u.Phone,
		Role:      u.Role,
		Status:    u.Status,
		LastLogin: u.LastLogin,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}