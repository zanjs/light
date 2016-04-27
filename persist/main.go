package persist

import (
	"io/ioutil"
	"os"
	"regexp"
	"text/template"

	"github.com/wothing/log"
)

var re = regexp.MustCompile(`\$\{.+?}`)

func Main() {
	gofile := os.Getenv("GOFILE")

	meta := parse(gofile)
	meta.Explain()

	// js, _ := json.MarshalIndent(meta, "", "    ")
	// fmt.Printf("%s\n", js)

	persist, err := ioutil.ReadFile("../../persist/persist.tpl")
	if err != nil {
		log.Error(err)
	}

	filename := gofile[:len(gofile)-3] + "data.go"
	out, err := os.Create(filename)
	if err != nil {
		log.Error(err)
	}
	t := template.Must(template.New("demo").Parse(string(persist)))
	err = t.Execute(out, meta)
	if err != nil {
		log.Error(err)
	}
}
