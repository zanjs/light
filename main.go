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
	version = flag.Bool("v", false, "version")
)

func main() {
	flag.Parse()

	if *version {
		fmt.Println("gobatis v0.2.0")
	}

	log.SetLevel(log.Lwarn)
	log.SetFormat("2006-01-02 15:04:05.999 info examples/main.go:88 message")

	gofile := os.Getenv("GOFILE")
	if gofile == "" || !strings.HasSuffix(gofile, ".go") {
		fmt.Println("used by go:generate only")
		return
	}

	pwd, err := os.Getwd()
	checkError(err)
	fmt.Printf("Found go file: %s/%s\n", pwd, gofile)

	filename := gofile[:len(gofile)-3] + "impl.go"

	// Check modified time, if generated file newer than source file, skip!
	// gofi, err := os.Stat(pwd + "/" + gofile)
	// checkError(err)
	//
	// fi, err := os.Stat(filename)
	// if err != nil {
	// 	if _, ok := err.(*os.PathError); !ok {
	// 		panic(err)
	// 	}
	// } else {
	// 	if gofi.ModTime().Before(fi.ModTime()) {
	// 		fmt.Printf("Generate file: %s/%s, skip!\n", pwd, filename)
	// 		return
	// 	}
	// }

	os.Remove(filename)

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

	// log.Infof("format source")
	// pretty, err := format.Source(buf.Bytes())
	// checkError(err)

	log.Infof("write source to file")
	// ioutil.WriteFile(filename, pretty, 0644)
	ioutil.WriteFile(filename, buf.Bytes(), 0644)

	log.Infof("format and import using goimports tool")
	cmd := exec.Command("goimports", "-w", pwd+"/"+filename)
	err = cmd.Run()
	checkError(err)

	fmt.Printf("Generate file: %s/%s\n", pwd, filename)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
