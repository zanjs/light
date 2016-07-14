package main

import (
	"os"

	"github.com/wothing/log"
)

func main() {
	gofile := os.Getenv("GOFILE")
	log.Infof(gofile)

	filename := gofile[:len(gofile)-3] + "impl.go"
	os.Remove(filename)

	err := readGoFile(gofile)
	if err != nil {
		panic(err)
	}
	log.JSONIndent(m)

	// log.JSONIndent(uses)
}
