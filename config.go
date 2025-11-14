package gb

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// InitConfig 函数用于处理InitConfig相关逻辑。
func InitConfig(fp string, cfg any) (string, error) {
	if fp == "" || !IsPtr(cfg) {
		return "", errors.New("fp为空或cfg非指针")
	}

	env := os.Getenv("GB_ENV")
	if env == "" {
		env = os.Getenv("GO_ENV")
	}
	if env == "" {
		env = "dev"
	}

	file, err := os.ReadFile(fp)
	if err != nil {
		return "", err
	}

	switch filepath.Ext(fp) {
	case ".json":
		return env, json.Unmarshal(file, cfg)
	case ".yml", ".yaml":
		return env, yaml.Unmarshal(file, cfg)
	default:
		return "", errors.New("无效的文件扩展名")
	}
}
