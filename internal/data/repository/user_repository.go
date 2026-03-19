package repository

import (
	"gorm.io/gorm"

	"github.com/yz626/edu-chain/internal/data/db/models"
)

// UserRepository 用户仓储接口
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByPhone(phone string) (*models.User, error)
	Update(user *models.User) error
	Delete(id string) error
	FindByStatus(status models.UserStatus, page, pageSize int) ([]models.User, int64, error)
}

// userRepository 用户仓储实现
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository 创建用户仓储
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

// Create 创建用户
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID 根据ID查询用户
func (r *userRepository) FindByID(id string) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查询用户
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查询用户
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByPhone 根据手机号查询用户
func (r *userRepository) FindByPhone(phone string) (*models.User, error) {
	var user models.User
	err := r.db.Where("phone = ?", phone).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete 删除用户（软删除）
func (r *userRepository) Delete(id string) error {
	return r.db.Delete(&models.User{}, "id = ?", id).Error
}

// FindByStatus 根据状态分页查询用户
func (r *userRepository) FindByStatus(status models.UserStatus, page, pageSize int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	// 统计总数
	countQuery := r.db.Model(&models.User{}).Where("status = ?", status)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Where("status = ?", status).
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&users).Error

	return users, total, err
}

// =====================================================
// 证书仓储示例
// =====================================================

// CertificateRepository 证书仓储接口
type CertificateRepository interface {
	Create(cert *models.Certificate) error
	FindByID(id string) (*models.Certificate, error)
	FindByCertificateNo(certNo string) (*models.Certificate, error)
	FindByUserID(userID string) ([]models.Certificate, error)
	FindByOrganizationID(orgID string, page, pageSize int) ([]models.Certificate, int64, error)
	Update(cert *models.Certificate) error
	Delete(id string) error
}

// certificateRepository 证书仓储实现
type certificateRepository struct {
	db *gorm.DB
}

// NewCertificateRepository 创建证书仓储
func NewCertificateRepository(db *gorm.DB) CertificateRepository {
	return &certificateRepository{db: db}
}

// Create 创建证书
func (r *certificateRepository) Create(cert *models.Certificate) error {
	return r.db.Create(cert).Error
}

// FindByID 根据ID查询证书
func (r *certificateRepository) FindByID(id string) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.db.First(&cert, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// FindByCertificateNo 根据证书编号查询证书
func (r *certificateRepository) FindByCertificateNo(certNo string) (*models.Certificate, error) {
	var cert models.Certificate
	err := r.db.Where("certificate_no = ?", certNo).First(&cert).Error
	if err != nil {
		return nil, err
	}
	return &cert, nil
}

// FindByUserID 根据用户ID查询证书列表
func (r *certificateRepository) FindByUserID(userID string) ([]models.Certificate, error) {
	var certs []models.Certificate
	err := r.db.Where("user_id = ?", userID).Find(&certs).Error
	return certs, err
}

// FindByOrganizationID 根据组织ID分页查询证书
func (r *certificateRepository) FindByOrganizationID(orgID string, page, pageSize int) ([]models.Certificate, int64, error) {
	var certs []models.Certificate
	var total int64

	// 统计总数
	countQuery := r.db.Model(&models.Certificate{}).Where("organization_id = ?", orgID)
	if err := countQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	err := r.db.Where("organization_id = ?", orgID).
		Limit(pageSize).
		Offset(offset).
		Order("created_at DESC").
		Find(&certs).Error

	return certs, total, err
}

// Update 更新证书
func (r *certificateRepository) Update(cert *models.Certificate) error {
	return r.db.Save(cert).Error
}

// Delete 删除证书（软删除）
func (r *certificateRepository) Delete(id string) error {
	return r.db.Delete(&models.Certificate{}, "id = ?", id).Error
}
