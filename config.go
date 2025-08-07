package gb

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// InitConfig fp为配置文件路径,可以是json文件或者yml和yaml,cfg必须为指针
func InitConfig(fp string, cfg any) error {
	if fp == "" || !IsPtr(cfg) {
		return errors.New("fp为空或cfg非指针")
	}

	file, err := os.ReadFile(fp)
	if err != nil {
		return err
	}

	switch filepath.Ext(fp) {
	case ".json":
		return json.Unmarshal(file, cfg)
	case ".yml", ".yaml":
		return yaml.Unmarshal(file, cfg)
	default:
		return errors.New("无效的文件扩展名")
	}
}
