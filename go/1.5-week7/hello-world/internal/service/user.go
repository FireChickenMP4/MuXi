package service

import (
	"context"
	"fmt"
	"strconv"

	v1 "hello-world/api/user/v1"
	"hello-world/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/metadata"
)

// UserService 用户相关 service。
type UserService struct {
	v1.UnimplementedUserServer
	uc *biz.UserUsecase
}

// NewUserService new a user service.
func NewUserService(uc *biz.UserUsecase) *UserService {
	return &UserService{uc: uc}
}

// Register 用户注册。
func (s *UserService) Register(ctx context.Context, in *v1.RegisterRequest) (*v1.RegisterReply, error) {
	user, err := s.uc.Register(ctx, &biz.User{
		Username: in.GetUsername(),
		Password: in.GetPassword(),
		Email:    in.GetEmail(),
		Mobile:   in.GetMobile(),
	})
	if err != nil {
		return nil, err
	}

	return &v1.RegisterReply{
		User:  toUserInfo(user),
		Token: fmt.Sprintf("register-token-%d", user.ID),
	}, nil
}

// Login 用户登录。
func (s *UserService) Login(ctx context.Context, in *v1.LoginRequest) (*v1.LoginReply, error) {
	user, err := s.uc.Login(ctx, in.GetAccount(), in.GetPassword())
	if err != nil {
		return nil, err
	}

	return &v1.LoginReply{
		User:  toUserInfo(user),
		Token: fmt.Sprintf("login-token-%d", user.ID),
	}, nil
}

// GetProfile 获取当前用户信息。
func (s *UserService) GetProfile(ctx context.Context, in *v1.GetProfileRequest) (*v1.GetProfileReply, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.uc.GetProfile(ctx, userID)
	if err != nil {
		return nil, err
	}
	return &v1.GetProfileReply{
		User: toUserInfo(user),
	}, nil
}

// UpdateProfile 更新当前用户信息。
func (s *UserService) UpdateProfile(ctx context.Context, in *v1.UpdateProfileRequest) (*v1.UpdateProfileReply, error) {
	userID, err := currentUserID(ctx)
	if err != nil {
		return nil, err
	}
	user, err := s.uc.UpdateProfile(ctx, &biz.User{
		ID:       userID,
		Username: in.GetUsername(),
		Email:    in.GetEmail(),
		Mobile:   in.GetMobile(),
		Avatar:   in.GetAvatar(),
	})
	if err != nil {
		return nil, err
	}
	return &v1.UpdateProfileReply{
		User: toUserInfo(user),
	}, nil
}

func toUserInfo(user *biz.User) *v1.UserInfo {
	if user == nil {
		return nil
	}
	return &v1.UserInfo{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Avatar:    user.Avatar,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func currentUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return 0, errors.Unauthorized("UNAUTHORIZED", "missing user identity")
	}

	for _, key := range []string{"user_id", "x-user-id"} {
		if values := md[key]; len(values) > 0 && values[0] != "" {
			if id, err := strconv.ParseInt(values[0], 10, 64); err == nil && id > 0 {
				return id, nil
			}
		}
	}

	return 0, errors.Unauthorized("UNAUTHORIZED", "missing user identity")
}
