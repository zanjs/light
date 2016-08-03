package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/gotips/log"
	"golang.org/x/tools/imports"
)

var (
	db      = flag.String("db", "db", "variable of prefix Query/QueryRow/Exec")
	path    = flag.String("path", "", "path variable db")
	force   = flag.Bool("force", false, "not skip, force to rewrite impl file even if it newer than go file")
	version = flag.Bool("v", false, "version")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: gobatis [flags] [file.go]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Println("gobatis v0.2.5")
		return
	}

	log.SetLevel(log.Lwarn)
	log.SetFormat("2006-01-02 15:04:05.999 info examples/main.go:88 message")

	gofile := os.Getenv("GOFILE")
	if gofile == "" {
		if flag.NArg() >= 1 {
			gofile = flag.Arg(0)
		} else {
			flag.Usage()
		}
	}

	if !strings.HasSuffix(gofile, ".go") {
		fmt.Println("file suffix must match *.go")
		return
	}

	fmt.Printf("Found  go file: %s\n", gofile)

	target := gofile[:len(gofile)-3] + "impl.go"

	if skip := checkSkip(gofile, target); skip {
		fmt.Printf("Generated file: %s, skip!\n", target)
		return
	}

	buildDeps(gofile)

	os.Remove(target)

	log.Infof("start parse go file")
	mapper.Source = gofile
	parseGofile(gofile)
	// log.JSONIndent(mapper)

	log.Infof("preparse data")
	prepareData()
	// log.JSONIndent(mapper)

	tplName := "template.txt"
	tpl, err := Asset(tplName)
	checkError(err)
	log.Infof("parse template")
	t, err := template.New(tplName).Parse(string(tpl))
	checkError(err)
	log.Infof("render with template")
	var buf bytes.Buffer
	err = t.Execute(&buf, mapper)
	checkError(err)

	pretty, err := imports.Process(target, buf.Bytes(), nil)
	checkError(err)
	err = ioutil.WriteFile(target, pretty, 0644)
	checkError(err)
	fmt.Printf("Generated file: %s\n", target)
}

func checkSkip(gofile, target string) bool {

	if !*force {
		// Check modified time, if generated file newer than source file, skip!
		gofi, err := os.Stat(gofile)
		checkError(err)

		fi, err := os.Stat(target)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				panic(err)
			}
		} else {
			if gofi.ModTime().Before(fi.ModTime()) {
				return true
			}
		}
	}
	return false
}

func buildDeps(gofile string) {
	// !!! go/types parse lib files in $GOPATH/pkg/..., so must build deps package
	log.Infof("build deps using `go build -i -v`")
	cmd := exec.Command("go", "build", "-i", "-v", "./"+gofile)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out[23:])
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
