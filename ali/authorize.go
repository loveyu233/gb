package alipay

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
)

// SystemOauthToken 换取授权访问令牌接口 https://docs.open.alipay.com/api_9/alipay.system.oauth.token
func (c *Client) SystemOauthToken(param SystemOauthToken) (result *SystemOauthTokenRsp, err error) {
	err = c.doRequest("POST", param, &result)
	return result, err
}

// UserInfoShare 支付宝会员授权信息查询接口 https://docs.open.alipay.com/api_2/alipay.user.info.share
func (c *Client) UserInfoShare(param UserInfoShare) (result *UserInfoShareRsp, err error) {
	err = c.doRequest("POST", param, &result)
	return result, err
}

// DecodePhoneNumber 小程序获取会员手机号  https://opendocs.alipay.com/mini/api/getphonenumber
// 本方法用于解码小程序端 my.getPhoneNumber 获取的数据
func (c *Client) DecodePhoneNumber(content string) (string, error) {
	// Create cipher block
	block, err := aes.NewCipher(c.encryptKey)
	if err != nil {
		return "", err
	}

	// 128bit zero IV
	iv := make([]byte, 16)
	for i := 0; i < 16; i++ {
		iv[i] = 0
	}

	// Decode the base64 encoded content
	encryptedBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return "", err
	}

	// Check that the encrypted data is a multiple of the block size
	if len(encryptedBytes)%aes.BlockSize != 0 {
		return "", errors.New("ciphertext is not a multiple of the block size")
	}

	// Create CBC mode decryptor
	mode := cipher.NewCBCDecrypter(block, iv)

	// Decrypt the data in-place
	mode.CryptBlocks(encryptedBytes, encryptedBytes)

	// Unpad the decrypted data
	decryptedBytes, err := c.pkcs5Unpad(encryptedBytes)
	if err != nil {
		return "", err
	}

	return string(decryptedBytes), nil
}

// pkcs5Unpad removes PKCS5 padding
func (c *Client) pkcs5Unpad(src []byte) ([]byte, error) {
	length := len(src)
	unpadding := int(src[length-1])

	if unpadding > length {
		return nil, errors.New("unpadding error")
	}

	return src[:(length - unpadding)], nil
}
