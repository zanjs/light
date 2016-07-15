package main

import (
	"os"

	"github.com/wothing/log"
)

func main() {
	gofile := os.Getenv("GOFILE")
	log.Infof(gofile)

	mapper.Source = gofile

	filename := gofile[:len(gofile)-3] + "impl.go"
	os.Remove(filename)

	parseGofile(gofile)

	prepareData()

	log.JSONIndent(mapper)

	// tplName := "template.txt"
	// tpl, err := Asset(tplName)
	// checkError(err)
	//
	// t, err := template.New(tplName).Parse(string(tpl))
	// checkError(err)
	//
	// var buf bytes.Buffer
	// err = t.Execute(&buf, mapper)
	// checkError(err)
	//
	// pretty, err := format.Source(buf.Bytes())
	// checkError(err)
	//
	// ioutil.WriteFile(filename, pretty, 0644)
}

func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
