package gb

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordEncryption 用来对明文密码进行 bcrypt 加密。
func PasswordEncryption(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(fromPassword), nil
}

// PasswordCompare 用来比较加密密码与输入密码。
func PasswordCompare(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// PasswordValidateStrength 用来检测密码的长度和复杂度是否合规。
func PasswordValidateStrength(password string, minLen, maxLen int) bool {
	if len(password) < minLen || len(password) > maxLen {
		return false
	}

	var (
		hasUpper   bool
		hasLower   bool
		hasDigit   bool
		hasSpecial bool
	)

	for _, char := range password {
		if char == ' ' { // 如果包含空格，直接返回 false
			return false
		}

		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsDigit(char):
			hasDigit = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	return hasUpper && hasLower && hasDigit && hasSpecial
}
