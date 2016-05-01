package persist

import (
	"bytes"
	"go/format"
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
	// mjs, _ := json.MarshalIndent(meta, "", "    ")
	// log.Debugf("%s\n", mjs)

	data := prepare(meta)
	// js, _ := json.MarshalIndent(data, "", "    ")
	// log.Debugf("%s\n", js)

	tplFile := "../../persist/persist.tpl"
	t, err := template.ParseFiles(tplFile)
	if err != nil {
		log.Error(err)
	}
	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		log.Error(err)
	}

	pretty, err := format.Source(buf.Bytes())
	if err != nil {
		log.Fatal(err)
	}

	filename := gofile[:len(gofile)-3] + "impl.go"

	ioutil.WriteFile(filename, pretty, 0644)
	if err != nil {
		log.Fatal(err)
	}
}
