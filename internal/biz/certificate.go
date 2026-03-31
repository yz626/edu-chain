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

// CertificateBiz 证书业务
type CertificateBiz struct {
	*Base
}

func NewCertificateBiz(base *Base) *CertificateBiz {
	return &CertificateBiz{Base: base}
}

type CreateCertificateInput struct {
	CertificateNo  string
	TypeID         string
	UserID         string
	OrganizationID string
	Name           string
	IDCardNumber   string
	Major          string
	Degree         int32
	IssueDate      time.Time
	IssuedBy       string
	IssueReason    string
	TemplateID     string
	GraduationDate *time.Time
}

type CertificateListFilter struct {
	Page           int
	PageSize       int
	Keyword        string
	OrganizationID string
	UserID         string
	Status         int32
}

func (b *CertificateBiz) Create(ctx context.Context, in CreateCertificateInput) (*model.Certificate, error) {
	if in.CertificateNo == "" || in.TypeID == "" || in.UserID == "" || in.OrganizationID == "" || in.Name == "" {
		return nil, cerr.ErrInvalidParam.SetMessage("证书核心字段不能为空")
	}
	db := withCtx(ctx, b.db)

	var exists int64
	if err := db.Model(&model.Certificate{}).Where("certificate_no = ?", in.CertificateNo).Count(&exists).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	if exists > 0 {
		return nil, cerr.New(cerr.ErrCodeInvalidParam, "证书编号已存在")
	}

	if in.IssueDate.IsZero() {
		in.IssueDate = time.Now()
	}
	now := time.Now()
	cert := &model.Certificate{
		ID:             uuid.NewString(),
		CertificateNo:  in.CertificateNo,
		TypeID:         in.TypeID,
		UserID:         in.UserID,
		OrganizationID: in.OrganizationID,
		TemplateID:     in.TemplateID,
		Name:           in.Name,
		IDCardNumber:   in.IDCardNumber,
		Major:          in.Major,
		Degree:         in.Degree,
		IssueDate:      in.IssueDate,
		IssuedBy:       in.IssuedBy,
		IssueReason:    in.IssueReason,
		Status:         1,
		OnChainStatus:  1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
	if in.GraduationDate != nil {
		cert.GraduationDate = *in.GraduationDate
	}

	if err := db.Create(cert).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}
	return cert, nil
}

func (b *CertificateBiz) GetByID(ctx context.Context, id string) (*model.Certificate, error) {
	var cert model.Certificate
	err := withCtx(ctx, b.db).Where("id = ?", id).First(&cert).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.New(cerr.ErrCodeCertNotFound, "证书不存在")
		}
		return nil, cerr.ErrInternal.WithError(err)
	}
	return &cert, nil
}

func (b *CertificateBiz) GetByNo(ctx context.Context, certNo string) (*model.Certificate, error) {
	var cert model.Certificate
	err := withCtx(ctx, b.db).Where("certificate_no = ?", certNo).First(&cert).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, cerr.New(cerr.ErrCodeCertNotFound, "证书不存在")
		}
		return nil, cerr.ErrInternal.WithError(err)
	}
	return &cert, nil
}

func (b *CertificateBiz) List(ctx context.Context, filter CertificateListFilter) (*PageResult[*model.Certificate], error) {
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	db := withCtx(ctx, b.db).Model(&model.Certificate{})

	if filter.Keyword != "" {
		kw := "%" + strings.TrimSpace(filter.Keyword) + "%"
		db = db.Where("certificate_no LIKE ? OR name LIKE ?", kw, kw)
	}
	if filter.OrganizationID != "" {
		db = db.Where("organization_id = ?", filter.OrganizationID)
	}
	if filter.UserID != "" {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Status > 0 {
		db = db.Where("status = ?", filter.Status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	var items []*model.Certificate
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return &PageResult[*model.Certificate]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}

func (b *CertificateBiz) Revoke(ctx context.Context, certID, revokedBy, reason string) error {
	if certID == "" {
		return cerr.ErrInvalidParam.SetMessage("证书ID不能为空")
	}
	now := time.Now()
	result := withCtx(ctx, b.db).Model(&model.Certificate{}).
		Where("id = ?", certID).
		Updates(map[string]any{
			"status":        2,
			"revoked_by":    revokedBy,
			"revoke_reason": reason,
			"revoked_at":    now,
			"updated_at":    now,
		})
	if result.Error != nil {
		return cerr.ErrInternal.WithError(result.Error)
	}
	if result.RowsAffected == 0 {
		return cerr.New(cerr.ErrCodeCertNotFound, "证书不存在")
	}
	return nil
}

func (b *CertificateBiz) MarkOnChainSuccess(ctx context.Context, certID, txHash, certHash string, blockNo int64, blockTime time.Time) error {
	if certID == "" || txHash == "" {
		return cerr.ErrInvalidParam.SetMessage("cert_id 和 tx_hash 不能为空")
	}
	now := time.Now()
	return withCtx(ctx, b.db).Model(&model.Certificate{}).
		Where("id = ?", certID).
		Updates(map[string]any{
			"on_chain_status":      3,
			"blockchain_tx_hash":   txHash,
			"blockchain_cert_hash": certHash,
			"blockchain_block_no":  blockNo,
			"blockchain_timestamp": blockTime,
			"on_chain_at":          now,
			"updated_at":           now,
		}).Error
}
