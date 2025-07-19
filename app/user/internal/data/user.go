package data

import (
	"context"
	"time"

	commonv1 "github.com/Sakuya1998/domain-admin/api/common/v1"
	"github.com/Sakuya1998/domain-admin/app/user/internal/biz"

	"github.com/go-kratos/kratos/v2/log"
	"gorm.io/gorm"
)

// User GORM模型
type User struct {
	ID          int64             `gorm:"primaryKey;autoIncrement" json:"id"`
	Username    string            `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email       string            `gorm:"uniqueIndex;size:100;not null" json:"email"`
	Password    string            `gorm:"size:255;not null" json:"-"`
	Role        commonv1.UserRole `gorm:"type:int;default:1" json:"role"`
	Status      commonv1.Status   `gorm:"type:int;default:1" json:"status"`
	CreatedAt   time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	LastLoginAt time.Time         `gorm:"autoUpdateTime" json:"last_login_at"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// ToBizUser 转换为业务层用户对象
func (u *User) ToBizUser() *biz.User {
	return &biz.User{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Password: u.Password,
		Role:     u.Role,
		Status:   u.Status,
	}
}

// FromBizUser 从业务层用户对象转换
func (u *User) FromBizUser(bizUser *biz.User) {
	u.ID = bizUser.ID
	u.Username = bizUser.Username
	u.Email = bizUser.Email
	u.Password = bizUser.Password
	u.Role = bizUser.Role
	u.Status = bizUser.Status
}

// userRepo 用户仓储实现
type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo creates a new user repository.
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

// Save saves a user to the repository.
func (r *userRepo) Save(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("Save user: %v", user.Username)

	// 转换为数据模型
	dbUser := &User{}
	dbUser.FromBizUser(user)

	// 保存到数据库
	if err := r.data.db.WithContext(ctx).Create(dbUser).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to save user: %v", err)
		return nil, err
	}

	// 返回保存后的用户（包含生成的ID）
	return dbUser.ToBizUser(), nil
}

// Update updates a user in the repository.
func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("Update user: %v", user.ID)

	// 转换为数据模型
	dbUser := &User{}
	dbUser.FromBizUser(user)

	// 更新数据库记录
	if err := r.data.db.WithContext(ctx).Save(dbUser).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update user: %v", err)
		return nil, err
	}

	return dbUser.ToBizUser(), nil
}

// FindByID finds a user by ID.
func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("FindByID: %v", id)

	var dbUser User
	if err := r.data.db.WithContext(ctx).First(&dbUser, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find user by ID: %v", err)
		return nil, err
	}

	return dbUser.ToBizUser(), nil
}

// FindByUsername finds a user by username.
func (r *userRepo) FindByUsername(ctx context.Context, username string) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("FindByUsername: %v", username)

	var dbUser User
	if err := r.data.db.WithContext(ctx).Where("username = ?", username).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find user by username: %v", err)
		return nil, err
	}

	return dbUser.ToBizUser(), nil
}

// FindByEmail finds a user by email.
func (r *userRepo) FindByEmail(ctx context.Context, email string) (*biz.User, error) {
	r.log.WithContext(ctx).Infof("FindByEmail: %v", email)

	var dbUser User
	if err := r.data.db.WithContext(ctx).Where("email = ?", email).First(&dbUser).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		r.log.WithContext(ctx).Errorf("Failed to find user by email: %v", err)
		return nil, err
	}

	return dbUser.ToBizUser(), nil
}

// ListByPage returns a paginated list of users.
func (r *userRepo) ListByPage(ctx context.Context, page, size int32) ([]*biz.User, int32, error) {
	r.log.WithContext(ctx).Infof("ListByPage: page=%v, size=%v", page, size)

	var dbUsers []User
	var total int64

	// 获取总数
	if err := r.data.db.WithContext(ctx).Model(&User{}).Count(&total).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to count users: %v", err)
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * size
	if err := r.data.db.WithContext(ctx).Offset(int(offset)).Limit(int(size)).Find(&dbUsers).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to list users: %v", err)
		return nil, 0, err
	}

	// 转换为业务模型
	var users []*biz.User
	for _, dbUser := range dbUsers {
		users = append(users, dbUser.ToBizUser())
	}

	return users, int32(total), nil
}

// Delete deletes a user by ID.
func (r *userRepo) Delete(ctx context.Context, id int64) error {
	r.log.WithContext(ctx).Infof("Delete user: %v", id)

	if err := r.data.db.WithContext(ctx).Delete(&User{}, id).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to delete user: %v", err)
		return err
	}

	return nil
}

// UpdatePassword updates a user's password.
func (r *userRepo) UpdatePassword(ctx context.Context, id int64, hashedPassword string) error {
	r.log.WithContext(ctx).Infof("UpdatePassword for user: %v", id)

	if err := r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("password", hashedPassword).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update password: %v", err)
		return err
	}

	return nil
}

// UpdateStatus updates a user's status.
func (r *userRepo) UpdateStatus(ctx context.Context, id int64, status commonv1.Status) error {
	r.log.WithContext(ctx).Infof("UpdateStatus: %v, status=%v", id, status)

	if err := r.data.db.WithContext(ctx).Model(&User{}).Where("id = ?", id).Update("status", status).Error; err != nil {
		r.log.WithContext(ctx).Errorf("Failed to update status: %v", err)
		return err
	}

	return nil
}
