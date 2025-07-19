package biz

import (
	"context"

	pb "github.com/Sakuya1998/domain-admin/api/user/v1"
	commonv1 "github.com/Sakuya1998/domain-admin/api/common/v1"

	"github.com/go-kratos/kratos/v2/log"
)

// User is a User model.
type User struct {
	ID       int64
	Username string
	Email    string
	Password string
	Role     commonv1.UserRole
	Status   commonv1.Status
}

// UserRepo is a Greater repo.
type UserRepo interface {
	Save(context.Context, *User) (*User, error)
	Update(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	FindByUsername(context.Context, string) (*User, error)
	FindByEmail(context.Context, string) (*User, error)
	ListByPage(context.Context, int32, int32) ([]*User, int32, error)
	Delete(context.Context, int64) error
	UpdatePassword(context.Context, int64, string) error
	UpdateStatus(context.Context, int64, commonv1.Status) error
}

// UserUsecase is a User usecase.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// NewUserUsecase new a User usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{repo: repo, log: log.NewHelper(logger)}
}

// Register creates a User, and returns the new User.
func (uc *UserUsecase) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	uc.log.WithContext(ctx).Infof("Register: %v", req.Username)
	
	// Check if user already exists
	existingUser, err := uc.repo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, pb.ErrorUserAlreadyExistsWithMsg("user already exists")
	}
	
	// Check if email already exists
	existingUser, err = uc.repo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, pb.ErrorEmailAlreadyExistsWithMsg("email already exists")
	}
	
	// Create new user
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // In real implementation, hash the password
		Role:     commonv1.UserRole_USER_ROLE_USER,
		Status:   commonv1.Status_STATUS_ACTIVE,
	}
	
	savedUser, err := uc.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return &pb.RegisterReply{
		UserInfo: &pb.UserInfo{
			Id:       savedUser.ID,
			Username: savedUser.Username,
			Email:    savedUser.Email,
			Role:     savedUser.Role,
			Status:   savedUser.Status,
		},
	}, nil
}

// Login authenticates a user.
func (uc *UserUsecase) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	var user *User
	var err error
	
	// Handle oneof field for account
	if req.GetUsername() != "" {
		uc.log.WithContext(ctx).Infof("Login by username: %v", req.GetUsername())
		user, err = uc.repo.FindByUsername(ctx, req.GetUsername())
	} else if req.GetEmail() != "" {
		uc.log.WithContext(ctx).Infof("Login by email: %v", req.GetEmail())
		user, err = uc.repo.FindByEmail(ctx, req.GetEmail())
	} else {
		return nil, pb.ErrorUserNotFoundWithMsg("username or email is required")
	}
	
	if err != nil {
		return nil, pb.ErrorUserNotFoundWithMsg("user not found")
	}
	
	// In real implementation, verify hashed password
	if user.Password != req.Password {
		return nil, pb.ErrorInvalidCredentialsWithMsg("invalid credentials")
	}
	
	if user.Status != commonv1.Status_STATUS_ACTIVE {
		return nil, pb.ErrorUserInactiveWithMsg("user is inactive")
	}
	
	return &pb.LoginReply{
		Token: "mock-jwt-token", // In real implementation, generate JWT token
		UserInfo: &pb.UserInfo{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		},
	}, nil
}

// GetUserInfo gets user information by ID.
func (uc *UserUsecase) GetUserInfo(ctx context.Context, req *commonv1.IDRequest) (*pb.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("GetUserInfo: %v", req.Id)
	
	user, err := uc.repo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, pb.ErrorUserNotFoundWithMsg("user not found")
	}
	
	return &pb.UserInfo{
		Id:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Status:   user.Status,
	}, nil
}

