package gb

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffMain 用来使用 diffmatchpatch 计算两个文本的差异。
func DiffMain(text1, text2 string, checklines ...bool) []diffmatchpatch.Diff {
	return diffmatchpatch.New().DiffMain(text1, text2, checklines[0])
}

// DiffPrettyHtml 用来返回高亮显示差异的 HTML 字符串。
func DiffPrettyHtml(text1, text2 string, checklines ...bool) string {
	if len(checklines) == 0 {
		checklines = []bool{false}
	}
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyHtml(dmp.DiffMain(text1, text2, checklines[0]))
}

// DiffPrettyText 用来返回可读的纯文本差异结果。
func DiffPrettyText(text1, text2 string, checklines ...bool) string {
	if len(checklines) == 0 {
		checklines = []bool{false}
	}
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(dmp.DiffMain(text1, text2, checklines[0]))
}
