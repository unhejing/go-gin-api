package servicepkg

import (
	"strings"
	"text/template"
)

func ParseTemplateOrPanic(t string) *template.Template {
	tpl, err := template.New("output_template").Parse(t)
	if err != nil {
		panic(err)
	}
	return tpl
}

// SQLColumnToHumpStyle sql转换成驼峰模式
func SQLColumnToHumpStyle(in string) (ret string) {
	for i := 0; i < len(in); i++ {
		if i > 0 && in[i-1] == '_' && in[i] != '_' {
			s := strings.ToUpper(string(in[i]))
			ret += s
		} else if in[i] == '_' {
			continue
		} else {
			ret += string(in[i])
		}
	}
	return
}

func Capitalize(s string) string {
	var upperStr string
	chars := strings.Split(s, "_")
	for _, val := range chars {
		vv := []rune(val)
		for i := 0; i < len(vv); i++ {
			if i == 0 {
				if vv[i] >= 97 && vv[i] <= 122 {
					vv[i] -= 32
				}
				upperStr += string(vv[i])
			} else {
				upperStr += string(vv[i])
			}
		}
	}
	return upperStr
}
