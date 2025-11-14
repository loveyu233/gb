package gb

import (
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// PasswordEncryption 函数用于处理PasswordEncryption相关逻辑。
func PasswordEncryption(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(fromPassword), nil
}

// PasswordCompare 函数用于处理PasswordCompare相关逻辑。
func PasswordCompare(hashedPassword, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)) == nil
}

// PasswordValidateStrength 函数用于处理PasswordValidateStrength相关逻辑。
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
