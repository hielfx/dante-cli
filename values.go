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

func calculateValues(files []string) (map[string]interface{}, error) {
	var values map[string]interface{}

	if err := mergo.Merge(&values, work.Env()); err != nil {
		log.Fatalf("template error: %s\n", err)
	}

	for i := len(files) - 1; i >= 0; i-- {
		file := files[i]
		var tmp map[string]interface{}
		//content := evaluateFile(file)

		source, err := ioutil.ReadFile(file)
		if err != nil {
			panic(err)
		}
		err = yaml.Unmarshal(source, &tmp)

		mergo.Map(&values, tmp)
		/*if err := mergo.Merge(&values, &tmp); err != nil {
			log.Fatalf("template error: %s\n", err)
		}*/

	}

	//fmt.Printf("%v OUT\n ", values)

	/*keys := reflect.ValueOf(values).MapKeys()
	strkeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		strkeys[i] = keys[i].String()
		//strvalues[i] =  value.(
	}*/

	for key, value := range values {
		var doc bytes.Buffer
		v := fmt.Sprintf("%s", value)
		k := fmt.Sprintf("%s", key)

		tpl, _ := loadString(k, v)
		tpl.Execute(&doc, values)
		values[key] = doc.String()
	}

	//fmt.Print(strings.Join(strkeys, ","))
	//fmt.Printf("%v OUT\n ", values)

	return values, nil
}

func evaluateFile(templatePath string) *bytes.Buffer {
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template_funcs)

	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("unable to parse template: %s", err)
	}

	buf := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(buf, filepath.Base(templatePath), &work.Values)
	if err != nil {
		log.Fatalf("template error: %s\n", err)
	}

	return buf
}
