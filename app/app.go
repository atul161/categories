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

//Create a new connection and returns thee database instance
func (app *App) InitializeAndRun(user, password, dbname, host string, port int) {
	conn := DB.NewConnection(&DB.Info{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Dbname:   dbname,
	})
	app.DB = conn.MakeConnection()
	app.Router = mux.NewRouter()
	//LogRequestHandler(app.Router)
	//Register category methods
	//All methods of API are registered
	app.RegisterCategoryMethods()
	app.RegisterProductMethods()
	app.RegisterVariantsMethods()
	app.Router.MethodNotAllowedHandler = MethodNotAllowedHandler()

	err := http.ListenAndServe(":8080", app.Router)
	if err != nil {
		panic("not able to listen:8080")
	}

	log.Println("server started at localhost:8080")

}

func MethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Method not allowed")
	})
}
