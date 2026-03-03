// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.4.25;

/**
 * @title CertificateRegistry
 * @dev 证书注册表智能合约 - 用于在FISCO BCOS上管理教育证书
 */
contract CertificateRegistry {
    
    // 证书事件
    // 证书颁发事件
    event CertificateIssued(
        string certId, 
        address indexed issuer,
        uint256 timestamp
    );
    
    // 撤销证书事件
    event CertificateRevoked(
        string certId,
        string reason,
        uint256 timestamp
    );
    
    // 存储证书的映射 - 使用单独的映射存储每个字段以避免栈溢出
    // key: 证书ID
    mapping(string => string) private certStudentIds;
    mapping(string => string) private certStudentNames;
    mapping(string => string) private certCourseNames;
    mapping(string => string) private certCourseScores;
    mapping(string => string) private certIssuers;
    mapping(string => uint256) private certIssueDates;
    mapping(string => bool) private certRevoked;
    mapping(string => string) private certRevokeReasons;
    
    // 地址到机构名称的映射
    // key: 区块链地址
    // value: 机构名称
    mapping(address => string) private issuers;
    
    // 证书总数
    uint256 private certificateCount;
    
    // 合约所有者
    address private owner;
    
    // 修饰符：仅合约所有者
    modifier onlyOwner() {
        require(msg.sender == owner, "Only owner can call this function");
        _;
    }
    
    // 修饰符：仅授权颁发机构
    modifier onlyIssuer() {
        require(bytes(issuers[msg.sender]).length > 0, "Not an authorized issuer");
        _;
    }
    
    // 构造函数
    constructor() {
        owner = msg.sender;
    }
    
    /**
     * @dev 添加授权颁发机构
     * @param _issuer 机构地址
     * @param _name 机构名称
     */
    function addIssuer(address _issuer, string _name) public onlyOwner {
        issuers[_issuer] = _name;
    }
    
    /**
     * @dev 自动生成证书ID
     * @return string 生成的证书ID
     */
    function generateCertId() private view returns (string) {
        // 使用 block.number + block.timestamp + 颁发者 + 证书数量 生成唯一ID
        bytes32 hash = keccak256(
            block.number,
            block.timestamp,
            msg.sender,
            certificateCount + 1
        );
        
        // 转换为十六进制字符串
        return bytes32ToString(hash);
    }
    
    /**
     * @dev bytes32 转 string
     */
    function bytes32ToString(bytes32 _bytes32) private pure returns (string) {
        bytes memory bytesString = new bytes(64);
        for (uint8 i = 0; i < 32; i++) {
            bytes1 b = bytes1(_bytes32[i]);
            bytes1 char = bytes1(uint8(b) / 16);
            bytesString[i * 2] = charToHexString(char);
            char = bytes1(uint8(b) - 16 * uint8(char));
            bytesString[i * 2 + 1] = charToHexString(char);
        }
        return string(bytesString);
    }
    
    /**
     * @dev 单字节转十六进制字符
     */
    function charToHexString(bytes1 _char) private pure returns (bytes1) {
        if (_char >= 0x30 && _char <= 0x39) {
            return bytes1(uint8(_char) - 0x30);
        } else if (_char >= 0x61 && _char <= 0x66) {
            return bytes1(uint8(_char) - 0x61 + 10);
        } else if (_char >= 0x41 && _char <= 0x46) {
            return bytes1(uint8(_char) - 0x41 + 10);
        }
        return bytes1(0x30);
    }
    
    /**
     * @dev 检查证书是否存在
     */
    function certExists(string _certId) private view returns (bool) {
        return bytes(certStudentIds[_certId]).length > 0;
    }
    
    /**
     * @dev 自动颁发证书（合约自动生成ID）
     * @param _studentId 学号
     * @param _studentName 学生姓名
     * @param _courseName 课程名称
     * @param _courseScore 课程成绩
     * @param _issueDate 颁发日期
     * @return string 生成的证书ID
     */
    function issueCertificateAuto(
        string _studentId,
        string _studentName,
        string _courseName,
        string _courseScore,
        uint256 _issueDate
    ) public onlyIssuer returns (string) {
        // 自动生成证书ID
        string memory certId = generateCertId();
        
        // 检查ID是否已存在
        require(!certExists(certId), "Certificate ID collision");
        
        // 存储证书数据到各个映射
        certStudentIds[certId] = _studentId;
        certStudentNames[certId] = _studentName;
        certCourseNames[certId] = _courseName;
        certCourseScores[certId] = _courseScore;
        certIssuers[certId] = issuers[msg.sender];
        certIssueDates[certId] = _issueDate;
        certRevoked[certId] = false;
        
        certificateCount++;
        
        CertificateIssued(certId, msg.sender, block.timestamp);
        
        return certId;
    }
    
    /**
     * @dev 颁发证书（手动指定ID）
     * @param _certId 证书ID（手动传入）
     * @param _studentId 学号
     * @param _studentName 学生姓名
     * @param _courseName 课程名称
     * @param _courseScore 课程成绩
     * @param _issueDate 颁发日期
     */
    function issueCertificate(
        string _certId,
        string _studentId,
        string _studentName,
        string _courseName,
        string _courseScore,
        uint256 _issueDate
    ) public onlyIssuer {
        require(!certExists(_certId), "Certificate already exists");
        
        // 存储证书数据到各个映射
        certStudentIds[_certId] = _studentId;
        certStudentNames[_certId] = _studentName;
        certCourseNames[_certId] = _courseName;
        certCourseScores[_certId] = _courseScore;
        certIssuers[_certId] = issuers[msg.sender];
        certIssueDates[_certId] = _issueDate;
        certRevoked[_certId] = false;
        
        certificateCount++;
        
        CertificateIssued(_certId, msg.sender, block.timestamp);
    }
    
    /**
     * @dev 查询证书完整信息
     * @param _certId 证书ID
     */
    function queryCertificate(string _certId) public view returns (
        string,
        string,
        string,
        string,
        string
    ) {
        require(certExists(_certId), "Certificate not found");
        
        return (
            certStudentIds[_certId],
            certStudentNames[_certId],
            certCourseNames[_certId],
            certCourseScores[_certId],
            certIssuers[_certId]
        );
    }
    
    /**
     * @dev 查询证书颁发日期和状态
     */
    function queryCertificateExtra(string _certId) public view returns (
        uint256,
        bool
    ) {
        require(certExists(_certId), "Certificate not found");
        
        return (
            certIssueDates[_certId],
            certRevoked[_certId]
        );
    }
    
    /**
     * @dev 验证证书
     * @param _certId 证书ID
     * @return bool 证书是否存在且未撤销
     */
    function verifyCertificate(string _certId) public view returns (bool) {
        return certExists(_certId) && !getRevoked(_certId);
    }
    
    /**
     * @dev 撤销证书
     * @param _certId 证书ID
     * @param _reason 撤销原因
     */
    function revokeCertificate(string _certId, string _reason) public onlyIssuer {
        require(certExists(_certId), "Certificate not found");
        require(!getRevoked(_certId), "Certificate already revoked");
        
        certRevoked[_certId] = true;
        certRevokeReasons[_certId] = _reason;
        
        certificateCount--;
        
        CertificateRevoked(_certId, _reason, block.timestamp);
    }
    
    /**
     * @dev 获取证书数量
     * @return uint256 已颁发的证书总数
     */
    function getCertificateCount() public view returns (uint256) {
        return certificateCount;
    }
    
    /**
     * @dev 获取学号
     */
    function getStudentId(string _certId) public view returns (string) {
        return certStudentIds[_certId];
    }
    
    /**
     * @dev 获取学生姓名
     */
    function getStudentName(string _certId) public view returns (string) {
        return certStudentNames[_certId];
    }
    
    /**
     * @dev 获取课程名称
     */
    function getCourseName(string _certId) public view returns (string) {
        return certCourseNames[_certId];
    }
    
    /**
     * @dev 获取课程成绩
     */
    function getCourseScore(string _certId) public view returns (string) {
        return certCourseScores[_certId];
    }
    
    /**
     * @dev 获取颁发机构
     */
    function getIssuer(string _certId) public view returns (string) {
        return certIssuers[_certId];
    }
    
    /**
     * @dev 获取颁发日期
     */
    function getIssueDate(string _certId) public view returns (uint256) {
        return certIssueDates[_certId];
    }
    
    /**
     * @dev 获取撤销状态
     */
    function getRevoked(string _certId) public view returns (bool) {
        return certRevoked[_certId];
    }
    
    /**
     * @dev 获取撤销原因
     */
    function getRevokeReason(string _certId) public view returns (string) {
        return certRevokeReasons[_certId];
    }
    
    /**
     * @dev 获取合约所有者
     */
    function getOwner() public view returns (address) {
        return owner;
    }
    
    /**
     * @dev 获取机构名称
     */
    function getIssuerName(address _addr) public view returns (string) {
        return issuers[_addr];
    }
}
