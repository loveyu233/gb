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
	_, err := gb.InitGormDB(gb.GormConnConfig{
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
	var door []Door
	gb.DB.Scopes(gb.DB.ScopeToday()).Find(&door)
	gb.DB.Scopes(gb.DB.ScopeYesterday()).Find(&door)
	gb.DB.Scopes(gb.DB.ScopeCurrentMonth()).Find(&door)
	gb.DB.Scopes(gb.DB.ScopeLastMonth()).Find(&door)
	gb.DB.Scopes(gb.DB.ScopeNextMonth()).Find(&door)
	gb.DB.Scopes(gb.DB.ScopeCurrentYears()).Find(&door)
	gb.DB.Scopes(gb.DB.ScopeLastYears()).Find(&door)
	gb.DB.Scopes(gb.DB.ScopeNextYears()).Find(&door)

}
