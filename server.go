package main

import (
	app "category/app"
)

func main() {
	a := app.App{}
	a.InitializeAndRun("postgres", "postgres", "shopalyst", "localhost", 5432)
}
