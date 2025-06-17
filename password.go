package gb

import "golang.org/x/crypto/bcrypt"

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
