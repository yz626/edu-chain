package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// JSON JSON类型字段
type JSON map[string]interface{}

// Value 实现driver.Valuer接口
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return json.Marshal(j)
}

// Scan 实现Scanner接口
func (j *JSON) Scan(value interface{}) error {
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

// =====================================================
// 用户与权限模块 Models
// =====================================================

// UserStatus 用户状态
type UserStatus int

const (
	UserStatusNormal   UserStatus = 1 // 正常
	UserStatusDisabled UserStatus = 2 // 禁用
	UserStatusPending  UserStatus = 3 // 待审核
	UserStatusLocked   UserStatus = 4 // 锁定
)

// UserType 用户类型
type UserType int

const (
	UserTypeNormal UserType = 1 // 普通用户
	UserTypeAdmin  UserType = 2 // 管理员
	UserTypeSystem UserType = 3 // 系统用户
)

// UserSource 用户来源
type UserSource int

const (
	UserSourceRegister   UserSource = 1 // 注册
	UserSourceImport     UserSource = 2 // 导入
	UserSourceThirdParty UserSource = 3 // 第三方
	UserSourceAPI        UserSource = 4 // API
)

// Gender 性别
type Gender int

const (
	GenderUnknown Gender = 0 // 未知
	GenderMale    Gender = 1 // 男
	GenderFemale  Gender = 2 // 女
)

// User 用户基础表
type User struct {
	ID               string     `gorm:"type:varchar(36);primaryKey;comment:用户ID (UUID)" json:"id"`
	Username         string     `gorm:"type:varchar(64);uniqueIndex;not null;comment:用户名" json:"username"`
	Email            string     `gorm:"type:varchar(128);uniqueIndex;not null;comment:邮箱地址" json:"email"`
	Phone            *string    `gorm:"type:varchar(20);uniqueIndex;comment:手机号" json:"phone"`
	PasswordHash     string     `gorm:"type:varchar(255);not null;comment:加密后的密码哈希" json:"-"`
	Salt             *string    `gorm:"type:varchar(64);comment:密码盐值" json:"-"`
	Status           UserStatus `gorm:"type:tinyint;default:1;comment:用户状态: 1-正常, 2-禁用, 3-待审核, 4-锁定" json:"status"`
	UserType         UserType   `gorm:"type:tinyint;default:1;comment:用户类型: 1-普通用户, 2-管理员, 3-系统用户" json:"user_type"`
	Source           UserSource `gorm:"type:tinyint;default:1;comment:注册来源: 1-注册, 2-导入, 3-第三方, 4-API" json:"source"`
	RealName         *string    `gorm:"type:varchar(64);comment:真实姓名" json:"real_name"`
	Nickname         *string    `gorm:"type:varchar(64);comment:昵称" json:"nickname"`
	AvatarURL        *string    `gorm:"type:varchar(512);comment:头像URL" json:"avatar_url"`
	Gender           *Gender    `gorm:"type:tinyint;comment:性别: 0-未知, 1-男, 2-女" json:"gender"`
	BirthDate        *time.Time `gorm:"type:date;comment:出生日期" json:"birth_date"`
	IDCardNumber     *string    `gorm:"type:varchar(18);comment:身份证号" json:"id_card_number"`
	Bio              *string    `gorm:"type:text;comment:个人简介" json:"bio"`
	Address          *string    `gorm:"type:text;comment:地址" json:"address"`
	LastLoginAt      *time.Time `gorm:"type:datetime(3);comment:最后登录时间" json:"last_login_at"`
	LastLoginIP      *string    `gorm:"type:varchar(45);comment:最后登录IP" json:"last_login_ip"`
	LoginCount       int        `gorm:"type:int;default:0;comment:登录次数" json:"login_count"`
	FailedAttempts   int8       `gorm:"type:tinyint;default:0;comment:失败尝试次数" json:"failed_attempts"`
	LockedUntil      *time.Time `gorm:"type:datetime(3);comment:锁定截止时间" json:"locked_until"`
	EmailVerified    bool       `gorm:"type:tinyint(1);default:0;comment:邮箱是否已验证" json:"email_verified"`
	PhoneVerified    bool       `gorm:"type:tinyint(1);default:0;comment:手机是否已验证" json:"phone_verified"`
	TwoFactorEnabled bool       `gorm:"type:tinyint(1);default:0;comment:是否启用双因素认证" json:"two_factor_enabled"`
	TwoFactorSecret  *string    `gorm:"type:varchar(128);comment:双因素认证密钥" json:"-"`
	ExtraData        *JSON      `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt        time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt        time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt        *time.Time `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`
	Version          int        `gorm:"type:int;default:1;comment:版本号 (乐观锁)" json:"version"`
}

func (User) TableName() string {
	return "users"
}

// RoleType 角色类型
type RoleType int

const (
	RoleTypeSuperAdmin   RoleType = 1 // 超级管理员
	RoleTypeAuditor      RoleType = 2 // 审计员
	RoleTypeOperator     RoleType = 3 // 操作员
	RoleTypeSchoolAdmin  RoleType = 4 // 校管理员
	RoleTypeDeptAdmin    RoleType = 5 // 院系管理员
	RoleTypeEnterpriseHR RoleType = 6 // 企业HR
	RoleTypeStudent      RoleType = 7 // 学生
)

// Role 角色表
type Role struct {
	ID          string     `gorm:"type:varchar(36);primaryKey;comment:角色ID (UUID)" json:"id"`
	Name        string     `gorm:"type:varchar(64);not null;comment:角色名称" json:"name"`
	Code        string     `gorm:"type:varchar(64);uniqueIndex;not null;comment:角色代码 (唯一标识)" json:"code"`
	Type        RoleType   `gorm:"type:tinyint;default:1;comment:角色类型: 1-超级管理员, 2-审计员, 3-操作员, 4-校管理员, 5-院系管理员, 6-企业HR, 7-学生" json:"type"`
	Description *string    `gorm:"type:text;comment:角色描述" json:"description"`
	Level       int        `gorm:"type:int;default:100;comment:角色级别 (数值越小级别越高)" json:"level"`
	IsSystem    bool       `gorm:"type:tinyint(1);default:0;comment:是否系统角色" json:"is_system"`
	IsActive    bool       `gorm:"type:tinyint(1);default:1;comment:是否启用" json:"is_active"`
	Permissions *JSON      `gorm:"type:json;comment:权限列表 (JSON数组)" json:"permissions"`
	Priority    int        `gorm:"type:int;default:0;comment:排序优先级" json:"priority"`
	CreatedAt   time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission 权限表
type Permission struct {
	ID          string     `gorm:"type:varchar(36);primaryKey;comment:权限ID (UUID)" json:"id"`
	Code        string     `gorm:"type:varchar(64);uniqueIndex;not null;comment:权限代码" json:"code"`
	Name        string     `gorm:"type:varchar(64);not null;comment:权限名称" json:"name"`
	Resource    string     `gorm:"type:varchar(64);not null;comment:资源类型" json:"resource"`
	Action      string     `gorm:"type:varchar(32);not null;comment:操作类型" json:"action"`
	Description *string    `gorm:"type:text;comment:权限描述" json:"description"`
	PathPattern *string    `gorm:"type:varchar(256);comment:API路径匹配模式" json:"path_pattern"`
	Method      *string    `gorm:"type:varchar(16);comment:HTTP方法" json:"method"`
	IsAPI       bool       `gorm:"type:tinyint(1);default:1;comment:是否为API权限" json:"is_api"`
	IsMenu      bool       `gorm:"type:tinyint(1);default:0;comment:是否菜单权限" json:"is_menu"`
	MenuIcon    *string    `gorm:"type:varchar(64);comment:菜单图标" json:"menu_icon"`
	MenuOrder   int        `gorm:"type:int;default:0;comment:菜单排序" json:"menu_order"`
	ParentID    *string    `gorm:"type:varchar(36);comment:父级权限ID" json:"parent_id"`
	Level       int        `gorm:"type:int;default:100;comment:权限级别" json:"level"`
	IsActive    bool       `gorm:"type:tinyint(1);default:1;comment:是否启用" json:"is_active"`
	ExtraData   *JSON      `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt   time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt   *time.Time `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`
}

func (Permission) TableName() string {
	return "permissions"
}

// UserRole 用户角色关联表
type UserRole struct {
	ID        string     `gorm:"type:varchar(36);primaryKey;comment:记录ID (UUID)" json:"id"`
	UserID    string     `gorm:"type:varchar(36);not null;index;comment:用户ID" json:"user_id"`
	RoleID    string     `gorm:"type:varchar(36);not null;index;comment:角色ID" json:"role_id"`
	GrantedAt time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:授权时间" json:"granted_at"`
	GrantedBy *string    `gorm:"type:varchar(36);comment:授权人ID" json:"granted_by"`
	ExpiresAt *time.Time `gorm:"type:datetime(3);comment:过期时间" json:"expires_at"`
	IsActive  bool       `gorm:"type:tinyint(1);default:1;comment:是否有效" json:"is_active"`
	Remark    *string    `gorm:"type:varchar(256);comment:备注" json:"remark"`
	CreatedAt time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联表
type RolePermission struct {
	ID           string    `gorm:"type:varchar(36);primaryKey;comment:记录ID (UUID)" json:"id"`
	RoleID       string    `gorm:"type:varchar(36);not null;index;comment:角色ID" json:"role_id"`
	PermissionID string    `gorm:"type:varchar(36);not null;index;comment:权限ID" json:"permission_id"`
	IsGranted    bool      `gorm:"type:tinyint(1);default:1;comment:是否授予" json:"is_granted"`
	Conditions   *JSON     `gorm:"type:json;comment:授权条件" json:"conditions"`
	CreatedAt    time.Time `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// RefreshToken 令牌表
type RefreshToken struct {
	ID         string     `gorm:"type:varchar(36);primaryKey;comment:令牌ID (UUID)" json:"id"`
	UserID     string     `gorm:"type:varchar(36);not null;index;comment:用户ID" json:"user_id"`
	TokenHash  string     `gorm:"type:varchar(128);uniqueIndex;not null;comment:令牌哈希值" json:"-"`
	ExpiresAt  time.Time  `gorm:"type:datetime(3);not null;index;comment:过期时间" json:"expires_at"`
	IsRevoked  bool       `gorm:"type:tinyint(1);default:0;not null;comment:是否已撤销" json:"is_revoked"`
	RevokedAt  *time.Time `gorm:"type:datetime(3);comment:撤销时间" json:"revoked_at"`
	DeviceInfo *JSON      `gorm:"type:json;comment:设备信息" json:"device_info"`
	IPAddress  *string    `gorm:"type:varchar(45);comment:IP地址" json:"ip_address"`
	Location   *string    `gorm:"type:varchar(128);comment:地理位置" json:"location"`
	LastUsedAt *time.Time `gorm:"type:datetime(3);comment:最后使用时间" json:"last_used_at"`
	UseCount   int        `gorm:"type:int;default:0;comment:使用次数" json:"use_count"`
	CreatedAt  time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
