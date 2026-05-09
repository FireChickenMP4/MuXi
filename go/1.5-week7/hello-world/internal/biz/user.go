package biz

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/wire"
)

var (
	// ErrUserUsernameRequired 用户名不能为空。
	ErrUserUsernameRequired = errors.BadRequest("USERNAME_REQUIRED", "username is required")
	// ErrUserPasswordRequired 密码不能为空。
	ErrUserPasswordRequired = errors.BadRequest("PASSWORD_REQUIRED", "password is required")
	// ErrUserAccountRequired 登录账号不能为空。
	ErrUserAccountRequired = errors.BadRequest("ACCOUNT_REQUIRED", "account is required")
)

// User is the user domain model.
type User struct {
	ID        int64
	Username  string
	Email     string
	Mobile    string
	Avatar    string
	Password  string
	CreatedAt int64
	UpdatedAt int64
}

// UserRepo is the user repository interface.
type UserRepo interface {
	Create(context.Context, *User) (*User, error)
	FindByID(context.Context, int64) (*User, error)
	FindByAccount(context.Context, string) (*User, error)
	Update(context.Context, *User) (*User, error)
}

// UserUsecase handles user business rules.
type UserUsecase struct {
	repo UserRepo
	log  *log.Helper
}

// ProviderSet is biz providers.
var UserProviderSet = wire.NewSet(NewUserUsecase)

// NewUserUsecase new a user usecase.
func NewUserUsecase(repo UserRepo, logger log.Logger) *UserUsecase {
	return &UserUsecase{
		repo: repo,
		log:  log.NewHelper(logger),
	}
}

// Register creates a new user.
func (uc *UserUsecase) Register(ctx context.Context, user *User) (*User, error) {
	if user == nil {
		return nil, ErrUserUsernameRequired
	}
	if user.Username == "" {
		return nil, ErrUserUsernameRequired
	}
	if user.Password == "" {
		return nil, ErrUserPasswordRequired
	}

	now := time.Now().Unix()
	user.CreatedAt = now
	user.UpdatedAt = now

	uc.log.Infof("register user: username=%s email=%s mobile=%s", user.Username, user.Email, user.Mobile)
	return uc.repo.Create(ctx, user)
}

// Login verifies user credentials.
func (uc *UserUsecase) Login(ctx context.Context, account, password string) (*User, error) {
	if account == "" {
		return nil, ErrUserAccountRequired
	}
	if password == "" {
		return nil, ErrUserPasswordRequired
	}

	uc.log.Infof("login user: account=%s", account)
	user, err := uc.repo.FindByAccount(ctx, account)
	if err != nil {
		return nil, err
	}
	if user.Password != password {
		return nil, errors.Unauthorized("PASSWORD_INVALID", "password is invalid")
	}
	return user, nil
}

// GetProfile returns user profile.
func (uc *UserUsecase) GetProfile(ctx context.Context, id int64) (*User, error) {
	if id <= 0 {
		return nil, errors.BadRequest("USER_ID_REQUIRED", "user id is required")
	}
	return uc.repo.FindByID(ctx, id)
}

// UpdateProfile updates user profile.
func (uc *UserUsecase) UpdateProfile(ctx context.Context, user *User) (*User, error) {
	if user == nil {
		return nil, errors.BadRequest("USER_REQUIRED", "user is required")
	}
	if user.ID <= 0 {
		return nil, errors.BadRequest("USER_ID_REQUIRED", "user id is required")
	}

	user.UpdatedAt = time.Now().Unix()
	return uc.repo.Update(ctx, user)
}
