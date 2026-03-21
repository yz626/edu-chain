package config

import (
	"fmt"
	"time"
)

// Config 应用配置
type Config struct {
	Server     ServerConfig     `mapstructure:"server"`
	GRPC       GRPCConfig       `mapstructure:"grpc"`
	Database   DatabaseConfig   `mapstructure:"database"`
	Redis      RedisConfig      `mapstructure:"redis"`
	Logger     LoggerConfig     `mapstructure:"logger"`
	JWT        JWTConfig        `mapstructure:"jwt"`
	Blockchain BlockchainConfig `mapstructure:"blockchain"`
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

// BlockchainConfig 区块链配置
type BlockchainConfig struct {
	Enabled   bool            `mapstructure:"enabled"` // 是否启用区块链
	Type      string          `mapstructure:"type"`    // 区块链类型: fiscobcos, ethereum, fabric
	FISCOBCOS FISCOBCOSConfig `mapstructure:"fiscobcos"`
}

// FISCOBCOSConfig FISCO BCOS区块链配置
type FISCOBCOSConfig struct {
	NodeURL         string        `mapstructure:"node_url"`          // RPC节点URL
	GroupID         string        `mapstructure:"group_id"`          // 群组ID
	ChainID         int           `mapstructure:"chain_id"`          // 链ID
	AccountKey      string        `mapstructure:"account_key"`       // 账户私钥(十六进制)
	AccountKeyFile  string        `mapstructure:"account_key_file"`  // 账户私钥文件路径
	ContractAddress string        `mapstructure:"contract_address"`  // 合约地址
	ContractABI     string        `mapstructure:"contract_abi"`      // 合约ABI(JSON)
	ContractABIPath string        `mapstructure:"contract_abi_path"` // 合约ABI(JSON)文件地址
	Timeout         time.Duration `mapstructure:"timeout"`           // 请求超时时间
	MaxRetries      int           `mapstructure:"max_retries"`       // 最大重试次数
	ConfirmTimeout  time.Duration `mapstructure:"confirm_timeout"`   // 交易确认超时
	GasLimit        uint64        `mapstructure:"gas_limit"`         // Gas限制
	// 国密配置
	GMEnable         bool   `mapstructure:"gm_enable"`           // 是否启用国密
	GMAccountKey     string `mapstructure:"gm_account_key"`      // 国密私钥
	GMAccountCert    string `mapstructure:"gm_account_cert"`     // 国密证书
	GMAccountKeyFile string `mapstructure:"gm_account_key_file"` // 国密私钥文件
	GMSSLCA          string `mapstructure:"gm_ssl_ca"`           // 国密CA证书
	GMSSLCert        string `mapstructure:"gm_ssl_cert"`         // 国密节点证书
	GMSSLKey         string `mapstructure:"gm_ssl_key"`          // 国密节点私钥
}
