package biz

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yz626/edu-chain/internal/data/repository/model"
	"github.com/yz626/edu-chain/internal/utils/crypto"
	"github.com/yz626/edu-chain/internal/utils/jwts"
	cerr "github.com/yz626/edu-chain/pkg/errors"
	"gorm.io/gorm"
)

// UserBiz 用户业务
type UserBiz struct {
	*Base
}

func NewUserBiz(base *Base) *UserBiz {
	return &UserBiz{Base: base}
}

// RegisterInput 注册参数
type RegisterInput struct {
	Username       string
	Email          string
	Phone          string
	Password       string
	RealName       string
	OrganizationID string
	RoleID         string
}

// LoginInput 登录参数
type LoginInput struct {
	Identifier string // 用户名/邮箱/手机号
	Password   string
}

// LoginOutput 登录结果
type LoginOutput struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
	User         *model.User
}

// UserListFilter 用户列表过滤条件
type UserListFilter struct {
	Page           int
	PageSize       int
	Keyword        string
	OrganizationID string
	Status         int32
}

func (b *UserBiz) Register(ctx context.Context, in RegisterInput) (*model.User, error) {
	if in.Username == "" || in.Email == "" || in.Password == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("用户名、邮箱、密码不能为空")
	}
	if err := crypto.ValidatePassword(in.Password); err != nil {
		return nil, err
	}

	db := withCtx(ctx, b.db)

	var exists int64
	if err := db.Model(&model.User{}).
		Where("username = ? OR email = ?", in.Username, in.Email).
		Count(&exists).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	if exists > 0 {
		return nil, cerr.New(cerr.ErrCodeUserAlreadyExists, "用户名或邮箱已存在")
	}

	hash, err := crypto.HashPassword(in.Password)
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	now := time.Now()
	user := &model.User{
		ID:           uuid.NewString(),
		Username:     in.Username,
		Email:        in.Email,
		Phone:        in.Phone,
		PasswordHash: hash,
		Status:       1,
		UserType:     1,
		Source:       1,
		RealName:     in.RealName,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		if in.OrganizationID != "" {
			orgUser := &model.OrganizationUser{
				ID:             uuid.NewString(),
				OrganizationID: in.OrganizationID,
				UserID:         user.ID,
				IsPrimary:      true,
				JoinedAt:       now,
				CreatedAt:      now,
				UpdatedAt:      now,
			}
			if err := tx.Create(orgUser).Error; err != nil {
				return err
			}
		}

		if in.RoleID != "" {
			userRole := &model.UserRole{
				ID:        uuid.NewString(),
				UserID:    user.ID,
				RoleID:    in.RoleID,
				IsActive:  true,
				GrantedAt: now,
				CreatedAt: now,
				UpdatedAt: now,
			}
			if err := tx.Create(userRole).Error; err != nil {
				return err
			}
		}

		return nil
	})
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return user, nil
}

func (b *UserBiz) Login(ctx context.Context, in LoginInput) (*LoginOutput, error) {
	if in.Identifier == "" || in.Password == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("登录参数不能为空")
	}

	db := withCtx(ctx, b.db)
	var user model.User
	err := db.Where("username = ? OR email = ? OR phone = ?", in.Identifier, in.Identifier, in.Identifier).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.New(cerr.ErrCodeUserNotFound, "用户不存在")
		}
		return nil, cerr.ErrInternal.WithError(err)
	}

	if user.Status != 1 {
		return nil, cerr.New(cerr.ErrCodeUserDisabled, "用户状态不可登录")
	}

	ok, err := crypto.CheckPassword(in.Password, user.PasswordHash)
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	if !ok {
		return nil, cerr.New(cerr.ErrCodeInvalidPassword, "密码错误")
	}

	accessToken, err := b.jwt.GenerateToken((&jwtsUserClaimsAdapter{user: &user}).toUserClaims())
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	refreshToken, err := b.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	now := time.Now()
	_ = db.Model(&model.User{}).Where("id = ?", user.ID).Updates(map[string]any{
		"last_login_at": now,
		"login_count":   gorm.Expr("login_count + 1"),
		"updated_at":    now,
	}).Error

	return &LoginOutput{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    0,
		User:         &user,
	}, nil
}

