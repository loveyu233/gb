package gb

import (
	"encoding/json"
	"errors"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
)

// InitConfig fp为配置文件路径,可以是json文件或者yml和yaml,cfg必须为指针
func InitConfig(fp string, cfg any) error {
	if fp == "" || !IsPtr(cfg) {
		return errors.New("fp is not empty or config must be a pointer")
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
		return errors.New("invalid file extension")
	}
}