// UpdateUserInfo updates user information.
func (uc *UserUsecase) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("UpdateUserInfo: %v", req.Id)
	
	user, err := uc.repo.FindByID(ctx, req.Id)
	if err != nil {
		return nil, pb.ErrorUserNotFoundWithMsg("user not found")
	}
	
	if req.Email != "" {
		user.Email = req.Email
	}
	
	updatedUser, err := uc.repo.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return &pb.UserInfo{
		Id:       updatedUser.ID,
		Username: updatedUser.Username,
		Email:    updatedUser.Email,
		Role:     updatedUser.Role,
		Status:   updatedUser.Status,
	}, nil
}

// UpdateUserPassword updates user password.
func (uc *UserUsecase) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) error {
	uc.log.WithContext(ctx).Infof("UpdateUserPassword: %v", req.Id)
	
	user, err := uc.repo.FindByID(ctx, req.Id)
	if err != nil {
		return pb.ErrorUserNotFoundWithMsg("user not found")
	}
	
	// In real implementation, verify old password and hash new password
	if user.Password != req.OldPassword {
		return pb.ErrorInvalidCredentialsWithMsg("invalid old password")
	}
	
	return uc.repo.UpdatePassword(ctx, req.Id, req.NewPassword)
}

// UpdateUserStatus updates user status.
func (uc *UserUsecase) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) error {
	uc.log.WithContext(ctx).Infof("UpdateUserStatus: %v", req.Id)
	
	return uc.repo.UpdateStatus(ctx, req.Id, req.Status)
}

// GetUserList gets user list with pagination.
func (uc *UserUsecase) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListReply, error) {
	uc.log.WithContext(ctx).Infof("GetUserList: page=%v, size=%v", req.Page.Page, req.Page.PageSize)
	
	users, _, err := uc.repo.ListByPage(ctx, req.Page.Page, req.Page.PageSize)
	if err != nil {
		return nil, err
	}
	
	userInfos := make([]*pb.UserInfo, 0, len(users))
	for _, user := range users {
		userInfos = append(userInfos, &pb.UserInfo{
			Id:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     user.Role,
			Status:   user.Status,
		})
	}
	
	return &pb.GetUserListReply{
		Users: userInfos,
	}, nil
}

// AddUser adds a new user.
func (uc *UserUsecase) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.UserInfo, error) {
	uc.log.WithContext(ctx).Infof("AddUser: %v", req.Username)
	
	// Check if user already exists
	existingUser, err := uc.repo.FindByUsername(ctx, req.Username)
	if err == nil && existingUser != nil {
		return nil, pb.ErrorUserAlreadyExistsWithMsg("user already exists")
	}
	
	// Check if email already exists
	existingUser, err = uc.repo.FindByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, pb.ErrorEmailAlreadyExistsWithMsg("email already exists")
	}
	
	// Create new user
	user := &User{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password, // In real implementation, hash the password
		Role:     req.Role,
		Status:   commonv1.Status_STATUS_ACTIVE,
	}
	
	savedUser, err := uc.repo.Save(ctx, user)
	if err != nil {
		return nil, err
	}
	
	return &pb.UserInfo{
		Id:       savedUser.ID,
		Username: savedUser.Username,
		Email:    savedUser.Email,
		Role:     savedUser.Role,
		Status:   savedUser.Status,
	}, nil
}

// DeleteUser deletes a user.
func (uc *UserUsecase) DeleteUser(ctx context.Context, req *commonv1.IDRequest) error {
	uc.log.WithContext(ctx).Infof("DeleteUser: %v", req.Id)
	
	return uc.repo.Delete(ctx, req.Id)
}

// ResetUserPassword resets user password.
func (uc *UserUsecase) ResetUserPassword(ctx context.Context, req *pb.ResetUserPasswordRequest) error {
	uc.log.WithContext(ctx).Infof("ResetUserPassword: %v", req.Id)
	
	// In real implementation, generate a random password or send reset email
	newPassword := "reset123456" // This should be generated randomly
	return uc.repo.UpdatePassword(ctx, req.Id, newPassword)
}