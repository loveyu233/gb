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

// GetUUID 函数用于处理GetUUID相关逻辑。
func GetUUID() string {
	return uuid.NewString()
}

// GetXID 函数用于处理GetXID相关逻辑。
func GetXID() string {
	return xid.New().String()
}

var worker *Worker

var once sync.Once

// GetSnowflakeID 函数用于处理GetSnowflakeID相关逻辑。
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

// RandomString 函数用于处理RandomString相关逻辑。
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

// RandomStringNoErr 函数用于处理RandomStringNoErr相关逻辑。
func RandomStringNoErr() string {
	var charset = RandomCharacterSetAllStr().String()
	var seededRand = mrand.New(mrand.NewSource(Now().UnixNano()))
	b := make([]byte, 6)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

// RandomStringWithPrefix 函数用于处理RandomStringWithPrefix相关逻辑。
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

// String 方法用于处理String相关逻辑。
func (r RandomCharacterSet) String() string {
	return string(r)
}

// RandomCharacterSetAllStr 函数用于处理RandomCharacterSetAllStr相关逻辑。
func RandomCharacterSetAllStr() RandomCharacterSet {
	return "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
}

// RandomCharacterSetLowerStr 函数用于处理RandomCharacterSetLowerStr相关逻辑。
func RandomCharacterSetLowerStr() RandomCharacterSet {
	return "abcdefghijklmnopqrstuvwxyz"
}

// RandomCharacterSetLowerStrExcludeCharIO 函数用于处理RandomCharacterSetLowerStrExcludeCharIO相关逻辑。
func RandomCharacterSetLowerStrExcludeCharIO() RandomCharacterSet {
	return "abcdefghjklmnpqrstuvwxyz"
}

// RandomCharacterSetUpperStr 函数用于处理RandomCharacterSetUpperStr相关逻辑。
func RandomCharacterSetUpperStr() RandomCharacterSet {
	return "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
}

// RandomCharacterSetUpperStrExcludeCharIO 函数用于处理RandomCharacterSetUpperStrExcludeCharIO相关逻辑。
func RandomCharacterSetUpperStrExcludeCharIO() RandomCharacterSet {
	return "ABCDEFGHJKLMNPQRSTUVWXYZ"
}

// RandomCharacterSetNumberStr 函数用于处理RandomCharacterSetNumberStr相关逻辑。
func RandomCharacterSetNumberStr() RandomCharacterSet {
	return "0123456789"
}

// RandomCharacterSetNumberStrExcludeCharo1 函数用于处理RandomCharacterSetNumberStrExcludeCharo1相关逻辑。
func RandomCharacterSetNumberStrExcludeCharo1() RandomCharacterSet {
	return "23456789"
}

// RandomCharacterExcludeErrorPronCharacters 函数用于处理RandomCharacterExcludeErrorPronCharacters相关逻辑。
func RandomCharacterExcludeErrorPronCharacters() RandomCharacterSet {
	return "23456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghjklmnpqrstuvwxyz"
}

// Random 函数用于处理Random相关逻辑。
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

// RandomExcludeErrorPronCharacters 函数用于处理RandomExcludeErrorPronCharacters相关逻辑。
func RandomExcludeErrorPronCharacters(strLen int64) string {
	return Random(strLen, RandomCharacterExcludeErrorPronCharacters())
}
