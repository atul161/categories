package main

import (
	"category/app"
	"log"
	"os"
	"strconv"
)


func main() {
	var (
		user            = os.Getenv("user")
		password        = os.Getenv("password")
		dbname          = os.Getenv("dbname")
		host            = os.Getenv("host")
		port            = os.Getenv("port")
	)

	p, err := strconv.Atoi(port)
	if err != nil {
		panic(err)
	}
	a := app.App{}
	//server started at localhost
	err = a.InitializeAndRun(user, password, dbname, host, p)
	if err != nil {
		log.Println(err.Error())
	}

}
