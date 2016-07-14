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

	err := parseGofile(gofile)
	if err != nil {
		panic(err)
	}

	prepareData()

	log.JSONIndent(mapper)
}
