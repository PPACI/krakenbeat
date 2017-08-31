package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/PPACI/krakenbeat/beater"
)

func main() {
	err := beat.Run("krakenbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
