package examples

import (
	"fmt"
	"github.com/loveyu233/gb"
	"testing"
)

type DemoTest struct {
	gb.BaseModel
	gb.BaseDeleteAtContainsIndex
	Username string `gorm:"index:idx_username;uniqueIndex:deleted_unique_index"`
}

func TestCreate(t *testing.T) {
	db, err := gb.InitGormDB(gb.GormConnConfig{
		Username: "Username",
		Password: "Password",
		Host:     "Host",
		Port:     3306,
		Database: "Database",
		Params:   nil,
	}, gb.GormDefaultLogger())
	if err != nil {
		fmt.Println(err)
		return
	}
	//DB.AutoMigrate(new(DemoTest))
	db.Create(&DemoTest{Username: "b"})
}
