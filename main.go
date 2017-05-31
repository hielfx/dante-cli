package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"golang.org/x/net/context"
)

const defaultWaitRetryInterval = time.Second

type sliceVar []string
type hostFlagsVar []string

var logger = log.New(os.Stderr, "", 0)

type Context struct {
	Values map[string]interface{}
}

type HttpHeader struct {
	name  string
	value string
}

func (c *Context) Env() map[string]interface{} {
	env := make(map[string]interface{})
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}

var (
	buildVersion string
	version      bool
	poll         bool
	wg           sync.WaitGroup

	templatesFlag    sliceVar
	templateDirsFlag sliceVar

	valuesFlag sliceVar
	delimsFlag string
	delims     []string

	ctx    context.Context
	cancel context.CancelFunc

	work Context
)

func (i *hostFlagsVar) String() string {
	return fmt.Sprint(*i)
}

func (i *hostFlagsVar) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (s *sliceVar) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func (s *sliceVar) String() string {
	return strings.Join(*s, ",")
}

func usage() {
	println(`Usage: dante-cli [options] [command]

A simple CLI templating tool written in golang

Options:`)
	flag.PrintDefaults()

	println(`
Arguments:
  command - command to be executed
  `)

	println(`Examples:
`)
	println(`   Generate deployment.yml from k8b/deployment.yml`)
	println(`
   dante-cli -template k8b/deployment.yml:deployment.yml \
              --values value.yml
	`)

	println(`For more information, see https://github.com/jhidalgo3/dante-cli`)
}

func main() {

	flag.BoolVar(&version, "version", false, "show version")

	flag.Var(&templatesFlag, "template", "Template (/template:/dest). Can be passed multiple times. Does also support directories")
	flag.StringVar(&delimsFlag, "delims", "", `template tag delimiters. default "{{":"}}" `)

	flag.Var(&valuesFlag, "values", ` Values to template `)

	flag.Usage = usage
	flag.Parse()

	if version {
		fmt.Println(buildVersion)
		return
	}

	/*if flag.NArg() == 0 && flag.NFlag() == 0 {
		usage()
		os.Exit(1)
	}*/

	if delimsFlag != "" {
		delims = strings.Split(delimsFlag, ":")
		if len(delims) != 2 {
			log.Fatalf("bad delimiters argument: %s. expected \"left:right\"", delimsFlag)
		}
	}

	var files []string
	for _, t := range valuesFlag {

		files = append(files, t)
	}
	work = Context{}
	//work.Values, _ = loadValues(files)

	work.Values, _ = calculateValues(files)

	log.Println(work.Values)

	if len(templatesFlag) == 0 {
		tpl, err := loadFileOrStdin("")
		if err != nil {
			logError("Error occurred while loading data:", err)
		}

		err = executeTemplate(os.Stdout, tpl)
		if err != nil {
			logError("Error occurred while attempting to template:", err)
		}
	} else {
		for _, t := range templatesFlag {
			template, dest := t, ""
			if strings.Contains(t, ":") {
				parts := strings.Split(t, ":")
				if len(parts) != 2 {
					log.Fatalf("bad template argument: %s. expected \"/template:/dest\"", t)
				}
				template, dest = parts[0], parts[1]
			}

			fi, err := os.Stat(template)
			if err != nil {
				log.Fatalf("unable to stat %s, error: %s", template, err)
			}
			if fi.IsDir() {
				generateDir(template, dest)
			} else {
				generateFile(template, dest)
			}
		}
	}

}

func logError(msg string, err error) {
	logger.Println(msg)
	logger.Println(err.Error())
}
