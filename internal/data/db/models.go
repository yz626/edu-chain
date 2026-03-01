// Package db 提供数据库模型定义
// 这些模型对应 sql/migrations/ 中的表结构
// GORM 标签说明: https://gorm.io/docs/models.html
package db

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ==================== 基础类型定义 ====================

// JSONMap 用于存储 JSON 数据的类型
type JSONMap map[string]interface{}

// Scan 实现 Scanner 接口
func (j *JSONMap) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 Valuer 接口
func (j JSONMap) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// JSONStringArray 用于存储 JSON 字符串数组
type JSONStringArray []string

// Scan 实现 Scanner 接口
func (j *JSONStringArray) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 Valuer 接口
func (j JSONStringArray) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// ==================== 用户认证模块 ====================

// User 用户基础表
type User struct {
	ID               string     `gorm:"primaryKey;type:varchar(36);comment:用户ID (UUID)" json:"id"`
	Username         string     `gorm:"type:varchar(64);uniqueIndex:uk_users_username;not null;comment:用户名" json:"username"`
	Email            string     `gorm:"type:varchar(128);uniqueIndex:uk_users_email;not null;comment:邮箱地址" json:"email"`
	Phone            string     `gorm:"type:varchar(20);uniqueIndex:uk_users_phone;comment:手机号" json:"phone"`
	PasswordHash     string     `gorm:"type:varchar(255);not null;comment:加密后的密码哈希" json:"-"`
	Salt             string     `gorm:"type:varchar(64);comment:密码盐值" json:"-"`
	Status           int8       `gorm:"type:tinyint;default:1;comment:用户状态: 1-正常, 2-禁用, 3-待审核, 4-锁定" json:"status"`
	UserType         int8       `gorm:"type:tinyint;default:1;comment:用户类型: 1-普通用户, 2-管理员, 3-系统用户" json:"user_type"`
	Source           int8       `gorm:"type:tinyint;default:1;comment:注册来源: 1-注册, 2-导入, 3-第三方, 4-API" json:"source"`
	LastLoginAt      *time.Time `gorm:"type:datetime(3);comment:最后登录时间" json:"last_login_at"`
	LastLoginIP      string     `gorm:"type:varchar(45);comment:最后登录IP" json:"last_login_ip"`
	LoginCount       int        `gorm:"type:int;default:0;comment:登录次数" json:"login_count"`
	FailedAttempts   int8       `gorm:"type:tinyint;default:0;comment:失败尝试次数" json:"failed_attempts"`
	LockedUntil      *time.Time `gorm:"type:datetime(3);comment:锁定截止时间" json:"locked_until"`
	EmailVerified    bool       `gorm:"type:tinyint(1);default:false;comment:邮箱是否已验证" json:"email_verified"`
	PhoneVerified    bool       `gorm:"type:tinyint(1);default:false;comment:手机是否已验证" json:"phone_verified"`
	TwoFactorEnabled bool       `gorm:"type:tinyint(1);default:false;comment:是否启用双因素认证" json:"two_factor_enabled"`
	TwoFactorSecret  string     `gorm:"type:varchar(128);comment:双因素认证密钥" json:"-"`
	CreatedAt        time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"type:datetime(3);index:idx_users_deleted_at;comment:删除时间" json:"deleted_at"`
	Version          int        `gorm:"type:int;default:1;comment:版本号 (乐观锁)" json:"version"`

	// Relations
	Profile       *UserProfile    `gorm:"foreignKey:UserID" json:"profile,omitempty"`
	Roles         []*Role         `gorm:"many2many:user_roles;foreignKey:ID;references:ID;joinForeignKey:UserID;joinReferences:RoleID" json:"roles,omitempty"`
	Organizations []*Organization `gorm:"many2many:organization_users;foreignKey:ID;references:ID;joinForeignKey:UserID;joinReferences:OrganizationID" json:"organizations,omitempty"`
	LoginLogs     []LoginLog      `gorm:"foreignKey:UserID" json:"login_logs,omitempty"`
	RefreshTokens []RefreshToken  `gorm:"foreignKey:UserID" json:"refresh_tokens,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// UserProfile 用户扩展信息表
type UserProfile struct {
	ID               string     `gorm:"primaryKey;type:varchar(36);comment:记录ID" json:"id"`
	UserID           string     `gorm:"type:varchar(36);uniqueIndex:uk_user_profiles_user_id;not null;comment:用户ID" json:"user_id"`
	RealName         string     `gorm:"type:varchar(64);comment:真实姓名" json:"real_name"`
	Nickname         string     `gorm:"type:varchar(64);comment:昵称" json:"nickname"`
	AvatarURL        string     `gorm:"type:varchar(512);comment:头像URL" json:"avatar_url"`
	Gender           *int8      `gorm:"type:tinyint;comment:性别: 0-未知, 1-男, 2-女" json:"gender"`
	BirthDate        *time.Time `gorm:"type:date;comment:出生日期" json:"birth_date"`
	IDCardNumber     string     `gorm:"type:varchar(18);comment:身份证号" json:"id_card_number"`
	IDCardFront      string     `gorm:"type:varchar(512);comment:身份证正面图片URL" json:"id_card_front"`
	IDCardBack       string     `gorm:"type:varchar(512);comment:身份证背面图片URL" json:"id_card_back"`
	Bio              string     `gorm:"type:text;comment:个人简介" json:"bio"`
	Address          string     `gorm:"type:text;comment:地址" json:"address"`
	Hometown         string     `gorm:"type:varchar(128);comment:籍贯" json:"hometown"`
	BloodType        *int8      `gorm:"type:tinyint;comment:血型" json:"blood_type"`
	EmergencyContact string     `gorm:"type:varchar(128);comment:紧急联系人" json:"emergency_contact"`
	EmergencyPhone   string     `gorm:"type:varchar(20);comment:紧急联系电话" json:"emergency_phone"`
	ExtraData        JSONMap    `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt        time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (UserProfile) TableName() string {
	return "user_profiles"
}

// Role 角色表
type Role struct {
	ID          string     `gorm:"primaryKey;type:varchar(36);comment:角色ID" json:"id"`
	Name        string     `gorm:"type:varchar(64);not null;comment:角色名称" json:"name"`
	Code        string     `gorm:"type:varchar(64);uniqueIndex:uk_roles_code;not null;comment:角色代码" json:"code"`
	Type        int8       `gorm:"type:tinyint;not null;default:1;comment:角色类型: 1-超级管理员, 2-审计员, 3-操作员, 4-校管理员, 5-院系管理员, 6-企业HR, 7-学生" json:"type"`
	Description string     `gorm:"type:text;comment:角色描述" json:"description"`
	Level       int        `gorm:"type:int;default:100;comment:角色级别" json:"level"`
	IsSystem    bool       `gorm:"type:tinyint(1);default:false;comment:是否系统角色" json:"is_system"`
	IsActive    bool       `gorm:"type:tinyint(1);default:true;comment:是否启用" json:"is_active"`
	Permissions JSONMap    `gorm:"type:json;comment:权限列表" json:"permissions"`
	Priority    int        `gorm:"type:int;default:0;comment:排序优先级" json:"priority"`
	CreatedAt   time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`
	Version     int        `gorm:"type:int;default:1;comment:版本号" json:"version"`

	// Relations
	Users          []*User           `gorm:"many2many:user_roles;foreignKey:ID;references:ID;joinForeignKey:RoleID;joinReferences:UserID" json:"users,omitempty"`
	PermissionsRel []*RolePermission `gorm:"foreignKey:RoleID" json:"permissions_rel,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限表
type Permission struct {
	ID          string     `gorm:"primaryKey;type:varchar(36);comment:权限ID" json:"id"`
	Code        string     `gorm:"type:varchar(64);uniqueIndex:uk_permissions_code;not null;comment:权限代码" json:"code"`
	Name        string     `gorm:"type:varchar(64);not null;comment:权限名称" json:"name"`
	Resource    string     `gorm:"type:varchar(64);not null;comment:资源类型" json:"resource"`
	Action      string     `gorm:"type:varchar(32);not null;comment:操作类型" json:"action"`
	Description string     `gorm:"type:text;comment:权限描述" json:"description"`
	PathPattern string     `gorm:"type:varchar(256);comment:API路径匹配模式" json:"path_pattern"`
	Method      string     `gorm:"type:varchar(16);comment:HTTP方法" json:"method"`
	IsAPI       bool       `gorm:"type:tinyint(1);default:true;comment:是否为API权限" json:"is_api"`
	IsMenu      bool       `gorm:"type:tinyint(1);default:false;comment:是否菜单权限" json:"is_menu"`
	MenuIcon    string     `gorm:"type:varchar(64);comment:菜单图标" json:"menu_icon"`
	MenuOrder   int        `gorm:"type:int;default:0;comment:菜单排序" json:"menu_order"`
	ParentID    *string    `gorm:"type:varchar(36);comment:父级权限ID" json:"parent_id"`
	Level       int        `gorm:"type:int;default:100;comment:权限级别" json:"level"`
	IsActive    bool       `gorm:"type:tinyint(1);default:true;comment:是否启用" json:"is_active"`
	ExtraData   JSONMap    `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt   time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`

	// Relations
	RolePermissions []*RolePermission `gorm:"foreignKey:PermissionID" json:"role_permissions,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           string    `gorm:"primaryKey;type:varchar(36);comment:记录ID" json:"id"`
	RoleID       string    `gorm:"type:varchar(36);not null;comment:角色ID" json:"role_id"`
	PermissionID string    `gorm:"type:varchar(36);not null;comment:权限ID" json:"permission_id"`
	IsGranted    bool      `gorm:"type:tinyint(1);default:true;comment:是否授予" json:"is_granted"`
	Conditions   JSONMap   `gorm:"type:json;comment:授权条件" json:"conditions"`
	CreatedAt    time.Time `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`

	// Relations
	Role       *Role       `gorm:"foreignKey:RoleID" json:"role,omitempty"`
	Permission *Permission `gorm:"foreignKey:PermissionID" json:"permission,omitempty"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        string     `gorm:"primaryKey;type:varchar(36);comment:记录ID" json:"id"`
	UserID    string     `gorm:"type:varchar(36);not null;comment:用户ID" json:"user_id"`
	RoleID    string     `gorm:"type:varchar(36);not null;comment:角色ID" json:"role_id"`
	GrantedAt time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:授权时间" json:"granted_at"`
	GrantedBy *string    `gorm:"type:varchar(36);comment:授权人ID" json:"granted_by"`
	ExpiresAt *time.Time `gorm:"type:datetime(3);comment:过期时间" json:"expires_at"`
	IsActive  bool       `gorm:"type:tinyint(1);default:true;comment:是否有效" json:"is_active"`
	Remark    string     `gorm:"type:varchar(256);comment:备注" json:"remark"`
	CreatedAt time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Role *Role `gorm:"foreignKey:RoleID" json:"role,omitempty"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// LoginLog 登录日志表
type LoginLog struct {
	ID            string     `gorm:"primaryKey;type:varchar(36);comment:日志ID" json:"id"`
	UserID        *string    `gorm:"type:varchar(36);index:idx_login_logs_user;comment:用户ID" json:"user_id"`
	SessionID     *string    `gorm:"type:varchar(36);comment:会话ID" json:"session_id"`
	IPAddress     string     `gorm:"type:varchar(45);not null;comment:IP地址" json:"ip_address"`
	Port          *int       `gorm:"type:int;comment:端口号" json:"port"`
	UserAgent     *string    `gorm:"type:text;comment:User-Agent" json:"user_agent"`
	DeviceType    *string    `gorm:"type:varchar(32);comment:设备类型" json:"device_type"`
	Browser       *string    `gorm:"type:varchar(64);comment:浏览器名称" json:"browser"`
	OS            *string    `gorm:"type:varchar(64);comment:操作系统" json:"os"`
	Platform      *string    `gorm:"type:varchar(64);comment:平台信息" json:"platform"`
	Location      *string    `gorm:"type:varchar(128);comment:登录地点" json:"location"`
	Success       bool       `gorm:"type:tinyint(1);not null;default:false;comment:是否登录成功" json:"success"`
	FailureReason *string    `gorm:"type:varchar(128);comment:失败原因" json:"failure_reason"`
	LoginType     *int8      `gorm:"type:tinyint;default:1;comment:登录类型" json:"login_type"`
	CaptchaReq    bool       `gorm:"type:tinyint(1);default:false;comment:是否需要验证码" json:"captcha_required"`
	MFARequired   bool       `gorm:"type:tinyint(1);default:false;comment:是否需要MFA" json:"mfa_required"`
	MFAVerified   bool       `gorm:"type:tinyint(1);default:false;comment:MFA是否已验证" json:"mfa_verified"`
	RequestID     *string    `gorm:"type:varchar(36);comment:请求ID" json:"request_id"`
	TraceID       *string    `gorm:"type:varchar(36);comment:链路追踪ID" json:"trace_id"`
	LoginAt       time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:登录时间" json:"login_at"`
	LogoutAt      *time.Time `gorm:"type:datetime(3);comment:登出时间" json:"logout_at"`
	SessionDur    *int       `gorm:"type:int;comment:会话持续时间" json:"session_duration"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (LoginLog) TableName() string {
	return "login_logs"
}

// RefreshToken 刷新令牌表
type RefreshToken struct {
	ID             string     `gorm:"primaryKey;type:varchar(36);comment:令牌ID" json:"id"`
	UserID         string     `gorm:"type:varchar(36);not null;index:idx_refresh_tokens_user;comment:用户ID" json:"user_id"`
	TokenHash      string     `gorm:"type:varchar(128);not null;uniqueIndex:uk_refresh_tokens_token_hash;comment:令牌哈希值" json:"-"`
	TokenSignature *string    `gorm:"type:varchar(256);comment:令牌签名" json:"-"`
	ExpiresAt      time.Time  `gorm:"type:datetime(3);not null;comment:过期时间" json:"expires_at"`
	IsRevoked      bool       `gorm:"type:tinyint(1);not null;default:false;comment:是否已撤销" json:"is_revoked"`
	RevokedAt      *time.Time `gorm:"type:datetime(3);comment:撤销时间" json:"revoked_at"`
	RevokedBy      *string    `gorm:"type:varchar(36);comment:撤销人ID" json:"revoked_by"`
	RevokeReason   *string    `gorm:"type:varchar(256);comment:撤销原因" json:"revoke_reason"`
	DeviceType     *string    `gorm:"type:varchar(32);comment:设备类型" json:"device_type"`
	DeviceName     *string    `gorm:"type:varchar(128);comment:设备名称" json:"device_name"`
	IPAddress      *string    `gorm:"type:varchar(45);comment:IP地址" json:"ip_address"`
	UserAgent      *string    `gorm:"type:text;comment:User-Agent" json:"user_agent"`
	Location       *string    `gorm:"type:varchar(128);comment:地理位置" json:"location"`
	LastUsedAt     *time.Time `gorm:"type:datetime(3);comment:最后使用时间" json:"last_used_at"`
	UseCount       int        `gorm:"type:int;default:0;comment:使用次数" json:"use_count"`
	MaxUseCount    int        `gorm:"type:int;default:1;comment:最大使用次数" json:"max_use_count"`
	CreatedAt      time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`

	// Relations
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (RefreshToken) TableName() string {
	return "refresh_tokens"
}

// Captcha 验证码表
type Captcha struct {
	ID           string     `gorm:"primaryKey;type:varchar(36);comment:验证码ID" json:"id"`
	CaptchaKey   string     `gorm:"type:varchar(64);not null;uniqueIndex:uk_captchas_key;comment:验证码键" json:"captcha_key"`
	CaptchaCode  string     `gorm:"type:varchar(16);not null;comment:验证码内容" json:"-"`
	CaptchaImage *string    `gorm:"type:text;comment:验证码图片" json:"captcha_image"`
	ExpireAt     time.Time  `gorm:"type:datetime(3);not null;comment:过期时间" json:"expire_at"`
	IsUsed       bool       `gorm:"type:tinyint(1);default:false;comment:是否已使用" json:"is_used"`
	UsedAt       *time.Time `gorm:"type:datetime(3);comment:使用时间" json:"used_at"`
	UsedCount    int8       `gorm:"type:tinyint;default:0;comment:已使用次数" json:"used_count"`
	MaxUseCount  int8       `gorm:"type:tinyint;default:1;comment:最大使用次数" json:"max_use_count"`
	IPAddress    *string    `gorm:"type:varchar(45);comment:请求IP" json:"ip_address"`
	CreatedAt    time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
}

// TableName 指定表名
func (Captcha) TableName() string {
	return "captchas"
}
