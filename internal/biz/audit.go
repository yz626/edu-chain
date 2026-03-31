package biz

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/yz626/edu-chain/internal/data/repository/model"
	cerr "github.com/yz626/edu-chain/pkg/errors"
)

// AuditBiz 审计业务
type AuditBiz struct {
	*Base
}

func NewAuditBiz(base *Base) *AuditBiz {
	return &AuditBiz{Base: base}
}

type CreateAuditLogInput struct {
	UserID         string
	Username       string
	OrganizationID string
	Module         string
	Action         string
	ResourceType   string
	ResourceID     string
	ResourceName   string
	Description    string
	RequestData    string
	ResponseData   string
	IPAddress      string
	UserAgent      string
	Location       string
	DeviceInfo     string
	Success        bool
	ErrorMessage   string
	TraceID        string
	DurationMs     int64
}

type AuditLogFilter struct {
	Page      int
	PageSize  int
	UserID    string
	Module    string
	Action    string
	Success   *bool
	StartTime *time.Time
	EndTime   *time.Time
}

func (b *AuditBiz) CreateLog(ctx context.Context, in CreateAuditLogInput) error {
	if strings.TrimSpace(in.Module) == "" || strings.TrimSpace(in.Action) == "" {
		return cerr.ErrInvalidParam.SetMessage("module 和 action 不能为空")
	}
	log := &model.AuditLog{
		ID:             uuid.NewString(),
		UserID:         in.UserID,
		Username:       in.Username,
		OrganizationID: in.OrganizationID,
		Module:         in.Module,
		Action:         in.Action,
		ResourceType:   in.ResourceType,
		ResourceID:     in.ResourceID,
		ResourceName:   in.ResourceName,
		Description:    in.Description,
		RequestData:    in.RequestData,
		ResponseData:   in.ResponseData,
		IPAddress:      in.IPAddress,
		UserAgent:      in.UserAgent,
		Location:       in.Location,
		DeviceInfo:     in.DeviceInfo,
		LoginSuccess:   in.Success,
		ErrorMessage:   in.ErrorMessage,
		TraceID:        in.TraceID,
		DurationMs:     in.DurationMs,
		CreatedAt:      time.Now(),
	}
	return withCtx(ctx, b.db).Create(log).Error
}

func (b *AuditBiz) ListLogs(ctx context.Context, filter AuditLogFilter) (*PageResult[*model.AuditLog], error) {
	page, pageSize := normalizePage(filter.Page, filter.PageSize)
	db := withCtx(ctx, b.db).Model(&model.AuditLog{})

	if filter.UserID != "" {
		db = db.Where("user_id = ?", filter.UserID)
	}
	if filter.Module != "" {
		db = db.Where("module = ?", filter.Module)
	}
	if filter.Action != "" {
		db = db.Where("action = ?", filter.Action)
	}
	if filter.Success != nil {
		db = db.Where("login_success = ?", *filter.Success)
	}
	if filter.StartTime != nil {
		db = db.Where("created_at >= ?", *filter.StartTime)
	}
	if filter.EndTime != nil {
		db = db.Where("created_at <= ?", *filter.EndTime)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	var items []*model.AuditLog
	offset := (page - 1) * pageSize
	if err := db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&items).Error; err != nil {
		return nil, cerr.ErrInternal.WithError(err)
	}

	return &PageResult[*model.AuditLog]{
		Items:      items,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: calcTotalPages(total, pageSize),
	}, nil
}
