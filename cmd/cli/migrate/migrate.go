package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/internal/data/db"
	"github.com/yz626/edu-chain/internal/data/db/models"
)

func main() {
	// 解析命令行参数
	migrate := flag.Bool("migrate", false, "Run database migration")
	seed := flag.Bool("seed", false, "Seed database with initial data")
	flag.Parse()

	// 加载配置
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	if err := db.Init(&cfg.Database); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// 执行迁移
	if *migrate {
		runMigration()
	}

	// 执行数据填充
	if *seed {
		runSeed()
	}

	// 如果没有指定任何参数，显示帮助
	if !*migrate && !*seed {
		fmt.Println("Usage:")
		fmt.Println("  migrate -migrate    Run database migration")
		fmt.Println("  migrate -seed       Seed database with initial data")
		fmt.Println("  migrate -migrate -seed  Run migration and seed")
	}
}

func runMigration() {
	fmt.Println("Starting database migration...")

	// 自动迁移所有模型
	err := db.GetDB().AutoMigrate(
		// 用户与权限模块
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.UserRole{},
		&models.RolePermission{},
		&models.RefreshToken{},

		// 组织架构模块
		&models.Organization{},
		&models.OrganizationUser{},
		&models.Department{},

		// 证书管理模块
		&models.CertificateType{},
		&models.CertificateTemplate{},
		&models.Certificate{},
		&models.CertificateBatch{},

		// 区块链模块
		&models.BlockchainNetwork{},
		&models.BlockchainTransaction{},
		&models.SmartContract{},

		// 验证服务模块
		&models.Verification{},

		// 审计日志模块
		&models.AuditLog{},

		// 系统管理模块
		&models.SystemConfig{},
		&models.Dictionary{},
		&models.FileRecord{},
		&models.JobQueue{},
	)

	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Database migration completed successfully!")
}

func runSeed() {
	fmt.Println("Seeding database with initial data...")
	// 这里可以添加初始数据
	// 例如：创建默认管理员角色、权限等
	fmt.Println("Database seeding completed successfully!")
}
