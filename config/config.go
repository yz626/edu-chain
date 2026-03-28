package config

import (
	"fmt"
)

// Config 应用配置
type Config struct {
	Server   ServerConfig   `mapstructure:"server"`
	GRPC     GRPCConfig     `mapstructure:"grpc"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Logger   LoggerConfig   `mapstructure:"logger"`
	JWT      JWTConfig      `mapstructure:"jwt"`
}

// ================================================================
// 区块链配置（从独立文件 config/blockchain.yaml 加载）
// ================================================================

// BlockchainConfig 区块链总配置
type BlockchainConfig struct {
	Enabled bool              `mapstructure:"enabled"` // 是否启用区块链
	Nodes   []string          `mapstructure:"nodes"`   // 节点地址列表
	GroupID string            `mapstructure:"group_id"` // 群组 ID
	ChainID string            `mapstructure:"chain_id"` // 链 ID
	TLS     BlockchainTLS     `mapstructure:"tls"`
	Account BlockchainAccount `mapstructure:"account"`
	Contract BlockchainContract `mapstructure:"contract"`

	// 连接参数
	Timeout       int `mapstructure:"timeout"`        // 请求超时（秒）
	RetryTimes    int `mapstructure:"retry_times"`    // 失败重试次数
	RetryInterval int `mapstructure:"retry_interval"` // 重试间隔（秒）
}

// BlockchainTLS TLS 证书配置
type BlockchainTLS struct {
	Enabled    bool   `mapstructure:"enabled"`     // 是否启用 TLS
	CACert     string `mapstructure:"ca_cert"`     // CA 根证书路径
	ClientCert string `mapstructure:"client_cert"` // 客户端证书路径
	ClientKey  string `mapstructure:"client_key"`  // 客户端私钥路径
}

// BlockchainAccount 账户（签名交易用）配置
type BlockchainAccount struct {
	Key     string `mapstructure:"key"`      // 私钥十六进制（优先使用）
	KeyFile string `mapstructure:"key_file"` // 私钥文件路径（次选）
	Address string `mapstructure:"address"`  // 账户地址（可留空，由 key 推导）
}

// BlockchainContract 智能合约配置
type BlockchainContract struct {
	Address   string `mapstructure:"address"`    // 合约部署地址
	ABIFile   string `mapstructure:"abi_file"`   // ABI 文件路径
	OwnerName string `mapstructure:"owner_name"` // 部署者名称（首次部署用）
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string `mapstructure:"host"`          // 服务器地址
	Port         int    `mapstructure:"port"`          // 服务器端口
	Mode         string `mapstructure:"mode"`          // 运行模式 (debug, release)
	ReadTimeout  int    `mapstructure:"read_timeout"`  // 读取超时(秒)
	WriteTimeout int    `mapstructure:"write_timeout"` // 写入超时(秒)
}

// GRPCConfig gRPC服务器配置
type GRPCConfig struct {
	Host string `mapstructure:"host"` // gRPC服务器地址
	Port int    `mapstructure:"port"` // gRPC服务器端口
}

// Addr 获取服务器地址
func (c *ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Name         string `mapstructure:"name"`           // 数据库名称
	Host         string `mapstructure:"host"`           // 数据库地址
	Port         int    `mapstructure:"port"`           // 数据库端口
	Username     string `mapstructure:"username"`       // 数据库用户名
	Password     string `mapstructure:"password"`       // 数据库密码
	Database     string `mapstructure:"database"`       // 数据库名称
	SSLMode      string `mapstructure:"sslmode"`        // SSL模式
	MaxOpenConns int    `mapstructure:"max_open_conns"` // 最大连接数
	MaxIdleConns int    `mapstructure:"max_idle_conns"` // 最大空闲连接数
	MaxLifetime  int    `mapstructure:"max_lifetime"`   // 连接超时时间
	Timeout      int    `mapstructure:"timeout"`        // 连接超时时间
}

// DSN 获取数据库连接字符串
func (c *DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.Host, c.Port, c.Username, c.Password, c.Database, c.SSLMode,
	)
}

// MySQLDSN 获取MySQL连接字符串
func (c *DatabaseConfig) MySQLDSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Username, c.Password, c.Host, c.Port, c.Database,
	)
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Password string `mapstructure:"password"`
	DB       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level            string `mapstructure:"level"`             // 日志级别
	Format           string `mapstructure:"format"`            // 日志格式
	Directory        string `mapstructure:"directory"`         // 日志目录
	Console          bool   `mapstructure:"console"`           // 是否输出到控制台
	MaxSize          int    `mapstructure:"max_size"`          // 单个日志文件最大大小(MB)
	MaxAge           int    `mapstructure:"max_age"`           // 日志文件保留天数
	MaxBackups       int    `mapstructure:"max_backups"`       // 保留的日志文件数量
	Compress         bool   `mapstructure:"compress"`          // 是否压缩旧日志
	EnableStacktrace bool   `mapstructure:"enable_stacktrace"` // 是否启用堆栈跟踪
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret        string `mapstructure:"secret"`         // JWT密钥
	Expire        int    `mapstructure:"expire"`         // 访问令牌过期时间（秒）
	RefreshExpire int    `mapstructure:"refresh_expire"` // 刷新令牌过期时间（秒）
	Issuer        string `mapstructure:"issuer"`         // JWT签发者
}
