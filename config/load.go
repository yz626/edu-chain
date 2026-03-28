package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

var (
	// ConfigPath 主配置文件路径
	ConfigPath = "config/config.yaml"
	// BlockchainConfigPath 区块链配置文件路径
	BlockchainConfigPath = "config/blockchain.yaml"
)

// Load 加载主配置
func Load(path ...string) (*Config, error) {
	if len(path) > 1 {
		return nil, fmt.Errorf("too many config paths")
	}

	v := viper.New()
	setDefaults(v)
	v.AutomaticEnv()

	// 优先读取环境变量中的配置路径，其次用传入参数，最后使用默认路径
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		if len(path) > 0 {
			configPath = path[0]
		} else {
			configPath = ConfigPath
		}
	}

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Config file changed: %s\n", e.Name)
	})

	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &cfg, nil
}

// LoadBlockchain 加载区块链独立配置文件
// 区块链配置独立于主配置，便于在不同环境单独维护证书、节点等信息。
// 路径优先级：环境变量 BLOCKCHAIN_CONFIG_PATH > 传入参数 > 默认路径
func LoadBlockchain(path ...string) (*BlockchainConfig, error) {
	if len(path) > 1 {
		return nil, fmt.Errorf("too many blockchain config paths")
	}

	v := viper.New()
	setBlockchainDefaults(v)

	// 支持通过环境变量覆盖私钥，避免明文写入配置文件
	v.AutomaticEnv()
	v.SetEnvPrefix("BLOCKCHAIN")
	_ = v.BindEnv("account.key", "BLOCKCHAIN_ACCOUNT_KEY")

	configPath := os.Getenv("BLOCKCHAIN_CONFIG_PATH")
	if configPath == "" {
		if len(path) > 0 {
			configPath = path[0]
		} else {
			configPath = BlockchainConfigPath
		}
	}

	v.SetConfigFile(configPath)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read blockchain config file: %w", err)
		}
	}

	// 监听配置文件变更（节点地址、合约地址等运行时可热更新）
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("Blockchain config file changed: %s\n", e.Name)
	})

	var cfg struct {
		Blockchain BlockchainConfig `mapstructure:"blockchain"`
	}
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal blockchain config: %w", err)
	}

	return &cfg.Blockchain, nil
}

// setDefaults 主配置默认值
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")
	v.SetDefault("server.read_timeout", 30)
	v.SetDefault("server.write_timeout", 30)

	v.SetDefault("grpc.host", "0.0.0.0")
	v.SetDefault("grpc.port", 9090)

	v.SetDefault("database.name", "mysql")
	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 3306)
	v.SetDefault("database.username", "root")
	v.SetDefault("database.password", "123456")
	v.SetDefault("database.database", "edu_chain")
	v.SetDefault("database.sslmode", "disable")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)

	v.SetDefault("logger.level", "info")
	v.SetDefault("logger.format", "json")
	v.SetDefault("logger.directory", "logs")
	v.SetDefault("logger.console", true)
	v.SetDefault("logger.max_size", 100)
	v.SetDefault("logger.max_age", 30)
	v.SetDefault("logger.max_backups", 10)
	v.SetDefault("logger.compress", true)
	v.SetDefault("logger.enable_stacktrace", false)

	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.expire", 3600)
	v.SetDefault("jwt.refresh_expire", 604800)
	v.SetDefault("jwt.issuer", "edu-chain")
}

// setBlockchainDefaults 区块链配置默认值
func setBlockchainDefaults(v *viper.Viper) {
	v.SetDefault("blockchain.enabled", false) // 默认关闭，需显式开启
	v.SetDefault("blockchain.group_id", "group0")
	v.SetDefault("blockchain.chain_id", "chain0")
	v.SetDefault("blockchain.nodes", []string{"127.0.0.1:20200"})
	v.SetDefault("blockchain.tls.enabled", true)
	v.SetDefault("blockchain.tls.ca_cert", "config/certs/ca.crt")
	v.SetDefault("blockchain.tls.client_cert", "config/certs/sdk.crt")
	v.SetDefault("blockchain.tls.client_key", "config/certs/sdk.key")
	v.SetDefault("blockchain.account.key_file", "config/certs/account.key")
	v.SetDefault("blockchain.timeout", 30)
	v.SetDefault("blockchain.retry_times", 3)
	v.SetDefault("blockchain.retry_interval", 2)
	v.SetDefault("blockchain.contract.address", "0x0000000000000000000000000000000000000000")
	v.SetDefault("blockchain.contract.abi_file", "contracts/abi/CertificateRegistry.abi")
	v.SetDefault("blockchain.contract.owner_name", "教育部学历认证中心")
}
