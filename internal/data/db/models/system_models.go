package models

import (
	"time"
)

// =====================================================
// 审计日志模块 Models
// =====================================================

// AuditModule 审计模块
type AuditModule string

const (
	AuditModuleLogin        AuditModule = "login"        // 登录
	AuditModuleCertificate  AuditModule = "certificate"  // 证书
	AuditModuleVerification AuditModule = "verification" // 验证
	AuditModuleSystem       AuditModule = "system"       // 系统
	AuditModuleAPI          AuditModule = "api"          // API
)

// AuditLog 审计日志表
type AuditLog struct {
	ID             string      `gorm:"type:varchar(36);primaryKey;comment:日志ID (UUID)" json:"id"`
	UserID         *string     `gorm:"type:varchar(36);index;comment:操作用户ID" json:"user_id"`
	Username       *string     `gorm:"type:varchar(64);comment:用户名" json:"username"`
	OrganizationID *string     `gorm:"type:varchar(36);comment:所属组织ID" json:"organization_id"`
	Module         AuditModule `gorm:"type:varchar(64);not null;index;comment:模块名称: login, certificate, verification, system, api" json:"module"`
	Action         string      `gorm:"type:varchar(64);not null;index;comment:操作类型" json:"action"`
	ResourceType   *string     `gorm:"type:varchar(64);comment:资源类型" json:"resource_type"`
	ResourceID     *string     `gorm:"type:varchar(36);comment:资源ID" json:"resource_id"`
	ResourceName   *string     `gorm:"type:varchar(256);comment:资源名称" json:"resource_name"`
	Description    *string     `gorm:"type:text;comment:操作描述" json:"description"`
	RequestData    *JSON       `gorm:"type:json;comment:请求数据" json:"request_data"`
	ResponseData   *JSON       `gorm:"type:json;comment:响应数据" json:"response_data"`
	IPAddress      *string     `gorm:"type:varchar(45);comment:IP地址" json:"ip_address"`
	UserAgent      *string     `gorm:"type:text;comment:User-Agent" json:"user_agent"`
	Location       *string     `gorm:"type:varchar(128);comment:地理位置" json:"location"`
	DeviceInfo     *JSON       `gorm:"type:json;comment:设备信息" json:"device_info"`
	LoginSuccess   bool        `gorm:"type:tinyint(1);default:1;comment:是否成功" json:"login_success"`
	ErrorMessage   *string     `gorm:"type:text;comment:错误信息" json:"error_message"`
	TraceID        *string     `gorm:"type:varchar(36);index;comment:链路追踪ID" json:"trace_id"`
	DurationMs     *int        `gorm:"type:int;comment:耗时 (毫秒)" json:"duration_ms"`
	CreatedAt      time.Time   `gorm:"type:datetime(3);autoCreateTime;index;comment:创建时间" json:"created_at"`
}

func (AuditLog) TableName() string {
	return "audit_logs"
}

// =====================================================
// 系统管理模块 Models
// =====================================================

// ConfigType 配置类型
type ConfigType string

