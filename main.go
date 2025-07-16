package main

import (
	"log"

	"github.com/zpnst/containy/linux"
)

const (
	configyPath string = "configy.json"
)

func main() {
	configy, err := linux.ParseConfigy(configyPath)
	if err != nil {
		log.Panicln(err)
	}
	log.Printf("%+v", configy)
}
