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

// VerificationBiz 验证业务
type VerificationBiz struct {
	*Base
}

func NewVerificationBiz(base *Base) *VerificationBiz {
	return &VerificationBiz{Base: base}
}

type VerifyInput struct {
	CertificateNo    string
	Name             string
	IDCardNumber     string
	VerificationType int32
	InputType        int32
	Purpose          string
	VerifierID       string
	VerifierOrgID    string
	IP               string
	UserAgent        string
}

type VerificationListFilter struct {
	Page          int
	PageSize      int
	CertificateID string
	UserID        string
	Result        int32
}

func (b *VerificationBiz) Verify(ctx context.Context, in VerifyInput) (*model.Verification, *model.Certificate, error) {
	if in.CertificateNo == "" {
		return nil, nil, cerr.ErrInvalidParam.SetMessage("certificate_no不能为空")
	}

	db := withCtx(ctx, b.db)
	var cert model.Certificate
	err := db.Where("certificate_no = ?", in.CertificateNo).First(&cert).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil, cerr.New(cerr.ErrCodeCertNotFound, "证书不存在")
		}
		return nil, nil, cerr.ErrInternal.WithError(err)
	}

	result := int32(1) // 真实
	if cert.Status == 2 {
		result = 4 // 已撤销
	} else {
		if in.Name != "" && !strings.EqualFold(strings.TrimSpace(in.Name), strings.TrimSpace(cert.Name)) {
			result = 3 // 未匹配
		}
		if in.IDCardNumber != "" && cert.IDCardNumber != "" && in.IDCardNumber != cert.IDCardNumber {
			result = 3
		}
	}

	now := time.Now()
	verify := &model.Verification{
		ID:               uuid.NewString(),
		VerificationNo:   now.Format("20060102150405") + "-" + uuid.NewString()[:8],
		CertificateID:    cert.ID,
		VerifierID:       in.VerifierID,
		VerifierOrgID:    in.VerifierOrgID,
		VerificationType: in.VerificationType,
		Purpose:          in.Purpose,
		InputType:        in.InputType,
		Result:           result,
		RiskLevel:        1,
		Status:           3,
		IPAddress:        in.IP,
		UserAgent:        in.UserAgent,
		VerifiedAt:       now,
		CreatedAt:        now,
		UpdatedAt:        now,
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(verify).Error; err != nil {
			return err
		}
		if err := tx.Model(&model.Certificate{}).Where("id = ?", cert.ID).Updates(map[string]any{
			"verification_count": gorm.Expr("verification_count + 1"),
			"last_verified_at":   now,
			"updated_at":         now,
		}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, nil, cerr.ErrInternal.WithError(err)
	}

	return verify, &cert, nil
}

func (b *VerificationBiz) GetByID(ctx context.Context, id string) (*model.Verification, error) {
	var item model.Verification
	err := withCtx(ctx, b.db).Where("id = ?", id).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.New(cerr.ErrCodeInvalidParam, "验证记录不存在")
		}
		return nil, cerr.ErrInternal.WithError(err)
	}
	return &item, nil
}

func (b *VerificationBiz) List(ctx context.Context, filter VerificationListFilter) (*PageResult[*model.Verification], error) {
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	db := withCtx(ctx, b.db).Model(&model.Verification{})

	if filter.CertificateID != "" {
		db = db.Where("certificate_id = ?", filter.CertificateID)
	}
	if filter.UserID != "" {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Result > 0 {
		db = db.Where("result = ?", filter.Result)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	var items []*model.Verification
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return &PageResult[*model.Verification]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}
