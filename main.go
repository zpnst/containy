package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/zpnst/containy/linux"
)

const (
	configyPath string = "bundle/configy.json"
)

func main() {
	f, err := os.Open(configyPath)
	if err != nil {
		log.Panicln(err)
	}
	defer f.Close()

	b, err := io.ReadAll(f)
	if err != nil {
		log.Panicln(err)
	}

	var configy linux.Configy
	if err := json.Unmarshal(b, &configy); err != nil {
		log.Panicln(err)
	}

	log.Printf("%+v", configy)
}
