package gb

import (
	"crypto/rand"
	"errors"
	"math/big"
	mrand "math/rand"

	"github.com/google/uuid"
	"github.com/loveyu233/gb/snowflake"
	"github.com/rs/xid"
)

// GetUUID 长度为36的字符串
func GetUUID() string {
	return uuid.NewString()
}

// GetXID 长度为20的字符串
func GetXID() string {
	return xid.New().String()
}

// GetSnowflakeID 长度为18的数字
func GetSnowflakeID() int64 {
	return snowflake.GetId()
}

// RandomString 获取指定长度的随机字符串
func RandomString(length int) (string, error) {
	var charset = RandomCharacterSetAllStr().String()
	result := make([]byte, length)
	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}
	return string(result), nil
}

func RandomStringNoErr() string {
	var charset = RandomCharacterSetAllStr().String()
	var seededRand = mrand.New(mrand.NewSource(Now().UnixNano()))
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomStringWithPrefix 带指定前后缀的指定长度字符串
func RandomStringWithPrefix(length int, prefix, suffix string) (string, error) {
	if length <= len(prefix)+len(suffix) {
		return "", errors.New("prefix + suffix <= length")
	}

	randomLength := length - len(prefix) - len(suffix)
	randomPart, err := RandomString(randomLength)
	if err != nil {
		return "", err
	}

	return prefix + randomPart + suffix, nil
}

type RandomCharacterSet string

func (r RandomCharacterSet) String() string {
	return string(r)
}

func RandomCharacterSetAllStr() RandomCharacterSet {
	return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
}

// RandomCharacterSetLowerStr 获取全部小写字母
func RandomCharacterSetLowerStr() RandomCharacterSet {
	return "abcdefghijklmnopqrstuvwxyz"
}

// RandomCharacterSetLowerStrExcludeCharIO 获取全部小写字母，排除掉i和o
func RandomCharacterSetLowerStrExcludeCharIO() RandomCharacterSet {
	return "abcdefghjklmnpqrstuvwxyz"
}

// RandomCharacterSetUpperStr 获取全部大写字母
func RandomCharacterSetUpperStr() RandomCharacterSet {
	return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

// RandomCharacterSetUpperStrExcludeCharIO 获取全部大写字母，排除掉I和O
func RandomCharacterSetUpperStrExcludeCharIO() RandomCharacterSet {
	return "ABCDEFGHJKLMNPQRSTUVWXYZ"
}

// RandomCharacterSetNumberStr 获取全部数字
func RandomCharacterSetNumberStr() RandomCharacterSet {
	return "0123456789"
}

// RandomCharacterSetNumberStrExcludeCharo1 获取全部数字，排除掉1
func RandomCharacterSetNumberStrExcludeCharo1() RandomCharacterSet {
	return "23456789"
}

// RandomCharacterExcludeErrorPronCharacters 排除易看错的字符 01IOio
func RandomCharacterExcludeErrorPronCharacters() RandomCharacterSet {
	return "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjklmnpqrstuvwxyz"
}

// Random 随机返回长度为strLen,随机字符为characterSet,characterSet不写默认为全部字符串
func Random(strLen int64, characterSet ...RandomCharacterSet) string {
	var charset string
	if len(characterSet) == 0 {
		charset = RandomCharacterSetAllStr().String()
	} else {
		for i := range characterSet {
			charset += characterSet[i].String()
		}
	}
	var seededRand = mrand.New(mrand.NewSource(Now().UnixNano()))
	b := make([]byte, strLen)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomExcludeErrorPronCharacters 返回长度为strLen的随机字符串，排除易看错的字符 01IOio
func RandomExcludeErrorPronCharacters(strLen int64) string {
	return Random(strLen, RandomCharacterExcludeErrorPronCharacters())
}
