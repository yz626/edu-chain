package biz

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yz626/edu-chain/internal/data/repository/model"
	cerr "github.com/yz626/edu-chain/pkg/errors"
	"gorm.io/gorm"
)

// OrganizationBiz 组织业务
type OrganizationBiz struct {
	*Base
}

func NewOrganizationBiz(base *Base) *OrganizationBiz {
	return &OrganizationBiz{Base: base}
}

type CreateOrganizationInput struct {
	Name         string
	Code         string
	Type         int32
	Address      string
	Website      string
	ContactName  string
	ContactPhone string
	ContactEmail string
}

type UpdateOrganizationInput struct {
	ID           string
	Name         string
	Address      string
	Website      string
	ContactName  string
	ContactPhone string
	ContactEmail string
	LogoURL      string
	IsActive     *bool
}

type OrganizationListFilter struct {
	Page     int
	PageSize int
	Keyword  string
	Type     int32
}

func (b *OrganizationBiz) Create(ctx context.Context, in CreateOrganizationInput) (*model.Organization, error) {
	if in.Name == "" || in.Code == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("组织名称和编码不能为空")
	}
	db := withCtx(ctx, b.db)

	var exists int64
	if err := db.Model(&model.Organization{}).Where("code = ?", in.Code).Count(&exists).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	if exists > 0 {
		return nil, cerr.New(cerr.ErrCodeInvalidParam, "组织编码已存在")
	}

	now := time.Now()
	org := &model.Organization{
		ID:           uuid.NewString(),
		Name:         in.Name,
		Code:         in.Code,
		Type:         in.Type,
		Address:      in.Address,
		Website:      in.Website,
		ContactName:  in.ContactName,
		ContactPhone: in.ContactPhone,
		ContactEmail: in.ContactEmail,
		IsActive:     true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err := db.Create(org).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	return org, nil
}

func (b *OrganizationBiz) GetByID(ctx context.Context, id string) (*model.Organization, error) {
	var org model.Organization
	err := withCtx(ctx, b.db).Where("id = ?", id).First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.New(cerr.ErrCodeInvalidParam, "组织不存在")
		}
		return nil, cerr.ErrInternal.WithError(err)
	}
	return &org, nil
}

func (b *OrganizationBiz) Update(ctx context.Context, in UpdateOrganizationInput) (*model.Organization, error) {
	if in.ID == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("组织ID不能为空")
	}
	updates := map[string]any{"updated_at": time.Now()}
	if in.Name != "" {
		updates["name"] = in.Name
	}
	if in.Address != "" {
		updates["address"] = in.Address
	}
	if in.Website != "" {
		updates["website"] = in.Website
	}
	if in.ContactName != "" {
		updates["contact_name"] = in.ContactName
	}
	if in.ContactPhone != "" {
		updates["contact_phone"] = in.ContactPhone
	}
	if in.ContactEmail != "" {
		updates["contact_email"] = in.ContactEmail
	}
	if in.LogoURL != "" {
		updates["logo_url"] = in.LogoURL
	}
	if in.IsActive != nil {
		updates["is_active"] = *in.IsActive
	}

	if err := withCtx(ctx, b.db).Model(&model.Organization{}).Where("id = ?", in.ID).Updates(updates).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	return b.GetByID(ctx, in.ID)
}

func (b *OrganizationBiz) Delete(ctx context.Context, id string) error {
	if id == "" {
		return cerr.ErrInvalidParam.SetMessage("组织ID不能为空")
	}
	return withCtx(ctx, b.db).Where("id = ?", id).Delete(&model.Organization{}).Error
}

func (b *OrganizationBiz) List(ctx context.Context, filter OrganizationListFilter) (*PageResult[*model.Organization], error) {
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	db := withCtx(ctx, b.db).Model(&model.Organization{})

	if filter.Keyword != "" {
		kw := "%" + strings.TrimSpace(filter.Keyword) + "%"
		db = db.Where("name LIKE ? OR code LIKE ?", kw, kw)
	}
	if filter.Type > 0 {
		db = db.Where("type = ?", filter.Type)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	var items []*model.Organization
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return &PageResult[*model.Organization]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

func (b *OrganizationBiz) AddMember(ctx context.Context, organizationID, userID, department, position string, isPrimary bool) error {
	if organizationID == "" || userID == "" {
		return cerr.ErrInvalidParam.SetMessage("organization_id 和 user_id 不能为空")
	}
	now := time.Now()
	item := &model.OrganizationUser{
		ID:             uuid.NewString(),
		OrganizationID: organizationID,
		UserID:         userID,
		DepartmentName: department,
		Position:       position,
		IsPrimary:      isPrimary,
		JoinedAt:       now,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	return withCtx(ctx, b.db).Create(item).Error
}

func (b *OrganizationBiz) RemoveMember(ctx context.Context, organizationID, userID string) error {
	if organizationID == "" || userID == "" {
		return cerr.ErrInvalidParam.SetMessage("organization_id 和 user_id 不能为空")
	}
	return withCtx(ctx, b.db).
		Where("organization_id = ? AND user_id = ?", organizationID, userID).
		Delete(&model.OrganizationUser{}).Error
}

func (b *OrganizationBiz) ListMembers(ctx context.Context, organizationID string, page, pageSize int) (*PageResult[*model.OrganizationUser], error) {
	if organizationID == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("organization_id不能为空")
	}
	page, pageSize = normalizePage(page, pageSize)
	db := withCtx(ctx, b.db).Model(&model.OrganizationUser{}).Where("organization_id = ?", organizationID)

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	var items []*model.OrganizationUser
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return &PageResult[*model.OrganizationUser]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}
