package examples

import (
	"fmt"
	"github.com/loveyu233/gb"
	"testing"
)

type Door struct {
	gb.BaseModel
	ID    int64  `gorm:"column:id"`
	Group string `gorm:"column:group"`
}

func (d Door) TableName() string {
	return "door"
}

func TestCreate(t *testing.T) {
	err := gb.NewGormDB(gb.GormConnConfig{
		Username: "root",
		Password: "",
		Host:     "127.0.0.1",
		Port:     3306,
		Database: "demo",
	}, gb.GormDefaultLogger())
	if err != nil {
		fmt.Println(err)
		return
	}

}
