package main

import (
	"log"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/internal/data/db"
	"gorm.io/gen"
)

func main() {
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

	// 生成实例
	g := gen.NewGenerator(gen.Config{
		// 相对执行`go run`时的路径, 会自动创建目录
		// 如果使用go run会以当前目录为起点 如果编辑器运行会以项目工程目录(根目录)为起点
		OutPath:      "internal/data/repository/query", // 查询代码目录
		ModelPkgPath: "/model",                         // 模型代码目录
		// WithDefaultQuery 生成默认查询结构体(作为全局变量使用), 即`Q`结构体和其字段(各表模型)
		// WithoutContext 生成没有context调用限制的代码供查询
		// WithQueryInterface 生成interface形式的查询代码(可导出), 如`Where()`方法返回的就是一个可导出的接口类型
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	// 设置目标 db
	g.UseDB(db.GetDB())

	// 生成所有表（使用默认类型映射）
	// GORM Gen 会将 JSON 类型映射为 string，这是正常行为
	g.GenerateAllTable()

	// 生成所有表的查询代码
	g.ApplyBasic(g.GenerateAllTable()...)

	g.Execute()
}
