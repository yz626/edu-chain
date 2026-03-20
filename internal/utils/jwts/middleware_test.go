package jwts

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/yz626/edu-chain/config"
	"github.com/yz626/edu-chain/pkg/constants"
)

var testJWTConfigForMiddleware = &config.JWTConfig{
	Secret:        "test-secret-key-123456",
	Expire:        3600,
	RefreshExpire: 604800,
	Issuer:        "edu-chain-test",
}

// setupGin 设置Gin测试环境
func setupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// TestAuthMiddlewareWithValidToken 测试带有效令牌的认证中间件
func TestAuthMiddlewareWithValidToken(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1, // 正常状态
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	r := setupGin()
	r.Use(jwtInstance.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id":   GetUserID(c),
			"username":  GetUsername(c),
			"email":     GetEmail(c),
			"user_type": GetUserType(c),
			"status":    GetStatus(c),
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusOK, w.Code)
	}
}

// TestAuthMiddlewareWithExpiredToken 测试带过期令牌的认证中间件
func TestAuthMiddlewareWithExpiredToken(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	// 创建一个已过期的令牌
	expiredConfig := &config.JWTConfig{
		Secret:        "test-secret-key-123456",
		Expire:        -3600, // 1小时前过期
		RefreshExpire: 604800,
		Issuer:        "edu-chain-test",
	}
	jwtInstanceExpired := New(expiredConfig)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	token, _ := jwtInstanceExpired.GenerateToken(user)

	r := setupGin()
	r.Use(jwtInstance.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusUnauthorized, w.Code)
	}
}

// TestAuthMiddlewareWithoutToken 测试不带令牌的认证中间件
func TestAuthMiddlewareWithoutToken(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	r := setupGin()
	r.Use(jwtInstance.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusUnauthorized, w.Code)
	}
}

// TestAuthMiddlewareWithInvalidTokenFormat 测试带无效令牌格式的认证中间件
func TestAuthMiddlewareWithInvalidTokenFormat(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	tests := []struct {
		name       string
		authHeader string
		expectCode int
	}{
		{
			name:       "无Bearer前缀",
			authHeader: "some-token",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "错误的prefix",
			authHeader: "Basic some-token",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "空token",
			authHeader: "Bearer ",
			expectCode: http.StatusUnauthorized,
		},
		{
			name:       "多余空格",
			authHeader: "Bearer  token with spaces",
			expectCode: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := setupGin()
			r.Use(jwtInstance.AuthMiddleware())
			r.GET("/test", func(c *gin.Context) {
				c.Status(http.StatusOK)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("Authorization", tt.authHeader)
			w := httptest.NewRecorder()

			r.ServeHTTP(w, req)

			if w.Code != tt.expectCode {
				t.Errorf("期望状态码 %d, 实际状态码 %d", tt.expectCode, w.Code)
			}
		})
	}
}

// TestAuthMiddlewareWithDisabledUser 测试禁用用户的认证中间件
func TestAuthMiddlewareWithDisabledUser(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   2, // 禁用状态
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	r := setupGin()
	r.Use(jwtInstance.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	// 禁用用户应该返回Forbidden
	if w.Code != http.StatusForbidden {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusForbidden, w.Code)
	}
}

// TestAuthMiddlewareWithDifferentUserStatuses 测试不同用户状态的认证中间件
func TestAuthMiddlewareWithDifferentUserStatuses(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	statuses := []int32{1, 2, 3, 4} // 正常、待审核、禁用、锁定
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

		r := setupGin()
		r.Use(jwtInstance.AuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		expectedCode := http.StatusOK
		if status != 1 {
			expectedCode = http.StatusForbidden
		}

		if w.Code != expectedCode {
			t.Errorf("状态 %d: 期望状态码 %d, 实际状态码 %d", status, expectedCode, w.Code)
		}
	}
}

// TestOptionalAuthMiddlewareWithValidToken 测试带有效令牌的可选认证中间件
func TestOptionalAuthMiddlewareWithValidToken(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

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

	r := setupGin()
	r.Use(jwtInstance.OptionalAuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"user_id":   GetUserID(c),
			"username":  GetUsername(c),
			"email":     GetEmail(c),
			"user_type": GetUserType(c),
			"status":    GetStatus(c),
		})
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusOK, w.Code)
	}
}

// TestOptionalAuthMiddlewareWithoutToken 测试不带令牌的可选认证中间件
func TestOptionalAuthMiddlewareWithoutToken(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	r := setupGin()
	r.Use(jwtInstance.OptionalAuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		// 验证没有设置用户信息
		userID := GetUserID(c)
		if userID != "" {
			t.Error("未认证时不应该设置用户ID")
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusOK, w.Code)
	}
}

