package gb

import (
	"bytes"
	"html/template"
)

// TemplateReplace 函数用于处理TemplateReplace相关逻辑。
func TemplateReplace(tmp string, replace any) (string, error) {
	tpl, err := template.New("fee").Parse(tmp)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err = tpl.Execute(buf, replace); err != nil {
		return "", err
	}

	return buf.String(), nil
}
