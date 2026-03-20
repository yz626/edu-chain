package crypto

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "TestPassword123"

	// 测试密码哈希
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	if hash == password {
		t.Fatal("Hash should not be equal to password")
	}

	// 测试空密码
	_, err = HashPassword("")
	if err != ErrPasswordEmpty {
		t.Fatalf("Expected ErrPasswordEmpty for empty password, got: %v", err)
	}
}

func TestHashPasswordWithCost(t *testing.T) {
	password := "TestPassword123"

	// 测试有效成本
	hash, err := HashPasswordWithCost(password, 10)
	if err != nil {
		t.Fatalf("HashPasswordWithCost failed: %v", err)
	}
	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	// 测试无效成本（低于最小值）
	hash, err = HashPasswordWithCost(password, 2)
	if err != nil {
		t.Fatalf("HashPasswordWithCost with low cost failed: %v", err)
	}
	// 应该使用默认成本
	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	// 测试无效成本（高于最大值）
	hash, err = HashPasswordWithCost(password, 100)
	if err != nil {
		t.Fatalf("HashPasswordWithCost with high cost failed: %v", err)
	}
	if hash == "" {
		t.Fatal("Hash should not be empty")
	}

	// 测试空密码
	_, err = HashPasswordWithCost("", 10)
	if err != ErrPasswordEmpty {
		t.Fatalf("Expected ErrPasswordEmpty for empty password, got: %v", err)
	}
}

func TestCheckPassword(t *testing.T) {
	password := "TestPassword123"

	// 生成哈希
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 测试正确密码
	match, err := CheckPassword(password, hash)
	if err != nil {
		t.Fatalf("CheckPassword failed: %v", err)
	}
	if !match {
		t.Fatal("Password should match")
	}

	// 测试错误密码
	match, err = CheckPassword("WrongPassword", hash)
	if err != nil {
		t.Fatalf("CheckPassword failed: %v", err)
	}
	if match {
		t.Fatal("Password should not match")
	}

	// 测试空密码
	_, err = CheckPassword("", hash)
	if err != ErrPasswordEmpty {
		t.Fatalf("Expected ErrPasswordEmpty for empty password, got: %v", err)
	}

	// 测试空哈希
	_, err = CheckPassword(password, "")
	if err != ErrPasswordEmpty {
		t.Fatalf("Expected ErrPasswordEmpty for empty hash, got: %v", err)
	}
}

func TestValidatePasswordStrength(t *testing.T) {
	tests := []struct {
		password   string
		wantErr    bool
		wantStrong bool // 是否期望强密码
	}{
		{"", true, false},
		{"abc", true, false},
		{"abcdefgh", true, false},               // 没有大写
		{"ABCDEFGH", true, false},               // 没有小写
		{"abcd1234", true, false},               // 没有大写字母
		{"Abcd1234!", false, true},              // 有效密码
		{"Password123!", false, true},           // 强密码
		{"VeryLongPassword123!@#", false, true}, // 非常强密码
		{"p@ssw0rd", true, false},               // 太短
	}

	for _, tt := range tests {
		strength, err := ValidatePasswordStrength(tt.password)
		if (err != nil) != tt.wantErr {
			t.Errorf("ValidatePasswordStrength(%q) error = %v, wantErr %v",
				tt.password, err, tt.wantErr)
			continue
		}

		if !tt.wantErr && tt.wantStrong && strength < StrengthStrong {
			t.Logf("Password %q strength: %d (expected strong)", tt.password, strength)
		}
	}
}

func TestValidatePassword(t *testing.T) {
	// 测试有效密码
	err := ValidatePassword("Abcd1234!")
	if err != nil {
		t.Errorf("ValidatePassword should accept valid password: %v", err)
	}

	// 测试无效密码（太短）
	err = ValidatePassword("Abc1!")
	if err != ErrPasswordTooShort {
		t.Errorf("Expected ErrPasswordTooShort, got: %v", err)
	}

	// 测试无效密码（没有大写）
	err = ValidatePassword("abcd1234!")
	if err != ErrPasswordNoUpperCase {
		t.Errorf("Expected ErrPasswordNoUpperCase, got: %v", err)
	}

	// 测试无效密码（没有小写）
	err = ValidatePassword("ABCD1234!")
	if err != ErrPasswordNoLowerCase {
		t.Errorf("Expected ErrPasswordNoLowerCase, got: %v", err)
	}

	// 测试无效密码（没有数字）
	err = ValidatePassword("Abcdefgh!")
	if err != ErrPasswordNoDigit {
		t.Errorf("Expected ErrPasswordNoDigit, got: %v", err)
	}
}

