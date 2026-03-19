package main

import (
	"log"

	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/internal/data/db"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化数据库连接
	if err := db.Init(conf); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
}
