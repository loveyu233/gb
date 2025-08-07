package examples

import (
	"testing"

	"github.com/loveyu233/gb"
)

func TestTemplateReplace(t *testing.T) {
	str, err := gb.TemplateReplace("<h1>{{.name}},{{.age}}</h1>", map[string]any{
		"name": "张三",
		"age":  30,
	})
	if err != nil {
		t.Fatal(err)
		return
	}
	t.Log(str)
}