const (
	ConfigTypeString  ConfigType = "string"
	ConfigTypeInteger ConfigType = "integer"
	ConfigTypeBoolean ConfigType = "boolean"
	ConfigTypeJSON    ConfigType = "json"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	ID          string     `gorm:"type:varchar(36);primaryKey;comment:配置ID (UUID)" json:"id"`
	Key         string     `gorm:"type:varchar(128);uniqueIndex;not null;comment:配置键 (唯一)" json:"key"`
	Value       string     `gorm:"type:text;not null;comment:配置值" json:"value"`
	Type        ConfigType `gorm:"type:varchar(32);default:string;comment:配置类型: string, integer, boolean, json" json:"type"`
	Description *string    `gorm:"type:text;comment:配置描述" json:"description"`
	Category    string     `gorm:"type:varchar(64);default:general;comment:配置分类" json:"category"`
	IsEncrypted bool       `gorm:"type:tinyint(1);default:0;comment:是否加密存储" json:"is_encrypted"`
	IsEditable  bool       `gorm:"type:tinyint(1);default:1;comment:是否可编辑" json:"is_editable"`
	SortOrder   int        `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	CreatedAt   time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

// Dictionary 字典表
type Dictionary struct {
	ID          string    `gorm:"type:varchar(36);primaryKey;comment:字典ID (UUID)" json:"id"`
	TypeCode    string    `gorm:"type:varchar(64);not null;index;comment:字典类型代码" json:"type_code"`
	Code        string    `gorm:"type:varchar(64);not null;comment:字典项代码" json:"code"`
	Name        string    `gorm:"type:varchar(128);not null;comment:字典项名称" json:"name"`
	Value       *string   `gorm:"type:text;comment:字典项值" json:"value"`
	ParentCode  *string   `gorm:"type:varchar(64);comment:父级代码" json:"parent_code"`
	Description *string   `gorm:"type:text;comment:描述" json:"description"`
	SortOrder   int       `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	IsActive    bool      `gorm:"type:tinyint(1);default:1;comment:是否启用" json:"is_active"`
	ExtraData   *JSON     `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt   time.Time `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (Dictionary) TableName() string {
	return "dictionaries"
}

// FileType 文件类型
type FileType string

const (
	FileTypePDF   FileType = "pdf"
	FileTypeImage FileType = "image"
	FileTypeExcel FileType = "excel"
	FileTypeWord  FileType = "word"
	FileTypeZIP   FileType = "zip"
)

// StorageType 存储类型
type StorageType int

const (
	StorageTypeLocal StorageType = 1 // 本地
	StorageTypeOSS   StorageType = 2 // OSS
	StorageTypeS3    StorageType = 3 // S3
	StorageTypeMinIO StorageType = 4 // MinIO
)

// FileRecord 文件记录表
type FileRecord struct {
	ID             string      `gorm:"type:varchar(36);primaryKey;comment:文件ID (UUID)" json:"id"`
	FileName       string      `gorm:"type:varchar(256);not null;comment:文件名" json:"file_name"`
	FilePath       string      `gorm:"type:varchar(512);not null;comment:文件路径" json:"file_path"`
	FileURL        string      `gorm:"type:varchar(512);not null;comment:文件访问URL" json:"file_url"`
	FileType       FileType    `gorm:"type:varchar(32);not null;index;comment:文件类型: pdf, image, excel, word, zip" json:"file_type"`
	MimeType       *string     `gorm:"type:varchar(64);comment:MIME类型" json:"mime_type"`
	FileSize       int64       `gorm:"type:bigint;not null;comment:文件大小 (字节)" json:"file_size"`
	FileHash       *string     `gorm:"type:varchar(64);index;comment:文件哈希值 (SHA256)" json:"file_hash"`
	StorageType    StorageType `gorm:"type:tinyint;default:1;comment:存储类型: 1-本地, 2-OSS, 3-S3, 4-MinIO" json:"storage_type"`
	StorageBucket  *string     `gorm:"type:varchar(128);comment:存储桶名称" json:"storage_bucket"`
	StorageKey     *string     `gorm:"type:varchar(256);comment:存储键" json:"storage_key"`
	UploadedBy     *string     `gorm:"type:varchar(36);index;comment:上传人ID" json:"uploaded_by"`
	OrganizationID *string     `gorm:"type:varchar(36);comment:组织ID" json:"organization_id"`
	RelatedType    *string     `gorm:"type:varchar(64);index;comment:关联业务类型" json:"related_type"`
	RelatedID      *string     `gorm:"type:varchar(36);index;comment:关联业务ID" json:"related_id"`
	Description    *string     `gorm:"type:varchar(512);comment:文件描述" json:"description"`
	IsActive       bool        `gorm:"type:tinyint(1);default:1;comment:是否有效" json:"is_active"`
	IsTemp         bool        `gorm:"type:tinyint(1);default:0;comment:是否临时文件" json:"is_temp"`
	ExpiresAt      *time.Time  `gorm:"type:datetime(3);comment:过期时间" json:"expires_at"`
	AccessCount    int         `gorm:"type:int;default:0;comment:访问次数" json:"access_count"`
	CreatedAt      time.Time   `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time   `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (FileRecord) TableName() string {
	return "file_records"
}

// JobStatus 任务状态
type JobStatus int

const (
	JobStatusWaiting  JobStatus = 1 // 等待
	JobStatusRunning  JobStatus = 2 // 执行中
	JobStatusSuccess  JobStatus = 3 // 成功
	JobStatusFailed   JobStatus = 4 // 失败
	JobStatusCanceled JobStatus = 5 // 已取消
	JobStatusRetrying JobStatus = 6 // 已重试
)

// JobQueue 任务队列表
type JobQueue struct {
	ID              string     `gorm:"type:varchar(36);primaryKey;comment:任务ID (UUID)" json:"id"`
	JobType         string     `gorm:"type:varchar(64);not null;index;comment:任务类型" json:"job_type"`
	JobName         string     `gorm:"type:varchar(128);not null;comment:任务名称" json:"job_name"`
	Payload         JSON       `gorm:"type:json;not null;comment:任务数据 (JSON)" json:"payload"`
	Priority        int8       `gorm:"type:tinyint;default:5;comment:优先级: 1-最高, 10-最低" json:"priority"`
	Status          JobStatus  `gorm:"type:tinyint;default:1;index;comment:任务状态: 1-等待, 2-执行中, 3-成功, 4-失败, 5-已取消, 6-已重试" json:"status"`
	RetryCount      int8       `gorm:"type:tinyint;default:0;comment:已重试次数" json:"retry_count"`
	MaxRetryCount   int8       `gorm:"type:tinyint;default:3;comment:最大重试次数" json:"max_retry_count"`
	NextRetryAt     *time.Time `gorm:"type:datetime(3);comment:下次重试时间" json:"next_retry_at"`
	StartedAt       *time.Time `gorm:"type:datetime(3);comment:开始执行时间" json:"started_at"`
	CompletedAt     *time.Time `gorm:"type:datetime(3);comment:完成时间" json:"completed_at"`
	ProgressPercent int8       `gorm:"type:tinyint;default:0;comment:进度百分比 (0-100)" json:"progress_percent"`
	ResultData      *JSON      `gorm:"type:json;comment:执行结果" json:"result_data"`
	ErrorMessage    *string    `gorm:"type:text;comment:错误信息" json:"error_message"`
	ScheduledAt     *time.Time `gorm:"type:datetime(3);comment:计划执行时间" json:"scheduled_at"`
	CreatedAt       time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (JobQueue) TableName() string {
	return "job_queues"
}
