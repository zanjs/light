package persist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"text/template"

	"github.com/wothing/log"
)

var re = regexp.MustCompile(`\$\{.+?}`)

func Main() {
	gofile := os.Getenv("GOFILE")

	meta, err := parseFile(gofile)
	if err != nil {
		log.Fatal(err)
	}

	data := prepare(meta)
	js, _ := json.MarshalIndent(data, "", "    ")
	fmt.Printf("%s\n", js)

	persist, err := ioutil.ReadFile("../../persist/persist.tpl")
	if err != nil {
		log.Error(err)
	}

	filename := gofile[:len(gofile)-3] + "impl.go"
	out, err := os.Create(filename)
	if err != nil {
		log.Error(err)
	}
	t := template.Must(template.New("persist").Parse(string(persist)))
	err = t.Execute(out, data)
	if err != nil {
		log.Error(err)
	}
}
