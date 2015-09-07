package main

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
)

//	模板字典
var templateMap map[string]*template.Template = make(map[string]*template.Template)

//	注册模板
func registerTemplates(layout string, templates ...string) {
	for _, t := range templates {
		templateMap[t] = template.Must(template.ParseFiles(layout, t))
	}
}

//	执行模板
func processTemplate(wr io.Writer, name string, data interface{}) error {
	t, found := templateMap[name]
	if !found {
		return fmt.Errorf("模板%s未注册", name)
	}

	return t.ExecuteTemplate(wr, filepath.Base(name), data)
}
