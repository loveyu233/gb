package gb

import (
	"golang.org/x/crypto/bcrypt"
	"unicode"
)

// PasswordEncryption 密码加密,返回字符串长度固定为60个字符,password不能超过72字节
func PasswordEncryption(password string) (string, error) {
	fromPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(fromPassword), nil
}

// PasswordCompare 判断密码是否正确
func PasswordCompare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// PasswordValidateStrength 校验密码强度,判断是否同时包含大小写数字和字符,不能有空格
func PasswordValidateStrength(password string, minLen int) bool {
	if len(password) < minLen {
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
