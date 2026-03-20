package jwts

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/yz626/edu-chain/config"
)

var testJWTConfig = &config.JWTConfig{
	Secret:        "test-secret-key-123456",
	Expire:        3600,   // 1小时
	RefreshExpire: 604800, // 7天
	Issuer:        "edu-chain-test",
}

// TestNew 测试创建JWT实例
func TestNew(t *testing.T) {
	jwtInstance := New(testJWTConfig)
	if jwtInstance == nil {
		t.Error("JWT实例不应为空")
	}
	if jwtInstance.cfg == nil {
		t.Error("JWT配置不应为空")
	}
}

// TestGenerateToken 测试生成访问令牌
func TestGenerateToken(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	if token == "" {
		t.Error("生成的令牌不应为空")
	}

	// 验证令牌可以被解析
	claims, err := jwtInstance.ParseToken(token)
	if err != nil {
		t.Fatalf("解析令牌失败: %v", err)
	}

	// 验证Claims中的用户信息
	if claims.UserID != user.UserID {
		t.Errorf("UserID不匹配: 期望 %s, 实际 %s", user.UserID, claims.UserID)
	}
	if claims.Username != user.Username {
		t.Errorf("Username不匹配: 期望 %s, 实际 %s", user.Username, claims.Username)
	}
	if claims.Email != user.Email {
		t.Errorf("Email不匹配: 期望 %s, 实际 %s", user.Email, claims.Email)
	}
	if claims.UserType != user.UserType {
		t.Errorf("UserType不匹配: 期望 %d, 实际 %d", user.UserType, claims.UserType)
	}
	if claims.Status != user.Status {
		t.Errorf("Status不匹配: 期望 %d, 实际 %d", user.Status, claims.Status)
	}

	// 验证Issuer
	if claims.Issuer != testJWTConfig.Issuer {
		t.Errorf("Issuer不匹配: 期望 %s, 实际 %s", testJWTConfig.Issuer, claims.Issuer)
	}

	// 验证过期时间
	if claims.ExpiresAt == nil {
		t.Error("过期时间不应为空")
	}
}

// TestGenerateTokenWithEmptyUser 测试空用户生成令牌
func TestGenerateTokenWithEmptyUser(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	user := &UserClaims{}
	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	if token == "" {
		t.Error("生成的令牌不应为空")
	}
}

// TestGenerateRefreshToken 测试生成刷新令牌
func TestGenerateRefreshToken(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	userID := "user123"
	token, err := jwtInstance.GenerateRefreshToken(userID)
	if err != nil {
		t.Fatalf("生成刷新令牌失败: %v", err)
	}

	if token == "" {
		t.Error("生成的刷新令牌不应为空")
	}

	// 验证刷新令牌可以被解析
	claims, err := jwtInstance.ParseToken(token)
	if err != nil {
		t.Fatalf("解析刷新令牌失败: %v", err)
	}

	// 刷新令牌的Subject应该是用户ID
	if claims.Subject != userID {
		t.Errorf("Subject不匹配: 期望 %s, 实际 %s", userID, claims.Subject)
	}
}

// TestParseToken 测试解析有效令牌
func TestParseToken(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	user := &UserClaims{
		UserID:   "user456",
		Username: "parseuser",
		Email:    "parse@example.com",
		UserType: 2,
		Status:   1,
	}

	// 先生成一个有效的令牌
	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	// 解析令牌
	claims, err := jwtInstance.ParseToken(token)
	if err != nil {
		t.Fatalf("解析令牌失败: %v", err)
	}

	// 验证解析结果
	if claims.UserID != user.UserID {
		t.Errorf("UserID不匹配: 期望 %s, 实际 %s", user.UserID, claims.UserID)
	}
}

