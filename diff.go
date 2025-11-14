package gb

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

func DiffMain(text1, text2 string, checklines ...bool) []diffmatchpatch.Diff {
	return diffmatchpatch.New().DiffMain(text1, text2, checklines[0])
}

func DiffPrettyHtml(text1, text2 string, checklines ...bool) string {
	if len(checklines) == 0 {
		checklines[0] = false
	}
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyHtml(dmp.DiffMain(text1, text2, checklines[0]))
}

func DiffPrettyText(text1, text2 string, checklines ...bool) string {
	if len(checklines) == 0 {
		checklines[0] = false
	}
	dmp := diffmatchpatch.New()
	return dmp.DiffPrettyText(dmp.DiffMain(text1, text2, checklines[0]))
}
