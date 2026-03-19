package models

import (
	"time"
)

// =====================================================
// 组织架构模块 Models
// =====================================================

// OrganizationType 组织类型
type OrganizationType int

const (
	OrgTypeRegulator  OrganizationType = 1 // 监管机构
	OrgTypeUniversity OrganizationType = 2 // 高校
	OrgTypeEnterprise OrganizationType = 3 // 企业
	OrgTypeGovernment OrganizationType = 4 // 政府机构
	OrgTypeOther      OrganizationType = 5 // 其他
)

// Organization 组织表
type Organization struct {
	ID              string           `gorm:"type:varchar(36);primaryKey;comment:组织ID (UUID)" json:"id"`
	Name            string           `gorm:"type:varchar(128);not null;comment:组织名称" json:"name"`
	Code            string           `gorm:"type:varchar(64);uniqueIndex;not null;comment:组织代码 (唯一标识)" json:"code"`
	Type            OrganizationType `gorm:"type:tinyint;default:5;comment:组织类型: 1-监管机构, 2-高校, 3-企业, 4-政府机构, 5-其他" json:"type"`
	ShortName       *string          `gorm:"type:varchar(64);comment:组织简称" json:"short_name"`
	LogoURL         *string          `gorm:"type:varchar(512);comment:Logo URL" json:"logo_url"`
	Description     *string          `gorm:"type:text;comment:组织描述" json:"description"`
	Address         *string          `gorm:"type:text;comment:地址" json:"address"`
	Province        *string          `gorm:"type:varchar(32);comment:省份" json:"province"`
	City            *string          `gorm:"type:varchar(32);comment:城市" json:"city"`
	District        *string          `gorm:"type:varchar(32);comment:区县" json:"district"`
	PostalCode      *string          `gorm:"type:varchar(10);comment:邮政编码" json:"postal_code"`
	Website         *string          `gorm:"type:varchar(256);comment:官方网站" json:"website"`
	Email           *string          `gorm:"type:varchar(128);comment:联系邮箱" json:"email"`
	Phone           *string          `gorm:"type:varchar(20);comment:联系电话" json:"phone"`
	ContactName     *string          `gorm:"type:varchar(64);comment:联系人姓名" json:"contact_name"`
	ContactPhone    *string          `gorm:"type:varchar(20);comment:联系人电话" json:"contact_phone"`
	ContactEmail    *string          `gorm:"type:varchar(128);comment:联系人邮箱" json:"contact_email"`
	Industry        *string          `gorm:"type:varchar(64);comment:所属行业" json:"industry"`
	Scale           *string          `gorm:"type:varchar(32);comment:组织规模" json:"scale"`
	EstablishedDate *time.Time       `gorm:"type:date;comment:成立日期" json:"established_date"`
	IsActive        bool             `gorm:"type:tinyint(1);default:1;comment:是否启用" json:"is_active"`
	IsVerified      bool             `gorm:"type:tinyint(1);default:0;comment:是否已认证" json:"is_verified"`
	VerifiedAt      *time.Time       `gorm:"type:datetime(3);comment:认证时间" json:"verified_at"`
	ParentID        *string          `gorm:"type:varchar(36);comment:上级组织ID" json:"parent_id"`
	Level           int              `gorm:"type:int;default:1;comment:组织层级" json:"level"`
	SortOrder       int              `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	ExtraData       *JSON            `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt       time.Time        `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time        `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt       *time.Time       `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`
}

func (Organization) TableName() string {
	return "organizations"
}

// EmploymentType 用工类型
type EmploymentType int

const (
	EmploymentTypeRegular   EmploymentType = 1 // 正式
	EmploymentTypeContract  EmploymentType = 2 // 合同
	EmploymentTypeIntern    EmploymentType = 3 // 实习
	EmploymentTypeOutsource EmploymentType = 4 // 外包
)

// WorkStatus 工作状态
type WorkStatus int

const (
	WorkStatusActive   WorkStatus = 1 // 在职
	WorkStatusResigned WorkStatus = 2 // 离职
	WorkStatusRetired  WorkStatus = 3 // 退休
	WorkStatusLeave    WorkStatus = 4 // 停薪留职
)