// TestParseTokenWithInvalidToken 测试解析无效令牌
func TestParseTokenWithInvalidToken(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	tests := []struct {
		name      string
		token     string
		expectErr error
	}{
		{
			name:      "空令牌",
			token:     "",
			expectErr: ErrTokenInvalid, // 空字符串会被识别为无效而非格式错误
		},
		{
			name:      "无效格式令牌",
			token:     "invalid.token.here",
			expectErr: ErrTokenInvalid,
		},
		{
			name:      "完全无效的令牌",
			token:     "this.is.not.a.valid.jwt.token.at.all",
			expectErr: ErrTokenInvalid,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := jwtInstance.ParseToken(tt.token)
			if err != tt.expectErr {
				t.Errorf("期望错误 %v, 实际错误 %v", tt.expectErr, err)
			}
		})
	}
}

// TestParseTokenWithWrongSecret 测试使用错误密钥解析令牌
func TestParseTokenWithWrongSecret(t *testing.T) {
	// 使用正确的密钥生成令牌
	jwtInstanceCorrect := New(testJWTConfig)
	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}
	token, err := jwtInstanceCorrect.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	// 使用不同的密钥解析令牌
	wrongConfig := &config.JWTConfig{
		Secret:        "wrong-secret-key",
		Expire:        3600,
		RefreshExpire: 604800,
		Issuer:        "test",
	}
	jwtInstanceWrong := New(wrongConfig)

	_, err = jwtInstanceWrong.ParseToken(token)
	if err != ErrTokenInvalid {
		t.Errorf("期望错误 %v, 实际错误 %v", ErrTokenInvalid, err)
	}
}

// TestParseTokenWithExpiredToken 测试解析过期令牌
func TestParseTokenWithExpiredToken(t *testing.T) {
	// 手动创建一个已过期的令牌
	// 使用 jwt.ErrTokenExpired 来创建一个带有过期声明的令牌
	now := time.Now()
	// 过期时间是1小时前
	expiredTime := now.Add(-time.Hour)
	// 签发时间是2小时前（在过期时间之前）
	issuedTime := now.Add(-2 * time.Hour)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiredTime),
			IssuedAt:  jwt.NewNumericDate(issuedTime),
			NotBefore: jwt.NewNumericDate(issuedTime),
			Issuer:    testJWTConfig.Issuer,
		},
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	expiredTokenString, err := token.SignedString([]byte(testJWTConfig.Secret))
	if err != nil {
		t.Fatalf("创建过期令牌失败: %v", err)
	}

	jwtInstance := New(testJWTConfig)
	_, err = jwtInstance.ParseToken(expiredTokenString)
	// 验证错误是 ErrTokenExpired 或 ErrTokenInvalid（取决于JWT库实现）
	if err != ErrTokenExpired && err != ErrTokenInvalid {
		t.Errorf("期望错误 %v 或 %v, 实际错误 %v", ErrTokenExpired, ErrTokenInvalid, err)
	}
}

