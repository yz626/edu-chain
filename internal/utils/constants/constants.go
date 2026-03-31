package constants

// =====================================================
// 用户相关常量
// =====================================================

// UserStatus 用户状态: 1-正常, 2-禁用, 3-待审核, 4-锁定
const (
	UserStatusNormal   = 1
	UserStatusDisabled = 2
	UserStatusPending  = 3
	UserStatusLocked   = 4
)

// UserType 用户类型: 1-普通用户, 2-管理员, 3-系统用户
const (
	UserTypeNormal = 1
	UserTypeAdmin  = 2
	UserTypeSystem = 3
)

// UserSource 用户来源: 1-注册, 2-导入, 3-第三方, 4-API
const (
	UserSourceRegister   = 1
	UserSourceImport     = 2
	UserSourceThirdParty = 3
	UserSourceAPI        = 4
)

// Gender 性别: 0-未知, 1-男, 2-女
const (
	GenderUnknown = 0
	GenderMale    = 1
	GenderFemale  = 2
)

// =====================================================
// 组织相关常量
// =====================================================

// OrganizationType 组织类型: 1-监管机构, 2-高校, 3-企业, 4-政府机构, 5-其他
const (
	OrganizationTypeSupervisor = 1
	OrganizationTypeUniversity = 2
	OrganizationTypeEnterprise = 3
	OrganizationTypeGovernment = 4
	OrganizationTypeOther      = 5
)

// =====================================================
// 角色相关常量
// =====================================================

// RoleType 角色类型: 1-超级管理员, 2-审计员, 3-操作员, 4-校管理员, 5-院系管理员, 6-企业HR, 7-学生
const (
	RoleTypeSuperAdmin   = 1
	RoleTypeAuditor      = 2
	RoleTypeOperator     = 3
	RoleTypeSchoolAdmin  = 4
	RoleTypeDeptAdmin    = 5
	RoleTypeEnterpriseHR = 6
	RoleTypeStudent      = 7
)

// =====================================================
// 部门相关常量
// =====================================================

// DepartmentType 部门类型: 1-职能部门, 2-教学部门, 3-行政部门, 4-其他
const (
	DepartmentTypeFunctional = 1
	DepartmentTypeTeaching   = 2
	DepartmentTypeAdmin      = 3
	DepartmentTypeOther      = 4
)

// =====================================================
// 区块链相关常量
// =====================================================

// BlockchainType 区块链类型: 1-Fabric, 2-Ethereum
const (
	BlockchainTypeFabric   = 1
	BlockchainTypeEthereum = 2
)

// NetworkStatus 网络状态: 1-正常, 2-维护中, 3-停用
const (
	NetworkStatusNormal   = 1
	NetworkStatusMaintain = 2
	NetworkStatusDisabled = 3
)

// TxType 交易类型: 1-存证, 2-撤销, 3-查询, 4-转让
const (
	TxTypeStore    = 1
	TxTypeRevoke   = 2
	TxTypeQuery    = 3
	TxTypeTransfer = 4
)

// TxStatus 交易状态: 1-待处理, 2-处理中, 3-成功, 4-失败, 5-超时
const (
	TxStatusPending  = 1
	TxStatusProgress = 2
	TxStatusSuccess  = 3
	TxStatusFailed   = 4
	TxStatusTimeout  = 5
)

// ContractStatus 合约状态: 1-正常, 2-停用, 3-已废弃
const (
	ContractStatusNormal     = 1
	ContractStatusDisabled   = 2
	ContractStatusDeprecated = 3
)

// =====================================================
// 证书相关常量
// =====================================================

// CertificateStatus 证书状态: 1-有效, 2-已撤销, 3-已挂失, 4-草稿, 5-审核中
const (
	CertificateStatusValid   = 1
	CertificateStatusRevoked = 2
	CertificateStatusLost    = 3
	CertificateStatusDraft   = 4
	CertificateStatusPending = 5
)

// CertificateOnChainStatus 证书上链状态: 1-待上链, 2-上链中, 3-已上链, 4-上链失败
const (
	OnChainStatusPending  = 1
	OnChainStatusProgress = 2
	OnChainStatusSuccess  = 3
	OnChainStatusFailed   = 4
)

// Degree 学位等级: 1-学士, 2-硕士, 3-博士, 4-无学位, 5-专科
const (
	DegreeBachelor  = 1
	DegreeMaster    = 2
	DegreeDoctor    = 3
	DegreeNone      = 4
	DegreeAssociate = 5
)

// EducationType 办学类型: 1-全日制, 2-非全日制, 3-成人教育
const (
	EducationTypeFullTime = 1
	EducationTypePartTime = 2
	EducationTypeAdult    = 3
)

// =====================================================
// 证书类型相关常量
// =====================================================

