package data

import (
	"context"
	"strings"

	"hello-world/internal/biz"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/log"
)

type userRepo struct {
	data *Data
	log  *log.Helper
}

// NewUserRepo creates a user repository implementation.
func NewUserRepo(data *Data, logger log.Logger) biz.UserRepo {
	return &userRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *userRepo) Create(ctx context.Context, user *biz.User) (*biz.User, error) {
	if user == nil {
		return nil, errors.BadRequest("USER_REQUIRED", "user is required")
	}
	if user.Username == "" {
		return nil, errors.BadRequest("USERNAME_REQUIRED", "username is required")
	}
	if user.Password == "" {
		return nil, errors.BadRequest("PASSWORD_REQUIRED", "password is required")
	}

	r.data.mu.Lock()
	defer r.data.mu.Unlock()

	if _, ok := r.data.accountIndex[user.Username]; ok {
		return nil, errors.Conflict("USERNAME_EXISTS", "username already exists")
	}
	if user.Email != "" {
		if _, ok := r.data.accountIndex[user.Email]; ok {
			return nil, errors.Conflict("EMAIL_EXISTS", "email already exists")
		}
	}
	if user.Mobile != "" {
		if _, ok := r.data.accountIndex[user.Mobile]; ok {
			return nil, errors.Conflict("MOBILE_EXISTS", "mobile already exists")
		}
	}

	if user.ID <= 0 {
		user.ID = r.data.nextUserID.Add(1)
	}
	record := toRecord(user)
	r.data.users[user.ID] = record
	r.data.accountIndex[user.Username] = user.ID
	if user.Email != "" {
		r.data.accountIndex[user.Email] = user.ID
	}
	if user.Mobile != "" {
		r.data.accountIndex[user.Mobile] = user.ID
	}

	r.log.Infof("create user: id=%d username=%s", user.ID, user.Username)
	return toDomain(record), nil
}

func (r *userRepo) FindByID(ctx context.Context, id int64) (*biz.User, error) {
	r.data.mu.RLock()
	defer r.data.mu.RUnlock()

	record, ok := r.data.users[id]
	if !ok {
		return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
	}
	return toDomain(record), nil
}

func (r *userRepo) FindByAccount(ctx context.Context, account string) (*biz.User, error) {
	account = strings.TrimSpace(account)
	if account == "" {
		return nil, errors.BadRequest("ACCOUNT_REQUIRED", "account is required")
	}

	r.data.mu.RLock()
	defer r.data.mu.RUnlock()

	if id, ok := r.data.accountIndex[account]; ok {
		record, ok := r.data.users[id]
		if ok {
			return toDomain(record), nil
		}
	}

	for _, record := range r.data.users {
		if record.Username == account || record.Email == account || record.Mobile == account {
			return toDomain(record), nil
		}
	}

	return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
}

func (r *userRepo) Update(ctx context.Context, user *biz.User) (*biz.User, error) {
	if user == nil {
		return nil, errors.BadRequest("USER_REQUIRED", "user is required")
	}
	if user.ID <= 0 {
		return nil, errors.BadRequest("USER_ID_REQUIRED", "user id is required")
	}

	r.data.mu.Lock()
	defer r.data.mu.Unlock()

	existing, ok := r.data.users[user.ID]
	if !ok {
		return nil, errors.NotFound("USER_NOT_FOUND", "user not found")
	}

	delete(r.data.accountIndex, existing.Username)
	if existing.Email != "" {
		delete(r.data.accountIndex, existing.Email)
	}
	if existing.Mobile != "" {
		delete(r.data.accountIndex, existing.Mobile)
	}

	updated := existing
	if user.Username != "" {
		updated.Username = user.Username
	}
	if user.Email != "" {
		updated.Email = user.Email
	}
	if user.Mobile != "" {
		updated.Mobile = user.Mobile
	}
	if user.Avatar != "" {
		updated.Avatar = user.Avatar
	}
	if user.Password != "" {
		updated.Password = user.Password
	}
	updated.UpdatedAt = user.UpdatedAt

	r.data.users[user.ID] = updated
	r.data.accountIndex[updated.Username] = updated.ID
	if updated.Email != "" {
		r.data.accountIndex[updated.Email] = updated.ID
	}
	if updated.Mobile != "" {
		r.data.accountIndex[updated.Mobile] = updated.ID
	}

	r.log.Infof("update user: id=%d username=%s", updated.ID, updated.Username)
	return toDomain(updated), nil
}

func toDomain(record UserRecord) *biz.User {
	return &biz.User{
		ID:        record.ID,
		Username:  record.Username,
		Email:     record.Email,
		Mobile:    record.Mobile,
		Avatar:    record.Avatar,
		Password:  record.Password,
		CreatedAt: record.CreatedAt,
		UpdatedAt: record.UpdatedAt,
	}
}

func toRecord(user *biz.User) UserRecord {
	return UserRecord{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Mobile:    user.Mobile,
		Avatar:    user.Avatar,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
