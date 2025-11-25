package gb

import (
	"bytes"
	"html/template"
)

// TemplateReplace 用来将字符串模板渲染成实际内容。
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
