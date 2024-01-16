package main

import (
	"log"

	"github.com/ffajarpratama/boiler-api/cmd/app"
)

func main() {
	err := app.Exec()
	if err != nil {
		log.Fatalf("[app-run-failed] \n%v\n", err)
	}
}
