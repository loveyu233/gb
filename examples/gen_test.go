package examples

import (
	"github.com/loveyu233/gb"
	"testing"
)

func TestGenConfig(t *testing.T) {
	err := gb.InitGormDB(gb.GormConnConfig{
		Username: "root",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "demo",
	}, gb.GormDefaultLogger(1))
	if err != nil {
		panic(err)
	}

	gb.DB.Gen(gb.WithGenOutFilePath("../gen/query"), gb.WithGenFieldType([]gb.GenFieldType{{
		ColumnName: "json_field",
		ColumnType: "types.User",
		Tags: map[string]string{
			"gorm": "column:json_field;serializer:json",
		},
	}}))
}
