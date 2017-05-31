package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
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

	//var secondRound map[string]string{}

	secondRound := map[string]string{}

	for key, value := range values {

		v := fmt.Sprintf("%s", value)
		k := fmt.Sprintf("%s", key)

		var doc bytes.Buffer
		tpl, _ := loadString(k, v)
		tpl.Execute(&doc, values)
		values[key] = doc.String()

		if strings.Contains(v, "{{") {
			secondRound[key] = v
		}

	}

	fmt.Println("Second round...")
	fmt.Printf("%v\n", secondRound)
	for key, value := range secondRound {
		v := fmt.Sprintf("%s", value)
		k := fmt.Sprintf("%s", key)

		var doc bytes.Buffer
		tpl, _ := loadString(k, v)
		tpl.Execute(&doc, values)
		values[key] = doc.String()

		fmt.Printf("%v=%v\n", key, values[key])
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