// OrganizationUser 组织成员关联表
type OrganizationUser struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;comment:记录ID (UUID)" json:"id"`
	OrganizationID  string         `gorm:"type:varchar(36);not null;index;comment:组织ID" json:"organization_id"`
	UserID          string         `gorm:"type:varchar(36);not null;index;comment:用户ID" json:"user_id"`
	DepartmentID    *string        `gorm:"type:varchar(36);comment:部门ID" json:"department_id"`
	DepartmentName  *string        `gorm:"type:varchar(128);comment:部门名称" json:"department_name"`
	Position        *string        `gorm:"type:varchar(64);comment:职位" json:"position"`
	PositionTitle   *string        `gorm:"type:varchar(64);comment:职称" json:"position_title"`
	EmployeeID      *string        `gorm:"type:varchar(64);comment:员工号" json:"employee_id"`
	HireDate        *time.Time     `gorm:"type:date;comment:入职日期" json:"hire_date"`
	ResignationDate *time.Time     `gorm:"type:date;comment:离职日期" json:"resignation_date"`
	EmploymentType  EmploymentType `gorm:"type:tinyint;default:1;comment:用工类型: 1-正式, 2-合同, 3-实习, 4-外包" json:"employment_type"`
	WorkStatus      WorkStatus     `gorm:"type:tinyint;default:1;comment:工作状态: 1-在职, 2-离职, 3-退休, 4-停薪留职" json:"work_status"`
	IsPrimary       bool           `gorm:"type:tinyint(1);default:0;comment:是否为主组织" json:"is_primary"`
	IsAdmin         bool           `gorm:"type:tinyint(1);default:0;comment:是否管理员" json:"is_admin"`
	JoinedAt        time.Time      `gorm:"type:datetime(3);autoCreateTime;comment:加入时间" json:"joined_at"`
	LeftAt          *time.Time     `gorm:"type:datetime(3);comment:离开时间" json:"left_at"`
	CreatedAt       time.Time      `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (OrganizationUser) TableName() string {
	return "organization_users"
}

// DepartmentType 部门类型
type DepartmentType int

const (
	DeptTypeFunctional DepartmentType = 1 // 职能部门
	DeptTypeTeaching   DepartmentType = 2 // 教学部门
	DeptTypeAdmin      DepartmentType = 3 // 行政部门
	DeptTypeOther      DepartmentType = 4 // 其他
)

// Department 部门表
type Department struct {
	ID             string         `gorm:"type:varchar(36);primaryKey;comment:部门ID (UUID)" json:"id"`
	OrganizationID string         `gorm:"type:varchar(36);not null;index;comment:所属组织ID" json:"organization_id"`
	ParentID       *string        `gorm:"type:varchar(36);index;comment:上级部门ID" json:"parent_id"`
	Name           string         `gorm:"type:varchar(128);not null;comment:部门名称" json:"name"`
	Code           *string        `gorm:"type:varchar(64);comment:部门代码" json:"code"`
	Type           DepartmentType `gorm:"type:tinyint;default:1;comment:部门类型: 1-职能部门, 2-教学部门, 3-行政部门, 4-其他" json:"type"`
	Description    *string        `gorm:"type:text;comment:部门描述" json:"description"`
	LeaderID       *string        `gorm:"type:varchar(36);comment:负责人ID" json:"leader_id"`
	SortOrder      int            `gorm:"type:int;default:0;comment:排序序号" json:"sort_order"`
	IsActive       bool           `gorm:"type:tinyint(1);default:1;comment:是否启用" json:"is_active"`
	Level          int            `gorm:"type:int;default:1;comment:部门层级" json:"level"`
	Path           *string        `gorm:"type:varchar(512);comment:部门路径" json:"path"`
	ExtraData      *JSON          `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt      time.Time      `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
	DeletedAt      *time.Time     `gorm:"type:datetime(3);index;comment:删除时间" json:"deleted_at"`
}

func (Department) TableName() string {
	return "departments"
}
