package main

import (
	"context"
	"fmt"
	"log"

	"github.com/yz626/edu-chain/config"
	fiscobcosabigen "github.com/yz626/edu-chain/internal/blockchain/fiscobcos-abigen"
)

func main() {
	// 加载配置
	blockchainConfig, err := config.LoadBlockchain()
	if err != nil {
		log.Fatalf("Failed to load blockchain config: %v", err)
	}

	// 初始化区块链客户端
	client, err := fiscobcosabigen.NewClient(blockchainConfig)
	if err != nil {
		log.Fatalf("Failed to initialize blockchain client: %v", err)
	}

	// 验证证书
	stats, err := client.GetStats(context.Background())
	if err != nil {
		log.Fatalf("Failed to get stats: %v", err)
	}

	fmt.Println("Stats:", stats)
}
