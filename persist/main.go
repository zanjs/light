package persist

import (
	"bytes"
	"io/ioutil"
	"os"
	"regexp"
	"text/template"

	"github.com/gotips/log"
)

var re = regexp.MustCompile(`\$\{.+?}`)

func Main() {
	gofile := os.Getenv("GOFILE")

	meta, err := parseFile(gofile)
	if err != nil {
		log.Fatal(err)
	}
	// log.JSONIndent(meta)

	data := prepare(meta)
	// log.JSONIndent(data)

	tplFile := "../../persist/persist.txt"
	t, err := template.ParseFiles(tplFile)
	if err != nil {
		log.Error(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		log.Error(err)
	}

	filename := gofile[:len(gofile)-3] + "impl.go"

	ioutil.WriteFile(filename, buf.Bytes(), 0644)

	// pretty, err := format.Source(buf.Bytes())
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// ioutil.WriteFile(filename, pretty, 0644)
}
