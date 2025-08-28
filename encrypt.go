package gb

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"io"
	"time"
)

type EncryptedResponse struct {
	Data      string `json:"data"`      // 加密的数据
	Timestamp int64  `json:"timestamp"` // 时间戳
	Nonce     string `json:"nonce"`     // 随机数，增加安全性
}

// AES-GCM 加密
func encryptAESGCM(plaintext []byte, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// 生成随机nonce
func generateNonce(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

// EncryptData 加密用户数据,key的长度需要为16位
func EncryptData(data any, custom func(now int64) (key, nonce string)) (*EncryptedResponse, error) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	now := time.Now().Unix()
	key, nonce := custom(now)

	// 加密数据
	encryptedData, err := encryptAESGCM(jsonData, []byte(key))
	if err != nil {
		return nil, err
	}

	return &EncryptedResponse{
		Data:      encryptedData,
		Timestamp: now,
		Nonce:     nonce,
	}, nil
}
