package app

import (
	DB "category/schema"
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

//This method will initialize the router nad activate all the end points
func (app *App) InitializeAndRun(user, password, dbname, host string, port int) error {
	//New connection will create a connection with the given value
	conn := DB.NewConnection(&DB.Info{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Dbname:   dbname,
	})
	//Make connection will create a new connection and check db is active or not
	app.DB = conn.MakeConnection()
	app.Router = mux.NewRouter()
	//LogRequestHandler(app.Router)
	//Register category methods

	//All methods of API are registered
	if err := app.RegisterCategoryMethods(); err != nil {
		return err
	}
	//Register products methods
	if err := app.RegisterProductMethods(); err != nil {
		return err
	}
	//Register Variants methods
	if err := app.RegisterVariantsMethods(); err != nil {
		return err
	}

	//If there is any other end point will called from Ui the this will throw an error
	app.Router.MethodNotAllowedHandler = MethodNotAllowedHandler()

	//serve at localhost
	err := http.ListenAndServe(":8080", app.Router)
	if err != nil {
		panic("not able to listen:8080")
	}

	log.Println("server started at localhost:8080")
	return nil
}

//Gives error if any handler comes and not registered
func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Method not allowed")
	})
}
