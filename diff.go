package gb

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

// DiffMain 函数用于处理DiffMain相关逻辑。
func DiffMain(text1, text2 string, checklines ...bool) []diffmatchpatch.Diff {
	return diffmatchpatch.New().DiffMain(text1, text2, checklines[0])
}

// DiffPrettyHtml 函数用于处理DiffPrettyHtml相关逻辑。
func DiffPrettyHtml(text1, text2 string, checklines ...bool) string {
	if len(checklines) == 0 {
		checklines = []bool{false}
	}
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyHtml(dmp.DiffMain(text1, text2, checklines[0]))
}

// DiffPrettyText 函数用于处理DiffPrettyText相关逻辑。
func DiffPrettyText(text1, text2 string, checklines ...bool) string {
	if len(checklines) == 0 {
		checklines = []bool{false}
	}
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(dmp.DiffMain(text1, text2, checklines[0]))
}
