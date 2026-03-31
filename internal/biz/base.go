package biz

import (
	"context"

	"github.com/google/wire"
	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/internal/data/db"
	"github.com/yz626/edu-chain/internal/utils/jwts"
	"github.com/yz626/edu-chain/pkg/logger"
	"gorm.io/gorm"
)

// ProviderSet 业务层依赖注入集合
var ProviderSet = wire.NewSet(
	NewBase,
	NewUserBiz,
	NewOrganizationBiz,
	NewCertificateBiz,
	NewVerificationBiz,
	NewAuditBiz,
)

// Base 业务层公共依赖
type Base struct {
	db  *gorm.DB
	log *logger.Logger
	jwt *jwts.JWT
}

func NewBase(cfg *config.Config, database *db.DateDB, log *logger.Logger) *Base {
	return &Base{
		db:  database.DB,
		log: log.Named("biz"),
		jwt: jwts.New(&cfg.JWT),
	}
}

// PageResult 通用分页结果
type PageResult[T any] struct {
	Items      []T   `json:"items"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

func normalizePage(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func calcTotalPages(total int64, pageSize int) int {
	if total == 0 {
		return 0
	}
	return int((total + int64(pageSize) - 1) / int64(pageSize))
}

func withCtx(ctx context.Context, db *gorm.DB) *gorm.DB {
	if ctx == nil {
		return db
	}
	return db.WithContext(ctx)
}