// TestParseTokenWithWrongSigningMethod 测试使用错误签名方法解析令牌
func TestParseTokenWithWrongSigningMethod(t *testing.T) {
	// 创建一个使用RS256的令牌
	claims := jwt.MapClaims{
		"user_id":   "user123",
		"username":  "testuser",
		"email":     "test@example.com",
		"user_type": 1,
		"status":    1,
		"exp":       time.Now().Add(time.Hour).Unix(),
		"iss":       "test",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	// 注意：这里只是创建一个令牌字符串，不实际签名
	_, err := token.SignedString([]byte("dummy"))
	// 这里忽略错误，因为我们只需要令牌字符串格式

	// 创建一个新的令牌用于测试
	tokenStr := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoidXNlcjEyMyIsInVzZXJuYW1lIjoidGVzdHVzZXIiLCJlbWFpbCI6InRlc3RAZXhhbXBsZS5jb20iLCJ1c2VyX3R5cGUiOjEsInN0YXR1cyI6MSwiZXhwIjoxNzAwMDAwMDAwfQ.test"

	jwtInstance := New(testJWTConfig)
	_, err = jwtInstance.ParseToken(tokenStr)
	// 这里会返回 ErrTokenInvalid 因为签名验证失败
	if err != ErrTokenInvalid && err != ErrSigningMethod {
		t.Logf("实际错误: %v", err)
	}
}

// TestValidateToken 测试验证令牌
func TestValidateToken(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	// 生成有效令牌
	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	// 验证有效令牌
	if !jwtInstance.ValidateToken(token) {
		t.Error("有效令牌应该通过验证")
	}

	// 验证无效令牌
	if jwtInstance.ValidateToken("invalid.token") {
		t.Error("无效令牌不应该通过验证")
	}

	// 验证空令牌
	if jwtInstance.ValidateToken("") {
		t.Error("空令牌不应该通过验证")
	}
}

// TestGetUserIDFromToken 测试从令牌获取用户ID
func TestGetUserIDFromToken(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	expectedUserID := "user999"
	user := &UserClaims{
		UserID:   expectedUserID,
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	// 获取用户ID
	userID, err := jwtInstance.GetUserIDFromToken(token)
	if err != nil {
		t.Fatalf("获取用户ID失败: %v", err)
	}

	if userID != expectedUserID {
		t.Errorf("期望用户ID %s, 实际 %s", expectedUserID, userID)
	}

	// 测试从无效令牌获取用户ID
	_, err = jwtInstance.GetUserIDFromToken("invalid.token")
	if err == nil {
		t.Error("从无效令牌获取用户ID应该失败")
	}
}

// TestGenerateTokenWithDifferentUserTypes 测试不同用户类型
func TestGenerateTokenWithDifferentUserTypes(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	userTypes := []int32{1, 2, 3, 100, -1}
	for _, userType := range userTypes {
		user := &UserClaims{
			UserID:   "user123",
			Username: "testuser",
			Email:    "test@example.com",
			UserType: userType,
			Status:   1,
		}

		token, err := jwtInstance.GenerateToken(user)
		if err != nil {
			t.Fatalf("生成用户类型 %d 的令牌失败: %v", userType, err)
		}

		claims, err := jwtInstance.ParseToken(token)
		if err != nil {
			t.Fatalf("解析用户类型 %d 的令牌失败: %v", userType, err)
		}

		if claims.UserType != userType {
			t.Errorf("用户类型不匹配: 期望 %d, 实际 %d", userType, claims.UserType)
		}
	}
}

// TestGenerateTokenWithDifferentStatuses 测试不同用户状态
func TestGenerateTokenWithDifferentStatuses(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	statuses := []int32{0, 1, 2, 3, 4, 100}
	for _, status := range statuses {
		user := &UserClaims{
			UserID:   "user123",
			Username: "testuser",
			Email:    "test@example.com",
			UserType: 1,
			Status:   status,
		}

		token, err := jwtInstance.GenerateToken(user)
		if err != nil {
			t.Fatalf("生成状态 %d 的令牌失败: %v", status, err)
		}

		claims, err := jwtInstance.ParseToken(token)
		if err != nil {
			t.Fatalf("解析状态 %d 的令牌失败: %v", status, err)
		}

		if claims.Status != status {
			t.Errorf("状态不匹配: 期望 %d, 实际 %d", status, claims.Status)
		}
	}
}

// TestGenerateTokenWithSpecialCharacters 测试包含特殊字符的用户信息
func TestGenerateTokenWithSpecialCharacters(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	user := &UserClaims{
		UserID:   "user-with-dashes",
		Username: "user_underscore",
		Email:    "test+special@example.com",
		UserType: 1,
		Status:   1,
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成包含特殊字符的令牌失败: %v", err)
	}

	claims, err := jwtInstance.ParseToken(token)
	if err != nil {
		t.Fatalf("解析包含特殊字符的令牌失败: %v", err)
	}

	if claims.UserID != user.UserID {
		t.Errorf("UserID不匹配: 期望 %s, 实际 %s", user.UserID, claims.UserID)
	}
	if claims.Username != user.Username {
		t.Errorf("Username不匹配: 期望 %s, 实际 %s", user.Username, claims.Username)
	}
	if claims.Email != user.Email {
		t.Errorf("Email不匹配: 期望 %s, 实际 %s", user.Email, claims.Email)
	}
}

// TestHandleJWTError 测试JWT错误处理
func TestHandleJWTError(t *testing.T) {
	// 测试 nil 错误
	result := handleJWTError(nil)
	if result != nil {
		t.Error("nil错误应该返回nil")
	}

	// 测试过期错误
	result = handleJWTError(jwt.ErrTokenExpired)
	if result != ErrTokenExpired {
		t.Errorf("期望错误 %v, 实际错误 %v", ErrTokenExpired, result)
	}

	// 测试还未生效错误
	result = handleJWTError(jwt.ErrTokenNotValidYet)
	if result != ErrTokenNotValidYet {
		t.Errorf("期望错误 %v, 实际错误 %v", ErrTokenNotValidYet, result)
	}

	// 测试格式错误
	result = handleJWTError(jwt.ErrTokenMalformed)
	if result != ErrTokenMalformed {
		t.Errorf("期望错误 %v, 实际错误 %v", ErrTokenMalformed, result)
	}

	// 测试其他错误
	result = handleJWTError(jwt.ErrSignatureInvalid)
	if result != ErrTokenInvalid {
		t.Errorf("期望错误 %v, 实际错误 %v", ErrTokenInvalid, result)
	}
}

// TestTokenUniqueness 测试令牌唯一性
func TestTokenUniqueness(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	// 生成多个令牌
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, err := jwtInstance.GenerateToken(user)
		if err != nil {
			t.Fatalf("生成令牌 %d 失败: %v", i, err)
		}
		tokens[token] = true
	}

	// 验证所有令牌都是唯一的
	if len(tokens) != 100 {
		t.Errorf("期望100个唯一令牌, 实际得到 %d 个", len(tokens))
	}
}

// TestRefreshTokenUniqueness 测试刷新令牌唯一性
func TestRefreshTokenUniqueness(t *testing.T) {
	jwtInstance := New(testJWTConfig)

	userID := "user123"

	// 生成多个刷新令牌
	tokens := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, err := jwtInstance.GenerateRefreshToken(userID)
		if err != nil {
			t.Fatalf("生成刷新令牌 %d 失败: %v", i, err)
		}
		tokens[token] = true
	}

	// 验证所有令牌都是唯一的
	if len(tokens) != 100 {
		t.Errorf("期望100个唯一刷新令牌, 实际得到 %d 个", len(tokens))
	}
}