func TestValidatePasswordWithCustomRules(t *testing.T) {
	// 使用自定义验证规则（不需要特殊字符）
	validator := PasswordValidator{
		MinLength:      6,
		MaxLength:      20,
		RequireUpper:   true,
		RequireLower:   true,
		RequireDigit:   true,
		RequireSpecial: false,
	}

	// 测试有效密码
	err := ValidatePasswordWithCustomRules("Abcd12", validator)
	if err != nil {
		t.Errorf("ValidatePasswordWithCustomRules should accept valid password: %v", err)
	}

	// 测试太短的密码
	err = ValidatePasswordWithCustomRules("Ab1", validator)
	if err != ErrPasswordTooShort {
		t.Errorf("Expected ErrPasswordTooShort, got: %v", err)
	}
}

func TestGetPasswordStrengthText(t *testing.T) {
	tests := []struct {
		strength int
		want     string
	}{
		{StrengthVeryWeak, "非常弱"},
		{StrengthWeak, "弱"},
		{StrengthFair, "一般"},
		{StrengthStrong, "强"},
		{StrengthVeryStrong, "非常强"},
		{99, "未知"},
	}

	for _, tt := range tests {
		got := GetPasswordStrengthText(tt.strength)
		if got != tt.want {
			t.Errorf("GetPasswordStrengthText(%d) = %v, want %v",
				tt.strength, got, tt.want)
		}
	}
}

func TestIsPasswordHashValid(t *testing.T) {
	// 生成有效的 bcrypt 哈希
	password := "TestPassword123"
	hash, _ := HashPassword(password)

	// 测试有效哈希
	if !IsPasswordHashValid(hash) {
		t.Error("Valid bcrypt hash should pass validation")
	}

	// 测试空哈希
	if IsPasswordHashValid("") {
		t.Error("Empty hash should be invalid")
	}

	// 测试无效格式
	if IsPasswordHashValid("invalid_hash") {
		t.Error("Invalid format should fail validation")
	}

	// 测试错误长度的哈希
	if IsPasswordHashValid("$2a$10$1234567890123456789012") {
		t.Error("Wrong length hash should be invalid")
	}
}

func TestNeedsRehash(t *testing.T) {
	// 生成高成本哈希
	password := "TestPassword123"
	hashHighCost, _ := HashPasswordWithCost(password, 15)

	// 生成低成本的哈希用于测试
	hashLowCost, _ := HashPasswordWithCost(password, 5)

	// 测试高成本哈希不需要重新哈希
	needsRehash, err := NeedsRehash(hashHighCost, 10)
	if err != nil {
		t.Fatalf("NeedsRehash failed: %v", err)
	}
	if needsRehash {
		t.Error("High cost hash should not need rehash")
	}

	// 测试低成本的哈希需要重新哈希
	needsRehash, err = NeedsRehash(hashLowCost, 10)
	if err != nil {
		t.Fatalf("NeedsRehash failed: %v", err)
	}
	if !needsRehash {
		t.Error("Low cost hash should need rehash")
	}

	// 测试空哈希
	_, err = NeedsRehash("", 10)
	if err != ErrPasswordEmpty {
		t.Fatalf("Expected ErrPasswordEmpty for empty hash, got: %v", err)
	}
}

func TestGenerateRandomPassword(t *testing.T) {
	// 测试默认长度
	password, err := GenerateRandomPassword(0, false)
	if err != nil {
		t.Fatalf("GenerateRandomPassword failed: %v", err)
	}
	if len(password) != 16 {
		t.Errorf("Default password length should be 16, got %d", len(password))
	}

	// 测试指定长度
	password, err = GenerateRandomPassword(12, false)
	if err != nil {
		t.Fatalf("GenerateRandomPassword failed: %v", err)
	}
	if len(password) != 12 {
		t.Errorf("Password length should be 12, got %d", len(password))
	}

	// 测试包含特殊字符
	password, err = GenerateRandomPassword(12, true)
	if err != nil {
		t.Fatalf("GenerateRandomPassword failed: %v", err)
	}
	if len(password) != 12 {
		t.Errorf("Password length should be 12, got %d", len(password))
	}

	// 测试无效长度（太短）
	password, err = GenerateRandomPassword(4, false)
	if err != nil {
		t.Fatalf("GenerateRandomPassword failed: %v", err)
	}
	// 应该使用默认长度
	if len(password) != 16 {
		t.Errorf("Password length should be 16, got %d", len(password))
	}

	// 测试无效长度（太长）
	password, err = GenerateRandomPassword(100, false)
	if err != nil {
		t.Fatalf("GenerateRandomPassword failed: %v", err)
	}
	// 应该使用默认长度
	if len(password) != 16 {
		t.Errorf("Password length should be 16, got %d", len(password))
	}
}

func TestBcryptConsistency(t *testing.T) {
	password := "ConsistentPassword123"

	// 多次哈希同一个密码应该产生不同的结果（因为盐值不同）
	hash1, _ := HashPassword(password)
	hash2, _ := HashPassword(password)

	if hash1 == hash2 {
		t.Error("Same password should produce different hashes due to salt")
	}

	// 但是两个哈希都应该能验证成功
	match1, _ := CheckPassword(password, hash1)
	match2, _ := CheckPassword(password, hash2)

	if !match1 || !match2 {
		t.Error("Both hashes should validate the password correctly")
	}
}
