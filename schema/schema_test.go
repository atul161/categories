package connection_test

import (
	"category/app"
	connection "category/schema"
	"testing"
)

func TestNewStore(t *testing.T) {
	a := app.App{}
	a.InitializeAndRun("postgres", "postgres", "shopalyst", "localhost", 5432)
	//New store will create the schema and table
	connection.NewStore(connection.Store{DB: a.DB})
}
