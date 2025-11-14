package gb

import (
	"github.com/k3a/html2text"
)

func HTML2Text(text string) string {
	return html2text.HTML2Text(text)
}
