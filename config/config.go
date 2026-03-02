package config

import (
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	RabbitMQ   RabbitMQConfig   `mapstructure:"rabbitmq"`
	Blockchain BlockchainConfig `mapstructure:"blockchain"`
	JWT        JWTConfig        `mapstructure:"jwt"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
	Mode string `mapstructure:"mode"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	Name     string `mapstructure:"name"`
	SSLMode  string `mapstructure:"sslmode"`
	MaxOpen  int    `mapstructure:"max_open"`
	MaxIdle  int    `mapstructure:"max_idle"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	VHost    string `mapstructure:"vhost"`
}

// BlockchainConfig 区块链配置
type BlockchainConfig struct {
	// 区块链类型: fiscobcos, ethereum, fabric
	Type      string          `mapstructure:"type"`
	FISCOBCOS FISCOBCOSConfig `mapstructure:"fiscobcos"`
	Ethereum  EthereumConfig  `mapstructure:"ethereum"`
	Fabric    FabricConfig    `mapstructure:"fabric"`
}

// FISCOBCOSConfig FISCO BCOS配置
type FISCOBCOSConfig struct {
	// 节点连接信息
	NodeURL  string   `mapstructure:"node_url"`  // 节点RPC地址，如: http://127.0.0.1:8545
	NodeURLs []string `mapstructure:"node_urls"` // 多个节点地址（群组模式）

	// 群组信息
	GroupID string `mapstructure:"group_id"` // 群组ID，默认: 1

	// 链上账户
	AccountKey     string `mapstructure:"account_key"`      // 账户私钥路径或十六进制
	AccountKeyFile string `mapstructure:"account_key_file"` // 私钥文件路径
	AccountCert    string `mapstructure:"account_cert"`     // 账户证书路径

	// 连接池配置
	MaxConns int `mapstructure:"max_conns"` // 最大连接数
	Timeout  int `mapstructure:"timeout"`   // 超时时间(秒)

	// 合约配置
	ContractAddress string `mapstructure:"contract_address"` // 部署的合约地址
	ContractABI     string `mapstructure:"contract_abi"`     // 合约ABI JSON

	// 链ID
	ChainID int `mapstructure:"chain_id"` // 链ID
}

// EthereumConfig 以太坊配置
type EthereumConfig struct {
	NetworkID   string `mapstructure:"network_id"`
	ChannelName string `mapstructure:"channel_name"`
	Chaincode   string `mapstructure:"chaincode"`
	PeerURL     string `mapstructure:"peer_url"`
	OrdererURL  string `mapstructure:"orderer_url"`
}

// FabricConfig Fabric配置
type FabricConfig struct {
	NetworkID   string `mapstructure:"network_id"`
	ChannelName string `mapstructure:"channel_name"`
	Chaincode   string `mapstructure:"chaincode"`
	PeerURL     string `mapstructure:"peer_url"`
	OrdererURL  string `mapstructure:"orderer_url"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"`
	Expiry int    `mapstructure:"expiry"`
}

// Load 加载配置
func Load() (*Config, error) {
	v := viper.New()
	setDefaults(v)
	v.AutomaticEnv()

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		configPath = "config.yaml"
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

// setDefaults 设置默认值
func setDefaults(v *viper.Viper) {
	v.SetDefault("server.host", "0.0.0.0")
	v.SetDefault("server.port", 8080)
	v.SetDefault("server.mode", "debug")

	v.SetDefault("database.host", "localhost")
	v.SetDefault("database.port", 5432)
	v.SetDefault("database.username", "postgres")
	v.SetDefault("database.password", "postgres")
	v.SetDefault("database.name", "edu_chain")
	v.SetDefault("database.sslmode", "disable")

	v.SetDefault("redis.host", "localhost")
	v.SetDefault("redis.port", 6379)

	v.SetDefault("jwt.secret", "your-secret-key")
	v.SetDefault("jwt.expiry", 86400)
}

// DSN 获取数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Name, c.SSLMode,
	)
}

// Addr 获取Redis地址
func (c *RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// Addr 获取服务器地址
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
