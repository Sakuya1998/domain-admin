package service

import (
	"context"

	pb "github.com/Sakuya1998/domain-admin/api/user/v1"
	commonv1 "github.com/Sakuya1998/domain-admin/api/common/v1"
	"github.com/Sakuya1998/domain-admin/app/user/internal/biz"
)

type UserService struct {
	pb.UnimplementedUserServer
	
	uc *biz.UserUsecase
}

func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{
		uc: uc,
	}
}

func (s *UserService) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterReply, error) {
	return s.uc.Register(ctx, req)
}

func (s *UserService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return s.uc.Login(ctx, req)
}

func (s *UserService) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutReply, error) {
	// TODO: Implement logout logic (e.g., invalidate token)
	return &pb.LogoutReply{}, nil
}

func (s *UserService) GetUserInfo(ctx context.Context, req *commonv1.IDRequest) (*pb.UserInfo, error) {
	return s.uc.GetUserInfo(ctx, req)
}

func (s *UserService) UpdateUserInfo(ctx context.Context, req *pb.UpdateUserInfoRequest) (*pb.UserInfo, error) {
	return s.uc.UpdateUserInfo(ctx, req)
}

func (s *UserService) UpdateUserPassword(ctx context.Context, req *pb.UpdateUserPasswordRequest) (*commonv1.EmptyReply, error) {
	err := s.uc.UpdateUserPassword(ctx, req)
	if err != nil {
		return nil, err
	}
	return &commonv1.EmptyReply{}, nil
}

func (s *UserService) UpdateUserStatus(ctx context.Context, req *pb.UpdateUserStatusRequest) (*commonv1.EmptyReply, error) {
	err := s.uc.UpdateUserStatus(ctx, req)
	if err != nil {
		return nil, err
	}
	return &commonv1.EmptyReply{}, nil
}

func (s *UserService) GetUserList(ctx context.Context, req *pb.GetUserListRequest) (*pb.GetUserListReply, error) {
	return s.uc.GetUserList(ctx, req)
}

func (s *UserService) AddUser(ctx context.Context, req *pb.AddUserRequest) (*pb.UserInfo, error) {
	return s.uc.AddUser(ctx, req)
}

func (s *UserService) DeleteUser(ctx context.Context, req *commonv1.IDRequest) (*commonv1.EmptyReply, error) {
	err := s.uc.DeleteUser(ctx, req)
	if err != nil {
		return nil, err
	}
	return &commonv1.EmptyReply{}, nil
}

func (s *UserService) ResetUserPassword(ctx context.Context, req *pb.ResetUserPasswordRequest) (*commonv1.EmptyReply, error) {
	err := s.uc.ResetUserPassword(ctx, req)
	if err != nil {
		return nil, err
	}
	return &commonv1.EmptyReply{}, nil
}
