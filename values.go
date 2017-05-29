package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"text/template"

	"github.com/imdario/mergo"
	yaml "gopkg.in/yaml.v2"
)

func loadValues(files []string) (map[string]interface{}, error) {
	var values map[string]interface{}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		var tmp map[string]interface{}
		source, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}

		err = yaml.Unmarshal(source, &tmp)
		mergo.Map(&values, tmp)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse standard input: %v", err)
		}
	}

	return values, nil
}

func calculateValues(files []string) (map[string]interface{}, error) {
	var values map[string]interface{}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		var tmp map[string]interface{}
		content := evaluateFile(file)

		err := yaml.Unmarshal(content.Bytes(), &tmp)

		mergo.Map(&values, tmp)
		if err != nil {
			return nil, fmt.Errorf("Failed to parse standard input: %v", err)
		}

	}

	return values, nil
}

func evaluateFile(templatePath string) *bytes.Buffer {
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template_funcs)

	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("unable to parse template: %s", err)
	}

	buf := new(bytes.Buffer)

	if err := mergo.Merge(&work.Values, work.Env()); err != nil {
		log.Fatalf("template error: %s\n", err)
	}

	err = tmpl.ExecuteTemplate(buf, filepath.Base(templatePath), &work.Values)
	if err != nil {
		log.Fatalf("template error: %s\n", err)
	}

	return buf
}
