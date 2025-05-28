package examples

import (
	"fmt"
	"github.com/loveyu233/gb"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

type DemoTest struct {
	gb.BaseModel
	gb.BaseDeleteAtContainsIndex
	Username string `gorm:"index:idx_username;uniqueIndex:deleted_unique_index"`
}

func TestCreate(t *testing.T) {
	var err error
	DB, err := gorm.Open(
		mysql.Open("username:password@tcp(host:3306)/dbname"),
		&gorm.Config{
			Logger:         nil,
			TranslateError: true,
		})
	if err != nil {
		panic(fmt.Errorf("链接数据库失败: %w", err))
	}
	//DB.AutoMigrate(new(DemoTest))
	DB.Create(&DemoTest{Username: "a"})
}
