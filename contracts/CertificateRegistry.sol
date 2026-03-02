// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.0;

/**
 * @title CertificateRegistry
 * @dev 证书注册表智能合约 - 用于在FISCO BCOS上管理教育证书
 */
contract CertificateRegistry {
    
    // 证书结构体
    struct Certificate {
        string certId;           // 证书唯一ID
        string studentId;        // 学号
        string studentName;      // 学生姓名
        string courseName;       // 课程名称
        string courseScore;      // 课程成绩
        string issuer;           // 颁发机构
        uint256 issueDate;       // 颁发日期(时间戳)
        bool revoked;            // 是否已撤销
        string revokeReason;     // 撤销原因
    }
    
    // 证书事件
    // 证书颁发事件
    event CertificateIssued(
        string indexed certId, 
        address indexed issuer,
        uint256 timestamp
    );
    
    // 撤销证书事件
    event CertificateRevoked(
        string indexed certId,
        string reason,
        uint256 timestamp
    );
    
    // 存储证书的映射
    // key: 证书ID
    // value: 证书结构体
    mapping(string => Certificate) public certificates;
    
    // 地址到机构名称的映射
    // key: 地址
    // value: 机构名称
    mapping(address => string) public issuers;
    
    // 合约所有者
    address public owner;
    
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
    function addIssuer(address _issuer, string memory _name) public onlyOwner {
        issuers[_issuer] = _name;
    }
    
    /**
     * @dev 颁发证书
     * @param _certId 证书ID
     * @param _studentId 学号
     * @param _studentName 学生姓名
     * @param _courseName 课程名称
     * @param _courseScore 课程成绩
     * @param _issueDate 颁发日期
     */
    function issueCertificate(
        string memory _certId,
        string memory _studentId,
        string memory _studentName,
        string memory _courseName,
        string memory _courseScore,
        uint256 _issueDate
    ) public onlyIssuer {
        require(bytes(certificates[_certId].certId).length == 0, "Certificate already exists");
        
        certificates[_certId] = Certificate({
            certId: _certId,
            studentId: _studentId,
            studentName: _studentName,
            courseName: _courseName,
            courseScore: _courseScore,
            issuer: issuers[msg.sender],
            issueDate: _issueDate,
            revoked: false,
            revokeReason: ""
        });
        
        emit CertificateIssued(_certId, msg.sender, block.timestamp);
    }
    
    /**
     * @dev 查询证书
     * @param _certId 证书ID
     */
    function queryCertificate(string memory _certId) public view returns (
        string memory,
        string memory,
        string memory,
        string memory,
        string memory,
        string memory,
        uint256,
        bool
    ) {
        Certificate memory cert = certificates[_certId];
        require(bytes(cert.certId).length > 0, "Certificate not found");
        
        return (
            cert.certId,
            cert.studentId,
            cert.studentName,
            cert.courseName,
            cert.courseScore,
            cert.issuer,
            cert.issueDate,
            cert.revoked
        );
    }
    
    /**
     * @dev 验证证书
     * @param _certId 证书ID
     * @return bool 证书是否存在且未撤销
     */
    function verifyCertificate(string memory _certId) public view returns (bool) {
        Certificate memory cert = certificates[_certId];
        return bytes(cert.certId).length > 0 && !cert.revoked;
    }
    
    /**
     * @dev 撤销证书
     * @param _certId 证书ID
     * @param _reason 撤销原因
     */
    function revokeCertificate(string memory _certId, string memory _reason) public onlyIssuer {
        require(bytes(certificates[_certId].certId).length > 0, "Certificate not found");
        require(!certificates[_certId].revoked, "Certificate already revoked");
        
        certificates[_certId].revoked = true;
        certificates[_certId].revokeReason = _reason;
        
        emit CertificateRevoked(_certId, _reason, block.timestamp);
    }
    
    /**
     * @dev 获取证书数量（演示用）
     */
    function getCertificateCount() public view returns (uint256) {
        return 0; // 实际实现需要额外存储
    }
}