// TestOptionalAuthMiddlewareWithInvalidToken 测试带无效令牌的可选认证中间件
func TestOptionalAuthMiddlewareWithInvalidToken(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	r := setupGin()
	r.Use(jwtInstance.OptionalAuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		// 验证没有设置用户信息
		userID := GetUserID(c)
		if userID != "" {
			t.Error("无效令牌时不应该设置用户ID")
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid.token")
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusOK, w.Code)
	}
}

// TestOptionalAuthMiddlewareWithDisabledUser 测试带禁用用户令牌的可选认证中间件
func TestOptionalAuthMiddlewareWithDisabledUser(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   2, // 禁用状态
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	r := setupGin()
	r.Use(jwtInstance.OptionalAuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		// 验证没有设置用户信息（因为状态不是1）
		userID := GetUserID(c)
		if userID != "" {
			t.Error("禁用用户不应该设置用户ID")
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusOK, w.Code)
	}
}

// TestGetUserID 测试获取用户ID
func TestGetUserID(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedUserID := "user123"

	c.Set(ContextKeyUserID, expectedUserID)
	userID := GetUserID(c)

	if userID != expectedUserID {
		t.Errorf("期望用户ID %s, 实际 %s", expectedUserID, userID)
	}

	// 测试未设置的情况
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	userID = GetUserID(c2)
	if userID != "" {
		t.Errorf("未设置时应该返回空字符串, 实际返回 %s", userID)
	}
}

// TestGetUsername 测试获取用户名
func TestGetUsername(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedUsername := "testuser"

	c.Set(ContextKeyUsername, expectedUsername)
	username := GetUsername(c)

	if username != expectedUsername {
		t.Errorf("期望用户名 %s, 实际 %s", expectedUsername, username)
	}

	// 测试未设置的情况
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	username = GetUsername(c2)
	if username != "" {
		t.Errorf("未设置时应该返回空字符串, 实际返回 %s", username)
	}
}

// TestGetEmail 测试获取邮箱
func TestGetEmail(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedEmail := "test@example.com"

	c.Set(ContextKeyEmail, expectedEmail)
	email := GetEmail(c)

	if email != expectedEmail {
		t.Errorf("期望邮箱 %s, 实际 %s", expectedEmail, email)
	}

	// 测试未设置的情况
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	email = GetEmail(c2)
	if email != "" {
		t.Errorf("未设置时应该返回空字符串, 实际返回 %s", email)
	}
}

// TestGetUserType 测试获取用户类型
func TestGetUserType(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedUserType := int32(2)

	c.Set(ContextKeyUserType, expectedUserType)
	userType := GetUserType(c)

	if userType != expectedUserType {
		t.Errorf("期望用户类型 %d, 实际 %d", expectedUserType, userType)
	}

	// 测试未设置的情况
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	userType = GetUserType(c2)
	if userType != 0 {
		t.Errorf("未设置时应该返回0, 实际返回 %d", userType)
	}
}

// TestGetStatus 测试获取用户状态
func TestGetStatus(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedStatus := int32(1)

	c.Set(ContextKeyStatus, expectedStatus)
	status := GetStatus(c)

	if status != expectedStatus {
		t.Errorf("期望状态 %d, 实际 %d", expectedStatus, status)
	}

	// 测试未设置的情况
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	status = GetStatus(c2)
	if status != 0 {
		t.Errorf("未设置时应该返回0, 实际返回 %d", status)
	}
}

// TestGetClaims 测试获取Claims
func TestGetClaims(t *testing.T) {
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	expectedClaims := &Claims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	c.Set(ContextKeyClaims, expectedClaims)
	claims := GetClaims(c)

	if claims == nil {
		t.Error("Claims不应该为空")
	}
	if claims.UserID != expectedClaims.UserID {
		t.Errorf("UserID不匹配: 期望 %s, 实际 %s", expectedClaims.UserID, claims.UserID)
	}

	// 测试未设置的情况
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	claims = GetClaims(c2)
	if claims != nil {
		t.Error("未设置时应该返回nil")
	}
}

