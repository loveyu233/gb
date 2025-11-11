package gb

import (
	"encoding/json"
	"errors"
	"flag"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// InitConfig fp为配置文件路径,可以是json文件或者yml和yaml,cfg必须为指针
func InitConfig(fp string, cfg any) (string, error) {
	if fp == "" || !IsPtr(cfg) {
		return "", errors.New("fp为空或cfg非指针")
	}

	var env string
	flag.StringVar(&env, "e", "dev", "运行环境 (local|test|prod)")
	flag.Parse()

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
