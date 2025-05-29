package examples

import (
	"fmt"
	"github.com/loveyu233/gb"
	"testing"
	"time"
)

type FromS struct {
	ID         int64
	CreatedAt  *time.Time
	UsernameGG string
}
type ToS struct {
	ID        int64
	CreatedAt *time.Time
	Username  string `copier:"UsernameGG"`
}

func TestCopy(t *testing.T) {
	now := time.Now()
	f := &FromS{
		ID:         1,
		CreatedAt:  &now,
		UsernameGG: "aaa",
	}
	to := &ToS{}
	t.Log(gb.DeepCopy(f, to))
	fmt.Printf("%+v\n", to)
	fmt.Printf("%#v\n", to)
}
