package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"text/template"

	sh "github.com/codeskyblue/go-sh"
	"github.com/imdario/mergo"
	"github.com/jwilder/gojq"
)

var template_funcs = template.FuncMap{
	"contains":  contains,
	"exists":    exists,
	"split":     split,
	"replace":   strings.Replace,
	"default":   defaultValue,
	"parseUrl":  parseUrl,
	"atoi":      strconv.Atoi,
	"add":       add,
	"isTrue":    isTrue,
	"lower":     strings.ToLower,
	"upper":     strings.ToUpper,
	"jsonQuery": jsonQuery,
	"shell":     shell,
	"commit":    commit,
	"branch":    branch,
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func contains(item map[string]string, key string) bool {
	if _, ok := item[key]; ok {
		return true
	}
	return false
}

func defaultValue(args ...interface{}) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("default called with no values!")
	}

	if len(args) > 0 {
		if args[0] != nil {
			return args[0].(string), nil
		}
	}

	if len(args) > 1 {
		if args[1] == nil {
			return "", fmt.Errorf("default called with nil default value!")
		}

		if _, ok := args[1].(string); !ok {
			return "", fmt.Errorf("default is not a string value. hint: surround it w/ double quotes.")
		}

		return args[1].(string), nil
	}

	return "", fmt.Errorf("default called with no default value")
}

func parseUrl(rawurl string) *url.URL {
	u, err := url.Parse(rawurl)
	if err != nil {
		log.Fatalf("unable to parse url %s: %s", rawurl, err)
	}
	return u
}

func add(arg1, arg2 int) int {
	return arg1 + arg2
}

func isTrue(s string) bool {
	b, err := strconv.ParseBool(strings.ToLower(s))
	if err == nil {
		return b
	}
	return false
}

func jsonQuery(jsonObj string, query string) (interface{}, error) {
	parser, err := gojq.NewStringQuery(jsonObj)
	if err != nil {
		return "", err
	}
	res, err := parser.Query(query)
	if err != nil {
		return "", err
	}
	return res, nil
}

func shell(cmd string, a ...interface{}) (string, error) {
	out, err := sh.Command(cmd, a...).Output()

	s := fmt.Sprintf("%s", out)
	s = strings.TrimSpace(s)
	return s, err
}

func commit() (string, error) {
	return shell("git", "rev-parse", "--short", "HEAD")
}

func branch() (string, error) {
	return shell("git", "rev-parse", "--abbrev-ref", "HEAD")
}

func split(in interface{}, c string, numToken int) string {
	s := fmt.Sprintf("%v", in)

	tokens := strings.Split(s, c)

	if numToken != -1 {
		return tokens[numToken]
	} else {
		return tokens[len(tokens)-1]
	}

}

func generateFile(templatePath, destPath string) bool {
	tmpl := template.New(filepath.Base(templatePath)).Funcs(template_funcs)

	if len(delims) > 0 {
		tmpl = tmpl.Delims(delims[0], delims[1])
	}
	tmpl, err := tmpl.ParseFiles(templatePath)
	if err != nil {
		log.Fatalf("unable to parse template: %s", err)
	}

	dest := os.Stdout
	if destPath != "" {
		dest, err = os.Create(destPath)
		if err != nil {
			log.Fatalf("unable to create %s", err)
		}
		defer dest.Close()
	}

	//fmt.Println(work.Values)
	if err := mergo.Merge(&work.Values, work.Env()); err != nil {
		log.Fatalf("template error: %s\n", err)
	}
	//fmt.Println(work.Values)

	err = tmpl.ExecuteTemplate(dest, filepath.Base(templatePath), &work.Values)
	if err != nil {
		log.Fatalf("template error: %s\n", err)
	}

	if fi, err := os.Stat(destPath); err == nil {
		if err := dest.Chmod(fi.Mode()); err != nil {
			log.Fatalf("unable to chmod temp file: %s\n", err)
		}
		if err := dest.Chown(int(fi.Sys().(*syscall.Stat_t).Uid), int(fi.Sys().(*syscall.Stat_t).Gid)); err != nil {
			log.Fatalf("unable to chown temp file: %s\n", err)
		}
	}

	return true
}

func generateDir(templateDir, destDir string) bool {
	if destDir != "" {
		fiDest, err := os.Stat(destDir)
		if err != nil {
			log.Fatalf("unable to stat %s, error: %s", destDir, err)
		}
		if !fiDest.IsDir() {
			log.Fatalf("if template is a directory, dest must also be a directory (or stdout)")
		}
	}

	files, err := ioutil.ReadDir(templateDir)
	if err != nil {
		log.Fatalf("bad directory: %s, error: %s", templateDir, err)
	}

	for _, file := range files {
		if destDir == "" {
			generateFile(filepath.Join(templateDir, file.Name()), "")
		} else {
			generateFile(filepath.Join(templateDir, file.Name()), filepath.Join(destDir, file.Name()))
		}
	}

	return true
}

func executeTemplate(out io.Writer, tpl *template.Template) error {
	if err := mergo.Merge(&work.Values, work.Env()); err != nil {
		log.Fatalf("template error: %s\n", err)
	}

	err := tpl.Execute(out, &work.Values)
	if err != nil {
		return fmt.Errorf("Failed to parse standard input: %v", err)
	}
	return nil
}
