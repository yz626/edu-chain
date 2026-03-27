// SPDX-License-Identifier: Apache-2.0
pragma solidity ^0.8.11;

/**
 * @title  CertificateRegistry
 * @notice EduChain 教育证书链上存证合约
 *
 * @dev 设计原则
 *   1. 合约只负责链上核心数据存证（证书哈希 + 状态），
 *      详细业务数据（姓名、专业等）存储在链下 MySQL，
 *      通过 certId（链下 UUID 转 bytes32）与链上记录关联。
 *   2. 结构清晰：Owner 管理 → 机构授权 → 证书颁发/撤销/恢复 → 查询验证。
 *   3. 避免过度复杂：不引入 NFT/Token，不做链上分页，不在链上存储大量字符串。
 *
 * 角色说明
 *   owner  : 合约部署者（监管方/教育部门），负责授权颁发机构
 *   issuer : 授权的颁发机构（高校），负责颁发/撤销/恢复证书
 */
contract CertificateRegistry {

    // =========================================================
    // 自定义错误（比 require(false, "msg") 节省 gas）
    // =========================================================

    error Unauthorized();                       // 调用者无权限
    error ZeroAddress();                        // 传入了零地址
    error EmptyParam();                         // 必填参数为空
    error BatchEmpty();                         // 批量数组长度为 0
    error ArrayLengthMismatch();                // 两个并行数组长度不同
    error ContractPaused();                     // 合约处于暂停状态
    error ContractNotPaused();                  // 合约未暂停（unpause 前置检查）
    error IssuerNotAuthorized(address issuer);  // 机构未被授权
    error CertAlreadyExists(bytes32 certId);    // 证书 ID 重复
    error CertNotFound(bytes32 certId);         // 证书不存在
    error CertAlreadyRevoked(bytes32 certId);   // 证书已撤销
    error CertNotRevoked(bytes32 certId);       // 证书未撤销（无法恢复）
    error NotCertIssuer(bytes32 certId);        // 非本证书的颁发机构

    // =========================================================
    // 数据结构
    // =========================================================

    /**
     * @dev 链上证书记录（只存核心字段）
     *
     * certHash     : 证书完整内容的 SHA-256 哈希，用于链下数据完整性校验
     * issuer       : 颁发机构区块链地址
     * issuedAt     : 颁发时间戳（Unix 秒，uint64 足够用到 year 2554）
     * revoked      : 是否已撤销
     * revokedAt    : 撤销时间戳（0 = 未撤销）
     * revokeReason : 撤销原因简述（详情存链下）
     */
    struct CertRecord {
        bytes32 certHash;
        address issuer;
        uint64  issuedAt;
        bool    revoked;
        uint64  revokedAt;
        string  revokeReason;
    }

    /**
     * @dev 颁发机构信息
     *
     * name         : 机构名称（便于链上核查）
     * authorized   : 当前是否处于授权状态
     * authorizedAt : 最近一次授权时间戳
     */
    struct IssuerInfo {
        string name;
        bool   authorized;
        uint64 authorizedAt;
    }

    // =========================================================
    // 状态变量
    // =========================================================

    /// @dev 合约所有者（监管方）
    address private _owner;

    /// @dev 紧急暂停标志
    bool private _paused;

    /// @dev certId => 证书记录
    mapping(bytes32 => CertRecord) private _certs;

    /// @dev certId => 是否已存在（快速检查，避免读整个 struct）
    mapping(bytes32 => bool) private _certExists;

    /// @dev 机构地址 => 机构信息
    mapping(address => IssuerInfo) private _issuers;

    /// @dev 已颁发证书总数
    uint256 private _totalIssued;

    /// @dev 已撤销证书净数（恢复后减回）
    uint256 private _totalRevoked;

    // =========================================================
    // 事件
    // =========================================================

    event OwnershipTransferred(address indexed previousOwner, address indexed newOwner);
    event Paused(address indexed by);
    event Unpaused(address indexed by);

    /// @notice 机构被授权（或重新授权）
    event IssuerAuthorized(address indexed issuer, string name, uint64 timestamp);

    /// @notice 机构授权被撤销
    event IssuerRevoked(address indexed issuer, uint64 timestamp);

    /// @notice 单张证书上链
    event CertificateIssued(
        bytes32 indexed certId,
        bytes32 indexed certHash,
        address indexed issuer,
        uint64          timestamp
    );

    /// @notice 证书被撤销
    event CertificateRevoked(
        bytes32 indexed certId,
        string          reason,
        address indexed by,
        uint64          timestamp
    );

    /// @notice 证书被恢复（撤销操作被撤回）
    event CertificateRestored(
        bytes32 indexed certId,
        address indexed by,
        uint64          timestamp
    );

    /// @notice 批量颁发完成
    event BatchIssued(address indexed issuer, uint256 count, uint64 timestamp);

    // =========================================================
    // 修饰符
    // =========================================================

    modifier onlyOwner() {
        if (msg.sender != _owner) revert Unauthorized();
        _;
    }

    modifier onlyIssuer() {
        if (!_issuers[msg.sender].authorized) revert IssuerNotAuthorized(msg.sender);
        _;
    }

    modifier whenNotPaused() {
        if (_paused) revert ContractPaused();
        _;
    }

    modifier whenPaused() {
        if (!_paused) revert ContractNotPaused();
        _;
    }

    modifier certMustExist(bytes32 certId) {
        if (!_certExists[certId]) revert CertNotFound(certId);
        _;
    }

    // =========================================================
    // 构造函数
    // =========================================================

    /**
     * @notice 部署合约
     * @param ownerName 部署机构名称（如 "Ministry of Education"）
     *                  部署者自动成为 owner 并获得颁发机构资格
     */
    constructor(string memory ownerName) {
        if (bytes(ownerName).length == 0) revert EmptyParam();
        _owner  = msg.sender;
        _paused = false;
        uint64 ts = uint64(block.timestamp);
        _issuers[msg.sender] = IssuerInfo({
            name:         ownerName,
            authorized:   true,
            authorizedAt: ts
        });
        emit IssuerAuthorized(msg.sender, ownerName, ts);
    }

    // =========================================================
    // Owner 管理
    // =========================================================

    /// @notice 返回当前合约所有者地址
    function owner() external view returns (address) {
        return _owner;
    }

    /// @notice 转移合约所有权
    function transferOwnership(address newOwner) external onlyOwner {
        if (newOwner == address(0)) revert ZeroAddress();
        address prev = _owner;
        _owner = newOwner;
        emit OwnershipTransferred(prev, newOwner);
    }

    // =========================================================
    // 紧急暂停
    // =========================================================

    /// @notice 暂停合约（仅 owner，紧急情况使用）
    function pause() external onlyOwner whenNotPaused {
        _paused = true;
        emit Paused(msg.sender);
    }

    /// @notice 恢复合约运行
    function unpause() external onlyOwner whenPaused {
        _paused = false;
        emit Unpaused(msg.sender);
    }

    /// @notice 查询合约是否已暂停
    function paused() external view returns (bool) {
        return _paused;
    }

    // =========================================================
    // 颁发机构管理（仅 owner）
    // =========================================================

    /// @notice 添加或重新授权一个颁发机构
    function addIssuer(address issuer, string calldata name) external onlyOwner {
        if (issuer == address(0)) revert ZeroAddress();
        if (bytes(name).length == 0) revert EmptyParam();
        uint64 ts = uint64(block.timestamp);
        _issuers[issuer] = IssuerInfo({ name: name, authorized: true, authorizedAt: ts });
        emit IssuerAuthorized(issuer, name, ts);
    }

    /// @notice 批量添加授权颁发机构
    function addIssuerBatch(address[] calldata issuers, string[] calldata names) external onlyOwner {
        uint256 len = issuers.length;
        if (len == 0) revert BatchEmpty();
        if (len != names.length) revert ArrayLengthMismatch();
        uint64 ts = uint64(block.timestamp);
        for (uint256 i = 0; i < len; ++i) {
            if (issuers[i] == address(0)) revert ZeroAddress();
            if (bytes(names[i]).length == 0) revert EmptyParam();
            _issuers[issuers[i]] = IssuerInfo({ name: names[i], authorized: true, authorizedAt: ts });
            emit IssuerAuthorized(issuers[i], names[i], ts);
        }
    }

    /// @notice 撤销机构授权（历史证书记录不受影响，链上证据永久保留）
    function revokeIssuer(address issuer) external onlyOwner {
        if (issuer == address(0)) revert ZeroAddress();
        if (!_issuers[issuer].authorized) revert IssuerNotAuthorized(issuer);
        _issuers[issuer].authorized = false;
        emit IssuerRevoked(issuer, uint64(block.timestamp));
    }

    /// @notice 查询机构授权信息
    function getIssuerInfo(address issuer)
        external view
        returns (bool authorized, string memory name, uint64 authorizedAt)
    {
        IssuerInfo storage info = _issuers[issuer];
        return (info.authorized, info.name, info.authorizedAt);
    }

    // =========================================================
    // 证书颁发（仅授权机构，合约未暂停时）
    // =========================================================

    /**
     * @notice 单张证书上链存证
     * @param certId   链下证书唯一标识（bytes32）
     *                 推荐：keccak256(abi.encodePacked(uuidString))
     * @param certHash 证书完整数据的 SHA-256 哈希（bytes32）
     */
    function issueCertificate(bytes32 certId, bytes32 certHash) external onlyIssuer whenNotPaused {
        if (certId   == bytes32(0)) revert EmptyParam();
        if (certHash == bytes32(0)) revert EmptyParam();
        if (_certExists[certId]) revert CertAlreadyExists(certId);
        uint64 ts = uint64(block.timestamp);
        _certs[certId] = CertRecord({
            certHash:     certHash,
            issuer:       msg.sender,
            issuedAt:     ts,
            revoked:      false,
            revokedAt:    0,
            revokeReason: ""
        });
        _certExists[certId] = true;
        unchecked { ++_totalIssued; }
        emit CertificateIssued(certId, certHash, msg.sender, ts);
    }

    /**
     * @notice 批量证书上链存证
     * @param certIds    证书 ID 数组
     * @param certHashes 对应的证书内容哈希数组（须与 certIds 等长）
     */
    function issueCertificateBatch(
        bytes32[] calldata certIds,
        bytes32[] calldata certHashes
    ) external onlyIssuer whenNotPaused {
        uint256 len = certIds.length;
        if (len == 0) revert BatchEmpty();
        if (len != certHashes.length) revert ArrayLengthMismatch();
        uint64 ts = uint64(block.timestamp);
        for (uint256 i = 0; i < len; ++i) {
            bytes32 cid  = certIds[i];
            bytes32 hash = certHashes[i];
            if (cid  == bytes32(0)) revert EmptyParam();
            if (hash == bytes32(0)) revert EmptyParam();
            if (_certExists[cid]) revert CertAlreadyExists(cid);
            _certs[cid] = CertRecord({
                certHash:     hash,
                issuer:       msg.sender,
                issuedAt:     ts,
                revoked:      false,
                revokedAt:    0,
                revokeReason: ""
            });
            _certExists[cid] = true;
            emit CertificateIssued(cid, hash, msg.sender, ts);
        }
        unchecked { _totalIssued += len; }
        emit BatchIssued(msg.sender, len, ts);
    }

    // =========================================================
    // 证书撤销与恢复（仅原颁发机构，合约未暂停时）
    // =========================================================

    /**
     * @notice 撤销证书
     * @dev  只有该证书的原颁发机构可以撤销，防止跨机构误操作
     * @param certId 证书 ID
     * @param reason 撤销原因简述（详情存链下）
     */
    function revokeCertificate(bytes32 certId, string calldata reason)
        external onlyIssuer whenNotPaused certMustExist(certId)
    {
        CertRecord storage rec = _certs[certId];
        if (rec.revoked)              revert CertAlreadyRevoked(certId);
        if (rec.issuer != msg.sender) revert NotCertIssuer(certId);
        uint64 ts = uint64(block.timestamp);
        rec.revoked      = true;
        rec.revokedAt    = ts;
        rec.revokeReason = reason;
        unchecked { ++_totalRevoked; }
        emit CertificateRevoked(certId, reason, msg.sender, ts);
    }

    /**
     * @notice 恢复被撤销的证书（撤销操作的撤回）
     * @dev  只有该证书的原颁发机构可以恢复
     * @param certId 证书 ID
     */
    function restoreCertificate(bytes32 certId)
        external onlyIssuer whenNotPaused certMustExist(certId)
    {
        CertRecord storage rec = _certs[certId];
        if (!rec.revoked)             revert CertNotRevoked(certId);
        if (rec.issuer != msg.sender) revert NotCertIssuer(certId);
        uint64 ts = uint64(block.timestamp);
        rec.revoked      = false;
        rec.revokedAt    = 0;
        rec.revokeReason = "";
        unchecked { --_totalRevoked; }
        emit CertificateRestored(certId, msg.sender, ts);
    }

    // =========================================================
    // 查询与验证（纯读函数，无 gas 消耗）
    // =========================================================

    /// @notice 查询证书是否存在
    function certExists(bytes32 certId) external view returns (bool) {
        return _certExists[certId];
    }

    /**
     * @notice 获取证书完整链上记录
     * @param  certId       证书 ID
     * @return certHash     证书内容哈希
     * @return issuer       颁发机构地址
     * @return issuedAt     颁发时间戳
     * @return revoked      是否已撤销
     * @return revokedAt    撤销时间戳（0 = 未撤销）
     * @return revokeReason 撤销原因
     */
    function getCertificate(bytes32 certId)
        external view certMustExist(certId)
        returns (
            bytes32       certHash,
            address       issuer,
            uint64        issuedAt,
            bool          revoked,
            uint64        revokedAt,
            string memory revokeReason
        )
    {
        CertRecord storage rec = _certs[certId];
        return (rec.certHash, rec.issuer, rec.issuedAt, rec.revoked, rec.revokedAt, rec.revokeReason);
    }

    /**
     * @notice 验证证书：确认链上哈希与传入哈希一致，且证书未被撤销
     * @param  certId   证书 ID
     * @param  certHash 待验证的证书内容哈希
     * @return valid    true = 哈希匹配且证书有效
     * @return revoked  true = 证书已被撤销
     */
    function verifyCertificate(bytes32 certId, bytes32 certHash)
        external view
        returns (bool valid, bool revoked)
    {
        if (!_certExists[certId]) return (false, false);
        CertRecord storage rec = _certs[certId];
        revoked = rec.revoked;
        valid   = (!revoked) && (rec.certHash == certHash);
    }

    /**
     * @notice 批量验证证书（一次调用验证多张，节省 RPC 调用次数）
     * @param  certIds    证书 ID 数组
     * @param  certHashes 对应的证书内容哈希数组（须与 certIds 等长）
     * @return valids     每张证书是否有效
     * @return revokeds   每张证书是否已撤销
     */
    function verifyCertificateBatch(
        bytes32[] calldata certIds,
        bytes32[] calldata certHashes
    ) external view returns (bool[] memory valids, bool[] memory revokeds) {
        uint256 len = certIds.length;
        if (len != certHashes.length) revert ArrayLengthMismatch();
        valids   = new bool[](len);
        revokeds = new bool[](len);
        for (uint256 i = 0; i < len; ++i) {
            if (!_certExists[certIds[i]]) continue;
            CertRecord storage rec = _certs[certIds[i]];
            revokeds[i] = rec.revoked;
            valids[i]   = (!rec.revoked) && (rec.certHash == certHashes[i]);
        }
    }

    // =========================================================
    // 统计信息
    // =========================================================

    /// @notice 返回全局统计数据
    function getStats() external view returns (uint256 totalIssued, uint256 totalRevoked) {
        return (_totalIssued, _totalRevoked);
    }
}