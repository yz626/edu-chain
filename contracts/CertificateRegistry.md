# CertificateRegistry 智能合约文档

## 1. 合约概述

### 基本信息

| 属性 | 值 |
|------|-----|
| 合约名称 | CertificateRegistry |
| 用途 | 证书注册表 - 用于在FISCO BCOS上管理教育证书 |
| 编译器版本 | ^0.4.25 |
| 许可证 | Apache-2.0 |

### 合约功能

CertificateRegistry 是一个教育证书管理智能合约，运行在FISCO BCOS区块链平台上。该合约提供了完整的证书生命周期管理功能，包括证书颁发、查询、验证和撤销。

---

## 2. 数据结构

### 2.1 证书存储

合约使用多个独立的 mapping 分别存储证书的各个字段，以避免Solidity中的栈溢出问题：

```solidity
mapping(string => string) private certStudentIds;      // 学号
mapping(string => string) private certStudentNames;    // 学生姓名
mapping(string => string) private certCourseNames;    // 课程名称
mapping(string => string) private certCourseScores;   // 课程成绩
mapping(string => string) private certIssuers;        // 颁发机构
mapping(string => uint256) private certIssueDates;     // 颁发日期
mapping(string => bool) private certRevoked;          // 撤销状态
mapping(string => string) private certRevokeReasons;  // 撤销原因
```

### 2.2 机构管理

```solidity
mapping(address => string) private issuers;  // 地址到机构名称的映射
```

### 2.3 其他状态变量

```solidity
uint256 private certificateCount;  // 证书总数
address private owner;             // 合约所有者
```

---

## 3. 访问控制

### 3.1 修饰符

| 修饰符 | 说明 |
|--------|------|
| `onlyOwner` | 仅合约所有者可调用 |
| `onlyIssuer` | 仅授权颁发机构可调用 |

### 3.2 权限说明

- **合约所有者**：拥有添加授权颁发机构的权限
- **授权颁发机构**：可以颁发、查询和撤销证书

---

## 4. 事件

### 4.1 CertificateIssued

证书颁发事件，在证书成功颁发时触发：

```solidity
event CertificateIssued(
    string certId,           // 证书ID
    address indexed issuer,  // 颁发机构地址
    uint256 timestamp        // 颁发时间
);
```

### 4.2 CertificateRevoked

证书撤销事件，在证书被撤销时触发：

```solidity
event CertificateRevoked(
    string certId,     // 撤销的证书ID
    string reason,     // 撤销原因
    uint256 timestamp  // 撤销时间
);
```

---

## 5. 函数说明

### 5.1 管理函数

#### addIssuer

添加授权颁发机构。

```solidity
function addIssuer(address _issuer, string _name) public onlyOwner
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_issuer` | address | 机构地址 |
| `_name` | string | 机构名称 |

---

### 5.2 证书颁发函数

#### issueCertificateAuto

自动颁发证书（系统自动生成证书ID）。

```solidity
function issueCertificateAuto(
    string _studentId,
    string _studentName,
    string _courseName,
    string _courseScore,
    uint256 _issueDate
) public onlyIssuer returns (string)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_studentId` | string | 学号 |
| `_studentName` | string | 学生姓名 |
| `_courseName` | string | 课程名称 |
| `_courseScore` | string | 课程成绩 |
| `_issueDate` | uint256 | 颁发日期 |

**返回值**：生成的证书ID

---

#### issueCertificate

手动指定证书ID颁发证书。

```solidity
function issueCertificate(
    string _certId,
    string _studentId,
    string _studentName,
    string _courseName,
    string _courseScore,
    uint256 _issueDate
) public onlyIssuer
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_certId` | string | 证书ID（手动指定） |
| `_studentId` | string | 学号 |
| `_studentName` | string | 学生姓名 |
| `_courseName` | string | 课程名称 |
| `_courseScore` | string | 课程成绩 |
| `_issueDate` | uint256 | 颁发日期 |

---

### 5.3 查询函数

#### queryCertificate

查询证书完整信息。

```solidity
function queryCertificate(string _certId) public view returns (
    string,  // 学号
    string,  // 学生姓名
    string,  // 课程名称
    string,  // 课程成绩
    string   // 颁发机构
)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_certId` | string | 证书ID |

---

#### queryCertificateExtra

查询证书颁发日期和撤销状态。

```solidity
function queryCertificateExtra(string _certId) public view returns (
    uint256,  // 颁发日期
    bool      // 撤销状态
)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_certId` | string | 证书ID |

---

#### verifyCertificate

验证证书有效性。

```solidity
function verifyCertificate(string _certId) public view returns (bool)
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_certId` | string | 证书ID |

**返回值**：证书是否存在且未撤销

---

### 5.4 撤销函数

#### revokeCertificate

撤销证书。

```solidity
function revokeCertificate(string _certId, string _reason) public onlyIssuer
```

| 参数 | 类型 | 说明 |
|------|------|------|
| `_certId` | string | 证书ID |
| `_reason` | string | 撤销原因 |

---

### 5.5 辅助查询函数

| 函数 | 返回值 | 说明 |
|------|--------|------|
| `getCertificateCount()` | uint256 | 获取证书总数 |
| `getStudentId(string)` | string | 获取学号 |
| `getStudentName(string)` | string | 获取学生姓名 |
| `getCourseName(string)` | string | 获取课程名称 |
| `getCourseScore(string)` | string | 获取课程成绩 |
| `getIssuer(string)` | string | 获取颁发机构 |
| `getIssueDate(string)` | uint256 | 获取颁发日期 |
| `getRevoked(string)` | bool | 获取撤销状态 |
| `getRevokeReason(string)` | string | 获取撤销原因 |
| `getOwner()` | address | 获取合约所有者 |
| `getIssuerName(address)` | string | 获取机构名称 |

---

## 6. 内部函数

### 6.1 generateCertId

自动生成唯一证书ID。

```solidity
function generateCertId() private view returns (string)
```

使用 `block.number + block.timestamp + 颁发者 + 证书数量` 通过 keccak256 哈希生成唯一ID。

### 6.2 bytes32ToString

将 bytes32 转换为十六进制字符串。

### 6.3 charToHexString

将单字节转换为十六进制字符。

### 6.4 certExists

检查证书是否存在。

---

## 7. 使用流程

### 7.1 初始化流程

1. 部署合约
2. 合约所有者调用 `addIssuer` 添加授权机构

### 7.2 证书颁发流程

1. 授权机构调用 `issueCertificate` 或 `issueCertificateAuto`
2. 合约触发 `CertificateIssued` 事件
3. 证书总数增加

### 7.3 证书验证流程

1. 调用方调用 `verifyCertificate` 验证证书
2. 返回证书是否存在且未撤销

### 7.4 证书撤销流程

1. 授权机构调用 `revokeCertificate`
2. 合约触发 `CertificateRevoked` 事件
3. 证书总数减少

---

## 8. 注意事项

1. **存储优化**：合约使用多个独立的 mapping 存储不同字段，避免栈溢出
2. **ID生成**：自动生成的证书ID基于区块信息和颁发者地址，确保唯一性
3. **撤销状态**：证书撤销后，证书总数会相应减少
4. **权限控制**：只有授权机构才能操作证书，合约所有者管理授权机构
