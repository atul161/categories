package main

import (
	app "category/app"
	"log"
)

func main() {
	a := app.App{}
	err := a.InitializeAndRun("postgres", "postgres", "shopalyst", "localhost", 5432)
	if err != nil {
		log.Println(err.Error())
	}

}
