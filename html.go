package gb

import (
	"github.com/k3a/html2text"
)

// HTML2Text 函数用于处理HTML2Text相关逻辑。
func HTML2Text(text string) string {
	return html2text.HTML2Text(text)
}
