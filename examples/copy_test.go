package examples

import (
	"fmt"
	"testing"
	"time"

	"github.com/loveyu233/gb"
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
	now := gb.Now()
	f := &FromS{
		ID:         1,
		CreatedAt:  &now,
		UsernameGG: "aaa",
	}
	to := ToS{}
	t.Log(gb.DeepCopy(f, &to))
	fmt.Printf("%+v\n", to)
	fmt.Printf("%#v\n", to)
}