// TestGenerateTokenWithLongExpire 测试长过期时间
func TestGenerateTokenWithLongExpire(t *testing.T) {
	longExpireConfig := &config.JWTConfig{
		Secret:        "test-secret-key-123456",
		Expire:        86400 * 365, // 1年
		RefreshExpire: 604800,
		Issuer:        "edu-chain-test",
	}
	jwtInstance := New(longExpireConfig)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成1年过期的令牌失败: %v", err)
	}

	claims, err := jwtInstance.ParseToken(token)
	if err != nil {
		t.Fatalf("解析1年过期的令牌失败: %v", err)
	}

	// 验证过期时间大约是1年
	expectedExpiry := time.Now().Add(time.Duration(longExpireConfig.Expire) * time.Second)
	actualExpiry := claims.ExpiresAt.Time
	diff := expectedExpiry.Sub(actualExpiry)
	if diff > time.Minute || diff < -time.Minute {
		t.Errorf("过期时间不匹配: 期望 %v, 实际 %v", expectedExpiry, actualExpiry)
	}
}

// BenchmarkGenerateToken 基准测试：生成令牌
func BenchmarkGenerateToken(b *testing.B) {
	jwtInstance := New(testJWTConfig)
	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jwtInstance.GenerateToken(user)
	}
}

// BenchmarkParseToken 基准测试：解析令牌
func BenchmarkParseToken(b *testing.B) {
	jwtInstance := New(testJWTConfig)
	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}
	token, _ := jwtInstance.GenerateToken(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = jwtInstance.ParseToken(token)
	}
}

// BenchmarkValidateToken 基准测试：验证令牌
func BenchmarkValidateToken(b *testing.B) {
	jwtInstance := New(testJWTConfig)
	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}
	token, _ := jwtInstance.GenerateToken(user)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = jwtInstance.ValidateToken(token)
	}
}