// TestIsAdmin 测试是否是管理员
func TestIsAdmin(t *testing.T) {
	tests := []struct {
		name        string
		userType    int32
		expectAdmin bool
	}{
		{
			name:        "普通用户",
			userType:    constants.UserTypeNormal,
			expectAdmin: false,
		},
		{
			name:        "管理员",
			userType:    constants.UserTypeAdmin,
			expectAdmin: true,
		},
		{
			name:        "系统用户",
			userType:    constants.UserTypeSystem,
			expectAdmin: false,
		},
		{
			name:        "未知类型",
			userType:    100,
			expectAdmin: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			c.Set(ContextKeyUserType, tt.userType)

			isAdmin := IsAdmin(c)
			if isAdmin != tt.expectAdmin {
				t.Errorf("期望IsAdmin=%v, 实际=%v", tt.expectAdmin, isAdmin)
			}
		})
	}

	// 测试未设置的情况
	c2, _ := gin.CreateTestContext(httptest.NewRecorder())
	if IsAdmin(c2) != false {
		t.Error("未设置时应该返回false")
	}
}

// TestContextKeyConstants 测试上下文键常量
func TestContextKeyConstants(t *testing.T) {
	if ContextKeyUserID != "user_id" {
		t.Errorf("ContextKeyUserID应该为 'user_id'")
	}
	if ContextKeyUsername != "username" {
		t.Errorf("ContextKeyUsername应该为 'username'")
	}
	if ContextKeyEmail != "email" {
		t.Errorf("ContextKeyEmail应该为 'email'")
	}
	if ContextKeyUserType != "user_type" {
		t.Errorf("ContextKeyUserType应该为 'user_type'")
	}
	if ContextKeyStatus != "status" {
		t.Errorf("ContextKeyStatus应该为 'status'")
	}
	if ContextKeyClaims != "jwt_claims" {
		t.Errorf("ContextKeyClaims应该为 'jwt_claims'")
	}
}

// TestAuthMiddlewareChain 测试中间件链
func TestAuthMiddlewareChain(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 2, // 管理员
		Status:   1,
	}

	token, err := jwtInstance.GenerateToken(user)
	if err != nil {
		t.Fatalf("生成令牌失败: %v", err)
	}

	callCount := 0

	r := setupGin()
	r.Use(jwtInstance.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		callCount++
		// 验证所有上下文值都被正确设置
		if GetUserID(c) != user.UserID {
			t.Errorf("UserID不匹配")
		}
		if GetUsername(c) != user.Username {
			t.Errorf("Username不匹配")
		}
		if GetEmail(c) != user.Email {
			t.Errorf("Email不匹配")
		}
		if GetUserType(c) != user.UserType {
			t.Errorf("UserType不匹配")
		}
		if GetStatus(c) != user.Status {
			t.Errorf("Status不匹配")
		}
		claims := GetClaims(c)
		if claims == nil {
			t.Error("Claims不应该为空")
		}
		if !IsAdmin(c) {
			t.Error("IsAdmin应该返回true")
		}
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if callCount != 1 {
		t.Errorf("中间件应该被调用1次, 实际调用 %d 次", callCount)
	}
}

// TestAuthMiddlewareAbort 测试中间件中止
func TestAuthMiddlewareAbort(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	r := setupGin()
	r.Use(jwtInstance.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		t.Error("中间件中止后不应该执行到这里")
	})
	r.GET("/fallback", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// 没有设置Authorization header
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	r.ServeHTTP(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("期望状态码 %d, 实际状态码 %d", http.StatusUnauthorized, w.Code)
	}
}

// TestOptionalAuthMiddlewareWithInvalidFormat 测试带无效格式令牌的可选中间件
func TestOptionalAuthMiddlewareWithInvalidFormat(t *testing.T) {
	jwtInstance := New(testJWTConfigForMiddleware)

	invalidFormats := []string{
		"Basic token",
		"Bearer",
		"Bearer  ",
		"token",
		"",
	}

	for _, format := range invalidFormats {
		r := setupGin()
		r.Use(jwtInstance.OptionalAuthMiddleware())
		r.GET("/test", func(c *gin.Context) {
			// 验证没有设置用户信息
			userID := GetUserID(c)
			if userID != "" {
				t.Errorf("无效格式 '%s' 时不应该设置用户ID", format)
			}
			c.Status(http.StatusOK)
		})

		req := httptest.NewRequest("GET", "/test", nil)
		if format != "" {
			req.Header.Set("Authorization", format)
		}
		w := httptest.NewRecorder()

		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("格式 '%s': 期望状态码 %d, 实际状态码 %d", format, http.StatusOK, w.Code)
		}
	}
}

// BenchmarkAuthMiddleware 基准测试：认证中间件
func BenchmarkAuthMiddleware(b *testing.B) {
	jwtInstance := New(testJWTConfigForMiddleware)

	user := &UserClaims{
		UserID:   "user123",
		Username: "testuser",
		Email:    "test@example.com",
		UserType: 1,
		Status:   1,
	}

	token, _ := jwtInstance.GenerateToken(user)

	r := setupGin()
	r.Use(jwtInstance.AuthMiddleware())
	r.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
	}
}
