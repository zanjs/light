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

	"github.com/wothing/log"
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
		fmt.Println("gobatis v0.2.4")
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

	filename := gofile[:len(gofile)-3] + "impl.go"

	if !*force {
		// Check modified time, if generated file newer than source file, skip!
		gofi, err := os.Stat(gofile)
		checkError(err)

		fi, err := os.Stat(filename)
		if err != nil {
			if _, ok := err.(*os.PathError); !ok {
				panic(err)
			}
		} else {
			if gofi.ModTime().Before(fi.ModTime()) {
				fmt.Printf("Generated file: %s, skip!\n", filename)
				return
			}
		}
	}

	// !!! go/types parse lib files in $GOPATH/pkg/..., so must build deps package
	log.Infof("build deps using `go build -i -v`")
	cmd := exec.Command("go", "build", "-i", "-v", "./"+gofile)
	out, err := cmd.CombinedOutput()
	fmt.Printf("%s", out[23:])
	checkError(err)

	os.Remove(filename)

	log.Infof("start parse go file")
	mapper.Source = gofile
	parseGofile(gofile)
	// log.JSONIndent(mapper)

	log.Infof("preparse data")
	prepareData()
	log.JSONIndent(mapper)

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

	// log.Infof("format source")
	// pretty, err := format.Source(buf.Bytes())
	// checkError(err)

	log.Infof("write source to file")
	// ioutil.WriteFile(filename, pretty, 0644)
	ioutil.WriteFile(filename, buf.Bytes(), 0644)

	fmt.Printf("Generated file: %s\n", filename)

	log.Infof("format and import using goimports tool")
	cmd = exec.Command("goimports", "-w", filename)
	out, err = cmd.CombinedOutput()
	fmt.Printf("%s\n", out)
	checkError(err)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