func (b *UserBiz) RefreshToken(ctx context.Context, refreshToken string) (*LoginOutput, error) {
	if refreshToken == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("refresh_token不能为空")
	}
	claims, err := b.jwt.ParseToken(refreshToken)
	if err != nil {
		return nil, err
	}

	userID := claims.Subject
	if userID == "" {
		userID = claims.UserID
	}
	return b.LoginByUserID(ctx, userID)
}

func (b *UserBiz) LoginByUserID(ctx context.Context, userID string) (*LoginOutput, error) {
	if userID == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("user_id不能为空")
	}
	var user model.User
	err := withCtx(ctx, b.db).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.New(cerr.ErrCodeUserNotFound, "用户不存在")
		}
		return nil, cerr.ErrInternal.WithError(err)
	}

	accessToken, err := b.jwt.GenerateToken((&jwtsUserClaimsAdapter{user: &user}).toUserClaims())
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	refreshToken, err := b.jwt.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return &LoginOutput{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresIn: 0, User: &user}, nil
}

func (b *UserBiz) GetByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := withCtx(ctx, b.db).Where("id = ?", userID).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.New(cerr.ErrCodeUserNotFound, "用户不存在")
		}
		return nil, cerr.ErrInternal.WithError(err)
	}
	return &user, nil
}

func (b *UserBiz) List(ctx context.Context, filter UserListFilter) (*PageResult[*model.User], error) {
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	db := withCtx(ctx, b.db).Model(&model.User{})

	if filter.Keyword != "" {
		kw := "%" + strings.TrimSpace(filter.Keyword) + "%"
		db = db.Where("username LIKE ? OR email LIKE ? OR real_name LIKE ?", kw, kw, kw)
	}
	if filter.OrganizationID != "" {
		db = db.Joins("JOIN organization_users ou ON ou.user_id = users.id AND ou.organization_id = ?", filter.OrganizationID)
	}
	if filter.Status > 0 {
		db = db.Where("status = ?", filter.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	var users []*model.User
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return &PageResult[*model.User]{
		Items:      users,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

func (b *UserBiz) UpdateProfile(ctx context.Context, userID, email, phone, realName string) (*model.User, error) {
	if userID == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("user_id不能为空")
	}
	updates := map[string]any{"updated_at": time.Now()}
	if email != "" {
		updates["email"] = email
	}
	if phone != "" {
		updates["phone"] = phone
	}
	if realName != "" {
		updates["real_name"] = realName
	}

	err := withCtx(ctx, b.db).Model(&model.User{}).Where("id = ?", userID).Updates(updates).Error
	if err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	return b.GetByID(ctx, userID)
}

func (b *UserBiz) UpdatePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	user, err := b.GetByID(ctx, userID)
	if err != nil {
		return err
	}

	ok, err := crypto.CheckPassword(oldPassword, user.PasswordHash)
	if err != nil {
		return cerr.ErrInternal.WithError(err)
	}
	if !ok {
		return cerr.New(cerr.ErrCodeInvalidPassword, "旧密码错误")
	}
	if err := crypto.ValidatePassword(newPassword); err != nil {
		return err
	}
	newHash, err := crypto.HashPassword(newPassword)
	if err != nil {
		return cerr.ErrInternal.WithError(err)
	}

	return withCtx(ctx, b.db).Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{"password_hash": newHash, "updated_at": time.Now()}).Error
}

func (b *UserBiz) Delete(ctx context.Context, userID string) error {
	if userID == "" {
		return cerr.ErrInvalidParam.SetMessage("user_id不能为空")
	}
	return withCtx(ctx, b.db).Where("id = ?", userID).Delete(&model.User{}).Error
}

type jwtsUserClaimsAdapter struct {
	user *model.User
}

func (a *jwtsUserClaimsAdapter) toUserClaims() *jwts.UserClaims {
	return &jwts.UserClaims{
		UserID:   a.user.ID,
		Username: a.user.Username,
		Email:    a.user.Email,
		UserType: a.user.UserType,
		Status:   a.user.Status,
	}
}
