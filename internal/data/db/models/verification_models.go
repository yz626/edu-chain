package models

import (
	"time"
)

// =====================================================
// 验证服务模块 Models
// =====================================================

// VerificationType 验证类型
type VerificationType int

const (
	VerificationTypeSelfQuery   VerificationType = 1 // 本人查询
	VerificationTypeOrgVerify   VerificationType = 2 // 机构验证
	VerificationTypeAdminQuery  VerificationType = 3 // 管理员查询
	VerificationTypeBatchVerify VerificationType = 4 // 批量验证
)

// InputType 输入类型
type InputType int

const (
	InputTypeCertNo InputType = 1 // 证书编号
	InputTypeIDCard InputType = 2 // 身份证号
	InputTypeNameID InputType = 3 // 姓名+身份证
	InputTypeQRCode InputType = 4 // 扫码验证
)

// VerifyResult 验证结果
type VerifyResult int

const (
	VerifyResultValid      VerifyResult = 1 // 真实
	VerifyResultSuspicious VerifyResult = 2 // 可疑
	VerifyResultNotMatch   VerifyResult = 3 // 未匹配
	VerifyResultRevoked    VerifyResult = 4 // 已撤销
)

// RiskLevel 风险等级
type RiskLevel int

const (
	RiskLevelLow      RiskLevel = 1 // 低
	RiskLevelMedium   RiskLevel = 2 // 中
	RiskLevelHigh     RiskLevel = 3 // 高
	RiskLevelVeryHigh RiskLevel = 4 // 极高
	RiskLevelFraud    RiskLevel = 5 // 已确认欺诈
)

// VerifyStatus 验证状态
type VerifyStatus int

const (
	VerifyStatusPending    VerifyStatus = 1 // 待验证
	VerifyStatusProcessing VerifyStatus = 2 // 验证中
	VerifyStatusCompleted  VerifyStatus = 3 // 已完成
	VerifyStatusExpired    VerifyStatus = 4 // 已过期
)

// Verification 验证记录表
type Verification struct {
	ID              string  `gorm:"type:varchar(36);primaryKey;comment:验证记录ID (UUID)" json:"id"`
	VerificationNo  string  `gorm:"type:varchar(64);uniqueIndex;not null;comment:验证流水号 (唯一)" json:"verification_no"`
	CertificateID   string  `gorm:"type:varchar(36);not null;index;comment:证书ID" json:"certificate_id"`
	UserID          *string `gorm:"type:varchar(36);index;comment:验证人用户ID" json:"user_id"`
	VerifierID      *string `gorm:"type:varchar(36);index;comment:验证操作人ID" json:"verifier_id"`
	VerifierOrgID   *string `gorm:"type:varchar(36);comment:验证人所属组织ID" json:"verifier_org_id"`
	VerifierOrgName *string `gorm:"type:varchar(128);comment:验证人所属组织名称" json:"verifier_org_name"`

	// 验证类型与用途
	VerificationType VerificationType `gorm:"type:tinyint;default:1;index;comment:验证类型: 1-本人查询, 2-机构验证, 3-管理员查询, 4-批量验证" json:"verification_type"`
	Purpose          *string          `gorm:"type:varchar(256);comment:验证用途" json:"purpose"`
	InputType        InputType        `gorm:"type:tinyint;default:1;comment:输入类型: 1-证书编号, 2-身份证号, 3-姓名+身份证, 4-扫码验证" json:"input_type"`
	InputData        *JSON            `gorm:"type:json;comment:输入数据" json:"input_data"`

	// 验证结果
	Result         VerifyResult `gorm:"type:tinyint;not null;index;comment:验证结果: 1-真实, 2-可疑, 3-未匹配, 4-已撤销" json:"result"`
	ResultDetails  *JSON        `gorm:"type:json;comment:验证结果详情" json:"result_details"`
	MatchedFields  *JSON        `gorm:"type:json;comment:匹配的字段列表" json:"matched_fields"`
	MismatchFields *JSON        `gorm:"type:json;comment:不匹配的字段列表" json:"mismatch_fields"`

	// 区块链验证
	BlockchainVerified bool  `gorm:"type:tinyint(1);default:0;comment:区块链是否已验证" json:"blockchain_verified"`
	BlockchainResult   *JSON `gorm:"type:json;comment:区块链验证结果" json:"blockchain_result"`

	// 报告信息
	ReportURL *string `gorm:"type:varchar(512);comment:验证报告URL" json:"report_url"`
	ReportID  *string `gorm:"type:varchar(36);comment:报告ID" json:"report_id"`

	// 风险评估
	RiskLevel   RiskLevel `gorm:"type:tinyint;default:1;comment:风险等级: 1-低, 2-中, 3-高, 4-极高, 5-已确认欺诈" json:"risk_level"`
	RiskFactors *JSON     `gorm:"type:json;comment:风险因素列表" json:"risk_factors"`

	// 审核信息
	ReviewedBy   *string    `gorm:"type:varchar(36);comment:审核人ID" json:"reviewed_by"`
	ReviewedAt   *time.Time `gorm:"type:datetime(3);comment:审核时间" json:"reviewed_at"`
	ReviewResult *int8      `gorm:"type:tinyint;comment:审核结果" json:"review_result"`
	ReviewNotes  *string    `gorm:"type:text;comment:审核备注" json:"review_notes"`

	// 请求信息
	IPAddress *string `gorm:"type:varchar(45);comment:请求IP地址" json:"ip_address"`
	UserAgent *string `gorm:"type:text;comment:User-Agent" json:"user_agent"`
	Location  *string `gorm:"type:varchar(128);comment:地理位置" json:"location"`

	// 状态
	Status     VerifyStatus `gorm:"type:tinyint;default:1;comment:状态: 1-待验证, 2-验证中, 3-已完成, 4-已过期" json:"status"`
	ExpiresAt  *time.Time   `gorm:"type:datetime(3);comment:过期时间" json:"expires_at"`
	VerifiedAt *time.Time   `gorm:"type:datetime(3);comment:验证完成时间" json:"verified_at"`

	// 审计字段
	CreatedAt time.Time `gorm:"type:datetime(3);autoCreateTime;index;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (Verification) TableName() string {
	return "verifications"
}
