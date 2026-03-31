package service

import (
	"context"
	"time"

	"github.com/google/wire"
	"github.com/yz626/edu-chain/internal/biz"
	"github.com/yz626/edu-chain/internal/data/repository/model"
)

// ProviderSet 服务层依赖注入集合
var ProviderSet = wire.NewSet(
	NewUserService,
	NewOrganizationService,
	NewCertificateService,
	NewVerificationService,
	NewAuditService,
)

// UserService 用户服务
type UserService struct {
	biz *biz.UserBiz
}

func NewUserService(userBiz *biz.UserBiz) *UserService {
	return &UserService{biz: userBiz}
}

func (s *UserService) Register(ctx context.Context, in biz.RegisterInput) (*model.User, error) {
	return s.biz.Register(ctx, in)
}

func (s *UserService) Login(ctx context.Context, in biz.LoginInput) (*biz.LoginOutput, error) {
	return s.biz.Login(ctx, in)
}

func (s *UserService) RefreshToken(ctx context.Context, token string) (*biz.LoginOutput, error) {
	return s.biz.RefreshToken(ctx, token)
}

func (s *UserService) GetByID(ctx context.Context, userID string) (*model.User, error) {
	return s.biz.GetByID(ctx, userID)
}

func (s *UserService) List(ctx context.Context, filter biz.UserListFilter) (*biz.PageResult[*model.User], error) {
	return s.biz.List(ctx, filter)
}

func (s *UserService) UpdateProfile(ctx context.Context, userID, email, phone, realName string) (*model.User, error) {
	return s.biz.UpdateProfile(ctx, userID, email, phone, realName)
}

func (s *UserService) UpdatePassword(ctx context.Context, userID, oldPassword, newPassword string) error {
	return s.biz.UpdatePassword(ctx, userID, oldPassword, newPassword)
}

func (s *UserService) Delete(ctx context.Context, userID string) error {
	return s.biz.Delete(ctx, userID)
}

// OrganizationService 组织服务
type OrganizationService struct {
	biz *biz.OrganizationBiz
}

func NewOrganizationService(orgBiz *biz.OrganizationBiz) *OrganizationService {
	return &OrganizationService{biz: orgBiz}
}

func (s *OrganizationService) Create(ctx context.Context, in biz.CreateOrganizationInput) (*model.Organization, error) {
	return s.biz.Create(ctx, in)
}

func (s *OrganizationService) Update(ctx context.Context, in biz.UpdateOrganizationInput) (*model.Organization, error) {
	return s.biz.Update(ctx, in)
}

func (s *OrganizationService) GetByID(ctx context.Context, id string) (*model.Organization, error) {
	return s.biz.GetByID(ctx, id)
}

func (s *OrganizationService) Delete(ctx context.Context, id string) error {
	return s.biz.Delete(ctx, id)
}

func (s *OrganizationService) List(ctx context.Context, filter biz.OrganizationListFilter) (*biz.PageResult[*model.Organization], error) {
	return s.biz.List(ctx, filter)
}

func (s *OrganizationService) AddMember(ctx context.Context, organizationID, userID, department, position string, isPrimary bool) error {
	return s.biz.AddMember(ctx, organizationID, userID, department, position, isPrimary)
}

func (s *OrganizationService) RemoveMember(ctx context.Context, organizationID, userID string) error {
	return s.biz.RemoveMember(ctx, organizationID, userID)
}

func (s *OrganizationService) ListMembers(ctx context.Context, organizationID string, page, pageSize int) (*biz.PageResult[*model.OrganizationUser], error) {
	return s.biz.ListMembers(ctx, organizationID, page, pageSize)
}

// CertificateService 证书服务
type CertificateService struct {
	biz *biz.CertificateBiz
}

func NewCertificateService(certBiz *biz.CertificateBiz) *CertificateService {
	return &CertificateService{biz: certBiz}
}

func (s *CertificateService) Create(ctx context.Context, in biz.CreateCertificateInput) (*model.Certificate, error) {
	return s.biz.Create(ctx, in)
}

func (s *CertificateService) GetByID(ctx context.Context, id string) (*model.Certificate, error) {
	return s.biz.GetByID(ctx, id)
}

func (s *CertificateService) GetByNo(ctx context.Context, certNo string) (*model.Certificate, error) {
	return s.biz.GetByNo(ctx, certNo)
}

func (s *CertificateService) List(ctx context.Context, filter biz.CertificateListFilter) (*biz.PageResult[*model.Certificate], error) {
	return s.biz.List(ctx, filter)
}

func (s *CertificateService) Revoke(ctx context.Context, certID, revokedBy, reason string) error {
	return s.biz.Revoke(ctx, certID, revokedBy, reason)
}

func (s *CertificateService) MarkOnChainSuccess(ctx context.Context, certID, txHash, certHash string, blockNo int64, blockTime int64) error {
	return s.biz.MarkOnChainSuccess(ctx, certID, txHash, certHash, blockNo, timeFromUnix(blockTime))
}

// VerificationService 验证服务
type VerificationService struct {
	biz *biz.VerificationBiz
}

func NewVerificationService(verificationBiz *biz.VerificationBiz) *VerificationService {
	return &VerificationService{biz: verificationBiz}
}

func (s *VerificationService) Verify(ctx context.Context, in biz.VerifyInput) (*model.Verification, *model.Certificate, error) {
	return s.biz.Verify(ctx, in)
}

func (s *VerificationService) GetByID(ctx context.Context, id string) (*model.Verification, error) {
	return s.biz.GetByID(ctx, id)
}

func (s *VerificationService) List(ctx context.Context, filter biz.VerificationListFilter) (*biz.PageResult[*model.Verification], error) {
	return s.biz.List(ctx, filter)
}

// AuditService 审计服务
type AuditService struct {
	biz *biz.AuditBiz
}

func NewAuditService(auditBiz *biz.AuditBiz) *AuditService {
	return &AuditService{biz: auditBiz}
}

func (s *AuditService) CreateLog(ctx context.Context, in biz.CreateAuditLogInput) error {
	return s.biz.CreateLog(ctx, in)
}

func (s *AuditService) ListLogs(ctx context.Context, filter biz.AuditLogFilter) (*biz.PageResult[*model.AuditLog], error) {
	return s.biz.ListLogs(ctx, filter)
}

func timeFromUnix(ts int64) time.Time {
	if ts <= 0 {
		return time.Time{}
	}
	return time.Unix(ts, 0)
}
