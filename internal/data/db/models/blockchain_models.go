package models

import (
	"time"
)

// =====================================================
// 区块链模块 Models
// =====================================================

// BlockchainType 区块链类型
type BlockchainType int

const (
	BlockchainTypeFabric   BlockchainType = 1 // Fabric
	BlockchainTypeEthereum BlockchainType = 2 // Ethereum
)

// NetworkStatus 网络状态
type NetworkStatus int

const (
	NetworkStatusNormal   NetworkStatus = 1 // 正常
	NetworkStatusMaintain NetworkStatus = 2 // 维护中
	NetworkStatusDisabled NetworkStatus = 3 // 停用
)

// BlockchainNetwork 区块链网络表
type BlockchainNetwork struct {
	ID          string         `gorm:"type:varchar(36);primaryKey;comment:网络ID (UUID)" json:"id"`
	Name        string         `gorm:"type:varchar(64);not null;comment:网络名称" json:"name"`
	Code        string         `gorm:"type:varchar(32);uniqueIndex;not null;comment:网络代码 (唯一标识)" json:"code"`
	Type        BlockchainType `gorm:"type:tinyint;default:1;comment:区块链类型: 1-Fabric, 2-Ethereum" json:"type"`
	ChainID     *int           `gorm:"type:int;comment:链ID" json:"chain_id"`
	EndpointURL string         `gorm:"type:varchar(256);not null;comment:节点RPC endpoint URL" json:"endpoint_url"`
	ExplorerURL *string        `gorm:"type:varchar(256);comment:区块链浏览器URL" json:"explorer_url"`
	Status      NetworkStatus  `gorm:"type:tinyint;default:1;comment:状态: 1-正常, 2-维护中, 3-停用" json:"status"`
	IsDefault   bool           `gorm:"type:tinyint(1);default:0;comment:是否默认网络" json:"is_default"`
	ExtraData   *JSON          `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt   time.Time      `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (BlockchainNetwork) TableName() string {
	return "blockchain_networks"
}

// TxType 交易类型
type TxType int

const (
	TxTypeStore    TxType = 1 // 存证
	TxTypeRevoke   TxType = 2 // 撤销
	TxTypeQuery    TxType = 3 // 查询
	TxTypeTransfer TxType = 4 // 转让
)

// TxStatus 交易状态
type TxStatus int

const (
	TxStatusPending  TxStatus = 1 // 待处理
	TxStatusProgress TxStatus = 2 // 处理中
	TxStatusSuccess  TxStatus = 3 // 成功
	TxStatusFailed   TxStatus = 4 // 失败
	TxStatusTimeout  TxStatus = 5 // 超时
)

// BlockchainTransaction 区块链交易表
type BlockchainTransaction struct {
	ID             string     `gorm:"type:varchar(36);primaryKey;comment:交易ID (UUID)" json:"id"`
	TxHash         string     `gorm:"type:varchar(128);uniqueIndex;not null;comment:交易哈希 (唯一)" json:"tx_hash"`
	NetworkID      string     `gorm:"type:varchar(36);not null;index;comment:区块链网络ID" json:"network_id"`
	CertificateID  *string    `gorm:"type:varchar(36);index;comment:关联证书ID" json:"certificate_id"`
	TxType         TxType     `gorm:"type:tinyint;default:1;comment:交易类型: 1-存证, 2-撤销, 3-查询, 4-转让" json:"tx_type"`
	FromAddress    *string    `gorm:"type:varchar(128);comment:发起方地址" json:"from_address"`
	ToAddress      *string    `gorm:"type:varchar(128);comment:接收方地址" json:"to_address"`
	Data           *string    `gorm:"type:text;comment:交易数据" json:"data"`
	Value          *string    `gorm:"type:decimal(38,0);comment:交易金额" json:"value"`
	GasUsed        *int64     `gorm:"type:bigint;comment:Gas消耗" json:"gas_used"`
	TxFee          *string    `gorm:"type:decimal(38,8);comment:交易手续费" json:"tx_fee"`
	BlockNumber    *int64     `gorm:"type:bigint;index;comment:区块高度" json:"block_number"`
	BlockHash      *string    `gorm:"type:varchar(128);comment:区块哈希" json:"block_hash"`
	BlockTimestamp *time.Time `gorm:"type:datetime(3);comment:区块时间戳" json:"block_timestamp"`
	Confirmations  int        `gorm:"type:int;default:0;comment:确认数" json:"confirmations"`
	Status         TxStatus   `gorm:"type:tinyint;default:1;index;comment:交易状态: 1-待处理, 2-处理中, 3-成功, 4-失败, 5-超时" json:"status"`
	ErrorMessage   *string    `gorm:"type:text;comment:错误信息" json:"error_message"`
	CertHash       *string    `gorm:"type:varchar(64);comment:证书数据哈希" json:"cert_hash"`
	OwnerAddress   *string    `gorm:"type:varchar(128);comment:所有者区块链地址" json:"owner_address"`
	TokenID        *string    `gorm:"type:varchar(64);comment:NFT Token ID" json:"token_id"`
	URI            *string    `gorm:"type:text;comment:元数据URI" json:"uri"`
	CreatedAt      time.Time  `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (BlockchainTransaction) TableName() string {
	return "blockchain_transactions"
}

// ContractStatus 合约状态
type ContractStatus int

const (
	ContractStatusNormal     ContractStatus = 1 // 正常
	ContractStatusDisabled   ContractStatus = 2 // 停用
	ContractStatusDeprecated ContractStatus = 3 // 已废弃
)

// SmartContract 智能合约表
type SmartContract struct {
	ID              string         `gorm:"type:varchar(36);primaryKey;comment:合约ID (UUID)" json:"id"`
	Name            string         `gorm:"type:varchar(64);not null;comment:合约名称" json:"name"`
	Code            string         `gorm:"type:varchar(32);uniqueIndex;not null;comment:合约代码 (唯一标识)" json:"code"`
	NetworkID       string         `gorm:"type:varchar(36);not null;index;comment:部署网络ID" json:"network_id"`
	ContractAddress string         `gorm:"type:varchar(128);not null;comment:合约地址" json:"contract_address"`
	ABI             *JSON          `gorm:"type:json;comment:ABI接口定义" json:"abi"`
	Bytecode        *string        `gorm:"type:text;comment:字节码" json:"bytecode"`
	Version         *string        `gorm:"type:varchar(32);comment:合约版本" json:"version"`
	DeployerAddress *string        `gorm:"type:varchar(128);comment:部署者地址" json:"deployer_address"`
	DeployedAt      *time.Time     `gorm:"type:datetime(3);comment:部署时间" json:"deployed_at"`
	Status          ContractStatus `gorm:"type:tinyint;default:1;comment:状态: 1-正常, 2-停用, 3-已废弃" json:"status"`
	IsVerified      bool           `gorm:"type:tinyint(1);default:0;comment:是否已验证源码" json:"is_verified"`
	VerifiedAt      *time.Time     `gorm:"type:datetime(3);comment:验证时间" json:"verified_at"`
	ExtraData       *JSON          `gorm:"type:json;comment:扩展数据" json:"extra_data"`
	CreatedAt       time.Time      `gorm:"type:datetime(3);autoCreateTime;comment:创建时间" json:"created_at"`
	UpdatedAt       time.Time      `gorm:"type:datetime(3);autoUpdateTime;comment:更新时间" json:"updated_at"`
}

func (SmartContract) TableName() string {
	return "smart_contracts"
}
