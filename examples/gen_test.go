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
		Database: "demo2",
	}, gb.GormDefaultLogger())
	if err != nil {
		panic(err)
	}

	gb.DB.Gen(gb.WithGenOutFilePath("../gen/query"), gb.WithGenUseTablesName("simple_table"))

}
