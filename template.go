package gb

import (
	"bytes"
	"html/template"
)

// TemplateReplace tmp中被替换的格式为 {{.key}}, replace最好为map,map[key]=value
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
