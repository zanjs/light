package main

import (
	"os"

	"github.com/wothing/log"
)

func main() {
	gofile := os.Getenv("GOFILE")

	filename := gofile[:len(gofile)-3] + "impl.go"
	os.Remove(filename)

	meta, err := parseFile(gofile)
	if err != nil {
		panic(err)
	}
	log.JSONIndent(meta)

}
