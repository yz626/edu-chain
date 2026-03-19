package models

import (
	"time"
)

// =====================================================
// 证书管理模块 Models
// =====================================================

// CertificateCategory 证书类别
type CertificateCategory int

const (
	CertCategoryDiploma    CertificateCategory = 1 // 毕业证书
	CertCategoryDegree     CertificateCategory = 2 // 学位证书
	CertCategoryTranscript CertificateCategory = 3 // 成绩单
	CertCategoryCredential CertificateCategory = 4 // 资格证书
	CertCategoryOther      CertificateCategory = 5 // 其他
)

// DegreeLevel 学位等级
type DegreeLevel int

const (
	DegreeLevelBachelor  DegreeLevel = 1 // 学士
	DegreeLevelMaster    DegreeLevel = 2 // 硕士
	DegreeLevelDoctor    DegreeLevel = 3 // 博士
	DegreeLevelNone      DegreeLevel = 4 // 无学位
	DegreeLevelAssociate DegreeLevel = 5 // 专科
)

// CertificateType 证书类型表
type CertificateType struct {
	ID          string              `gorm:"type:varchar(36);primaryKey;comment:类型ID (UUID)" json:"id"`
	Code        string              `gorm:"type:varchar(32);uniqueIndex;not null;comment:类型代码 (唯一标识)" json:"code"`
	Name        string              `gorm:"type:varchar(64);not null;comment:类型名称" json:"name"`
	Category    CertificateCategory `gorm:"type:tinyint;default:1;comment:证书类别: 1-毕业证书, 2-学位证书, 3-成绩单, 4-资格证书, 5-其他" json:"category"`
	DegreeLevel *DegreeLevel        `gorm:"type:tinyint;comment:学位等级: 1-学士, 2-硕士, 3-博士" json:"degree_level"`
	Description *string             `gorm:"type:text;comment:类型描述" json:"description"`
	IsActive    bool                `gorm:"type:tinyint(1);default:1;comment:是否启用" json:"is_active"`
	SortOrder   int                 `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	ExtraData   *JSON               `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt   time.Time           `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time           `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (CertificateType) TableName() string {
	return "certificate_types"
}

// CertificateTemplate 证书模板表
type CertificateTemplate struct {
	ID              string    `gorm:"type:varchar(36);primaryKey;comment:模板ID (UUID)" json:"id"`
	Name            string    `gorm:"type:varchar(128);not null;comment:模板名称" json:"name"`
	Code            string    `gorm:"type:varchar(64);uniqueIndex;not null;comment:模板代码 (唯一标识)" json:"code"`
	TypeID          string    `gorm:"type:varchar(36);not null;index;comment:证书类型ID" json:"type_id"`
	ThumbnailURL    *string   `gorm:"type:varchar(512);comment:缩略图URL" json:"thumbnail_url"`
	TemplateFileURL *string   `gorm:"type:varchar(512);comment:模板文件URL" json:"template_file_url"`
	TemplateData    *JSON     `gorm:"type:json;comment:模板配置数据" json:"template_data"`
	Fields          JSON      `gorm:"type:json;not null;comment:模板字段配置" json:"fields"`
	IsDefault       bool      `gorm:"type:tinyint(1);default:0;comment:是否默认模板" json:"is_default"`
	IsActive        bool      `gorm:"type:tinyint(1);default:1;comment:是否启用" json:"is_active"`
	Version         int       `gorm:"type:int;default:1;comment:版本号" json:"version"`
	OrganizationID  *string   `gorm:"type:varchar(36);comment:所属组织ID" json:"organization_id"`
	CreatedBy       *string   `gorm:"type:varchar(36);comment:创建人ID" json:"created_by"`
	CreatedAt       time.Time `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (CertificateTemplate) TableName() string {
	return "certificate_templates"
}

// EducationType 办学类型
type EducationType int

const (
	EducationTypeFullTime EducationType = 1 // 全日制
	EducationTypePartTime EducationType = 2 // 非全日制
	EducationTypeAdult    EducationType = 3 // 成人教育
)

// CertificateStatus 证书状态
type CertificateStatus int

const (
	CertStatusValid   CertificateStatus = 1 // 有效
	CertStatusRevoked CertificateStatus = 2 // 已撤销
	CertStatusLost    CertificateStatus = 3 // 已挂失
	CertStatusDraft   CertificateStatus = 4 // 草稿
	CertStatusPending CertificateStatus = 5 // 审核中
)

// OnChainStatus 上链状态
type OnChainStatus int

const (
	OnChainStatusPending    OnChainStatus = 1 // 待上链
	OnChainStatusInProgress OnChainStatus = 2 // 上链中
	OnChainStatusCompleted  OnChainStatus = 3 // 已上链
	OnChainStatusFailed     OnChainStatus = 4 // 上链失败
)

// Certificate 证书主表
type Certificate struct {
	ID             string  `gorm:"type:varchar(36);primaryKey;comment:证书ID (UUID)" json:"id"`
	CertificateNo  string  `gorm:"type:varchar(64);uniqueIndex;not null;comment:证书编号 (唯一)" json:"certificate_no"`
	TypeID         string  `gorm:"type:varchar(36);not null;index;comment:证书类型ID" json:"type_id"`
	UserID         string  `gorm:"type:varchar(36);not null;index;comment:持有人用户ID" json:"user_id"`
	OrganizationID string  `gorm:"type:varchar(36);not null;index;comment:发证机构ID" json:"organization_id"`
	TemplateID     *string `gorm:"type:varchar(36);comment:使用的模板ID" json:"template_id"`

	// 基本信息
	StudentNo    *string `gorm:"type:varchar(64);index;comment:学号" json:"student_no"`
	Name         string  `gorm:"type:varchar(64);not null;index;comment:证书持有人姓名" json:"name"`
	IDCardNumber *string `gorm:"type:varchar(18);comment:身份证号" json:"id_card_number"`
	Gender       *int8   `gorm:"type:tinyint;comment:性别: 1-男, 2-女" json:"gender"`

	// 学历信息
	Major          *string       `gorm:"type:varchar(128);comment:专业名称" json:"major"`
	MajorCode      *string       `gorm:"type:varchar(32);comment:专业代码" json:"major_code"`
	Degree         DegreeLevel   `gorm:"type:tinyint;default:1;comment:学位等级: 1-学士, 2-硕士, 3-博士, 4-无学位, 5-专科" json:"degree"`
	DegreeName     *string       `gorm:"type:varchar(32);comment:学位名称" json:"degree_name"`
	EducationLevel *int8         `gorm:"type:tinyint;comment:学历层次" json:"education_level"`
	EducationType  EducationType `gorm:"type:tinyint;default:1;comment:办学类型: 1-全日制, 2-非全日制, 3-成人教育" json:"education_type"`
	EnrollmentDate *time.Time    `gorm:"type:date;comment:入学日期" json:"enrollment_date"`
	GraduationDate *time.Time    `gorm:"type:date;index;comment:毕业日期" json:"graduation_date"`
	StudyPeriod    *string       `gorm:"type:varchar(64);comment:学习期限" json:"study_period"`
	SchoolSystem   *string       `gorm:"type:varchar(32);comment:学制" json:"school_system"`
	Campus         *string       `gorm:"type:varchar(128);comment:校区" json:"campus"`

	// 成绩信息
	GPA          *float64 `gorm:"type:decimal(4,2);comment:平均学分绩点" json:"gpa"`
	TotalCredits *float64 `gorm:"type:decimal(6,2);comment:总学分" json:"total_credits"`

	// 证书状态
	Status         CertificateStatus `gorm:"type:tinyint;default:1;index;comment:证书状态: 1-有效, 2-已撤销, 3-已挂失, 4-草稿, 5-审核中" json:"status"`
	IssueDate      time.Time         `gorm:"type:date;index;comment:发证日期" json:"issue_date"`
	ValidFromDate  *time.Time        `gorm:"type:date;comment:有效开始日期" json:"valid_from_date"`
	ValidUntilDate *time.Time        `gorm:"type:date;comment:有效结束日期" json:"valid_until_date"`

	// 区块链信息
	BlockchainTxHash    *string       `gorm:"type:varchar(128);index;comment:区块链交易哈希" json:"blockchain_tx_hash"`
	BlockchainCertHash  *string       `gorm:"type:varchar(64);comment:证书数据哈希" json:"blockchain_cert_hash"`
	BlockchainBlockNo   *int64        `gorm:"type:bigint;comment:区块链块号" json:"blockchain_block_no"`
	BlockchainTimestamp *time.Time    `gorm:"type:datetime(3);comment:区块链时间戳" json:"blockchain_timestamp"`
	OnChainStatus       OnChainStatus `gorm:"type:tinyint;default:1;comment:上链状态: 1-待上链, 2-上链中, 3-已上链, 4-上链失败" json:"on_chain_status"`
	OnChainAt           *time.Time    `gorm:"type:datetime(3);comment:上链时间" json:"on_chain_at"`

	// PDF信息
	PDFURL      *string    `gorm:"type:varchar(512);comment:PDF文件URL" json:"pdf_url"`
	PDFHash     *string    `gorm:"type:varchar(64);comment:PDF文件哈希" json:"pdf_hash"`
	PDFSignedBy *string    `gorm:"type:varchar(36);comment:PDF签名人ID" json:"pdf_signed_by"`
	PDFSignedAt *time.Time `gorm:"type:datetime(3);comment:PDF签名时间" json:"pdf_signed_at"`

	// 签发/撤销信息
	SignedBy     *string    `gorm:"type:varchar(36);comment:证书签发人ID" json:"signed_by"`
	IssuedBy     *string    `gorm:"type:varchar(36);comment:操作人ID" json:"issued_by"`
	IssueReason  *string    `gorm:"type:varchar(256);comment:颁发原因" json:"issue_reason"`
	RevokedAt    *time.Time `gorm:"type:datetime(3);comment:撤销时间" json:"revoked_at"`
	RevokedBy    *string    `gorm:"type:varchar(36);comment:撤销人ID" json:"revoked_by"`
	RevokeReason *string    `gorm:"type:text;comment:撤销原因" json:"revoke_reason"`

	// 验证统计
	VerificationCount int        `gorm:"type:int;default:0;comment:验证次数" json:"verification_count"`
	LastVerifiedAt    *time.Time `gorm:"type:datetime(3);comment:最后验证时间" json:"last_verified_at"`

	// 扩展数据
	ExtraData *JSON `gorm:"type:json;comment:扩展数据" json:"extra_data"`

	// 审计字段
	CreatedAt time.Time  `gorm:"type:datetime(3);autoCreateTime;index;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt *time.Time `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`
	Version   int        `gorm:"type:int;default:1;comment:版本号" json:"version"`
}

func (Certificate) TableName() string {
	return "certificates"
}

// BatchStatus 批次状态
type BatchStatus int

const (
	BatchStatusPending   BatchStatus = 1 // 待处理
	BatchStatusProgress  BatchStatus = 2 // 处理中
	BatchStatusCompleted BatchStatus = 3 // 已完成
	BatchStatusPaused    BatchStatus = 4 // 已暂停
	BatchStatusCanceled  BatchStatus = 5 // 已取消
)

// CertificateBatch 证书批次表
type CertificateBatch struct {
	ID             string      `gorm:"type:varchar(36);primaryKey;comment:批次ID (UUID)" json:"id"`
	BatchNo        string      `gorm:"type:varchar(64);uniqueIndex;not null;comment:批次号 (唯一)" json:"batch_no"`
	TypeID         string      `gorm:"type:varchar(36);not null;index;comment:证书类型ID" json:"type_id"`
	OrganizationID string      `gorm:"type:varchar(36);not null;index;comment:发证机构ID" json:"organization_id"`
	TemplateID     *string     `gorm:"type:varchar(36);comment:模板ID" json:"template_id"`
	Name           string      `gorm:"type:varchar(128);not null;comment:批次名称" json:"name"`
	Description    *string     `gorm:"type:text;comment:批次描述" json:"description"`
	TotalCount     int         `gorm:"type:int;default:0;comment:总数量" json:"total_count"`
	SuccessCount   int         `gorm:"type:int;default:0;comment:成功数量" json:"success_count"`
	FailCount      int         `gorm:"type:int;default:0;comment:失败数量" json:"fail_count"`
	Status         BatchStatus `gorm:"type:tinyint;default:1;index;comment:批次状态: 1-待处理, 2-处理中, 3-已完成, 4-已暂停, 5-已取消" json:"status"`
	ImportFileURL  *string     `gorm:"type:varchar(512);comment:导入文件URL" json:"import_file_url"`
	ExecutedBy     *string     `gorm:"type:varchar(36);comment:执行人ID" json:"executed_by"`
	StartedAt      *time.Time  `gorm:"type:datetime(3);comment:开始时间" json:"started_at"`
	CompletedAt    *time.Time  `gorm:"type:datetime(3);comment:完成时间" json:"completed_at"`
	CreatedAt      time.Time   `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (CertificateBatch) TableName() string {
	return "certificate_batches"
}
