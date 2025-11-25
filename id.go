package gb

import (
	"crypto/rand"
	"errors"
	"math/big"
	mrand "math/rand"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/xid"
)

// GetUUID 用来生成一个标准的 UUID 字符串。
func GetUUID() string {
	return uuid.NewString()
}

// GetXID 用来生成紧凑的 XID。
func GetXID() string {
	return xid.New().String()
}

var worker *Worker

var once sync.Once

// GetSnowflakeID 用来产生分布式雪花 ID。
func GetSnowflakeID() int64 {
	once.Do(func() {
		w, err := NewWorker(1)
		if err != nil {
			panic(err)
			return
		}
		worker = w
	})
	return worker.GetId()
}

// RandomString 用来生成指定长度的随机字符串。
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

// RandomStringNoErr 用来快速生成一个 6 位随机字符串。
func RandomStringNoErr() string {
	var charset = RandomCharacterSetAllStr().String()
	var seededRand = mrand.New(mrand.NewSource(Now().UnixNano()))
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomStringWithPrefix 用来生成带前后缀的随机字符串。
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

// String 用来返回字符集的实际内容。
func (r RandomCharacterSet) String() string {
	return string(r)
}

// RandomCharacterSetAllStr 用来返回包含大小写字母和数字的字符集。
func RandomCharacterSetAllStr() RandomCharacterSet {
	return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
}

// RandomCharacterSetLowerStr 用来返回全部小写字母字符集。
func RandomCharacterSetLowerStr() RandomCharacterSet {
	return "abcdefghijklmnopqrstuvwxyz"
}

// RandomCharacterSetLowerStrExcludeCharIO 用来返回去除易混淆字符的全小写字符集。
func RandomCharacterSetLowerStrExcludeCharIO() RandomCharacterSet {
	return "abcdefghjklmnpqrstuvwxyz"
}

// RandomCharacterSetUpperStr 用来返回全部大写字母字符集。
func RandomCharacterSetUpperStr() RandomCharacterSet {
	return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

// RandomCharacterSetUpperStrExcludeCharIO 用来返回去掉 I/O 的大写字符集。
func RandomCharacterSetUpperStrExcludeCharIO() RandomCharacterSet {
	return "ABCDEFGHJKLMNPQRSTUVWXYZ"
}

// RandomCharacterSetNumberStr 用来返回 0-9 的数字字符集。
func RandomCharacterSetNumberStr() RandomCharacterSet {
	return "0123456789"
}

// RandomCharacterSetNumberStrExcludeCharo1 用来返回去除 0/1 的数字字符集。
func RandomCharacterSetNumberStrExcludeCharo1() RandomCharacterSet {
	return "23456789"
}

// RandomCharacterExcludeErrorPronCharacters 用来构建去除易混淆字符的集合。
func RandomCharacterExcludeErrorPronCharacters() RandomCharacterSet {
	return "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjklmnpqrstuvwxyz"
}

// Random 用来从指定字符集中生成随机字符串。
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

// RandomExcludeErrorPronCharacters 用来生成没有易混淆字符的随机串。
func RandomExcludeErrorPronCharacters(strLen int64) string {
	return Random(strLen, RandomCharacterExcludeErrorPronCharacters())
}
