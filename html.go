package gb

import (
	"github.com/k3a/html2text"
)

// HTML2Text 用来将 HTML 文本转换成纯文本。
func HTML2Text(text string) string {
	return html2text.HTML2Text(text)
}