// CertificateCategory 证书类别: 1-毕业证书, 2-学位证书, 3-成绩单, 4-资格证书, 5-其他
const (
	CertificateCategoryGraduation = 1
	CertificateCategoryDegree     = 2
	CertificateCategoryTranscript = 3
	CertificateCategoryCredential = 4
	CertificateCategoryOther      = 5
)

// DegreeLevel 学位等级: 1-学士, 2-硕士, 3-博士
const (
	DegreeLevelBachelor = 1
	DegreeLevelMaster   = 2
	DegreeLevelDoctor   = 3
)

// =====================================================
// 证书批次相关常量
// =====================================================

// CertificateBatchStatus 证书批次状态: 1-待处理, 2-处理中, 3-已完成, 4-已暂停, 5-已取消
const (
	CertificateBatchStatusPending   = 1
	CertificateBatchStatusProgress  = 2
	CertificateBatchStatusCompleted = 3
	CertificateBatchStatusPaused    = 4
	CertificateBatchStatusCancelled = 5
)

// =====================================================
// 组织用户相关常量
// =====================================================

// EmploymentType 用工类型: 1-正式, 2-合同, 3-实习, 4-外包
const (
	EmploymentTypeRegular   = 1
	EmploymentTypeContract  = 2
	EmploymentTypeIntern    = 3
	EmploymentTypeOutsource = 4
)

// WorkStatus 工作状态: 1-在职, 2-离职, 3-退休, 4-停薪留职
const (
	WorkStatusActive   = 1
	WorkStatusResigned = 2
	WorkStatusRetired  = 3
	WorkStatusLeave    = 4
)

// =====================================================
// 文件相关常量
// =====================================================

// FileType 文件类型: pdf, image, excel, word, zip
const (
	FileTypePDF   = "pdf"
	FileTypeImage = "image"
	FileTypeExcel = "excel"
	FileTypeWord  = "word"
	FileTypeZip   = "zip"
)

// StorageType 存储类型: 1-本地, 2-OSS, 3-S3, 4-MinIO
const (
	StorageTypeLocal = 1
	StorageTypeOSS   = 2
	StorageTypeS3    = 3
	StorageTypeMinIO = 4
)

// =====================================================
// 任务队列相关常量
// =====================================================

// JobStatus 任务状态: 1-等待, 2-执行中, 3-成功, 4-失败, 5-已取消, 6-已重试
const (
	JobStatusWaiting   = 1
	JobStatusRunning   = 2
	JobStatusSuccess   = 3
	JobStatusFailed    = 4
	JobStatusCancelled = 5
	JobStatusRetried   = 6
)

// JobPriority 任务优先级: 1-最高, 10-最低
const (
	JobPriorityHighest = 1
	JobPriorityLowest  = 10
)

// =====================================================
// 验证相关常量
// =====================================================

// VerificationType 验证类型: 1-本人查询, 2-机构验证, 3-管理员查询, 4-批量验证
const (
	VerificationTypeSelfQuery   = 1
	VerificationTypeOrgVerify   = 2
	VerificationTypeAdminQuery  = 3
	VerificationTypeBatchVerify = 4
)

// VerificationInputType 验证输入类型: 1-证书编号, 2-身份证号, 3-姓名+身份证, 4-扫码验证
const (
	VerificationInputCertNo = 1
	VerificationInputIDCard = 2
	VerificationInputNameID = 3
	VerificationInputQRCode = 4
)

// VerificationResult 验证结果: 1-真实, 2-可疑, 3-未匹配, 4-已撤销
const (
	VerificationResultValid      = 1
	VerificationResultSuspicious = 2
	VerificationResultMismatch   = 3
	VerificationResultRevoked    = 4
)

// VerificationStatus 验证记录状态: 1-待验证, 2-验证中, 3-已完成, 4-已过期
const (
	VerificationStatusPending   = 1
	VerificationStatusProgress  = 2
	VerificationStatusCompleted = 3
	VerificationStatusExpired   = 4
)

// RiskLevel 风险等级: 1-低, 2-中, 3-高, 4-极高, 5-已确认欺诈
const (
	RiskLevelLow            = 1
	RiskLevelMedium         = 2
	RiskLevelHigh           = 3
	RiskLevelVeryHigh       = 4
	RiskLevelConfirmedFraud = 5
)

// ReviewResult 审核结果: 1-待审核, 2-已通过, 3-已拒绝
const (
	ReviewResultPending  = 1
	ReviewResultApproved = 2
	ReviewResultRejected = 3
)

// =====================================================
// 通用状态常量
// =====================================================

// ActiveStatus 启用状态: 0-禁用, 1-启用
const (
	ActiveStatusYes = 1
	ActiveStatusNo  = 0
)

// =====================================================
// 审计日志相关常量
// =====================================================

// AuditModule 审计模块: login, certificate, verification, system, api
const (
	AuditModuleLogin        = "login"
	AuditModuleCertificate  = "certificate"
	AuditModuleVerification = "verification"
	AuditModuleSystem       = "system"
	AuditModuleAPI          = "api"
)
