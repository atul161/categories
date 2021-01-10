package app

import (
	category "category/app/category"
	"category/helper"
	schema "category/schema"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Categories = "cat"
	Limit      = 10
)

func (app *App) RegisterCategoryMethods() {
	app.Router.HandleFunc("/category/{id}", app.getCategory).Methods("GET")
	app.Router.HandleFunc("/category/{id}", app.deleteCategory).Methods("DELETE")
	app.Router.HandleFunc("/category/{id}", app.updateCategory).Methods("PUT")
	app.Router.HandleFunc("/category",      app.createCategory).Methods("POST")

}

func (app *App) createCategory(resp http.ResponseWriter, req *http.Request) {
	log.Println("/category")
	//schema will create the schema if not exists
	schema.NewStore(schema.Store{DB: app.DB})
	var cat *category.CategoryResp
	//Decoding the request
	err := json.NewDecoder(req.Body).Decode(&cat)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusInternalServerError, resp)
		if err != nil {
			return
		}
		return
	}
	//Category nil check
	if cat == nil {
		resp, err = NewMessage("category must not be empty", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}
	// Subcategories id mapping
	err = MapWithId(cat, 0)
	if err != nil {
		resp, err = NewMessage("sub categories nesting limit exceeded", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}
	//Products id mapping and variant mapping
	for _, prod := range cat.Products {
		prod.Id = helper.New(Products)
		for _, variant := range prod.Variants {
			variant.Id = helper.New(Var)
		}
	}

	err = category.InsertCategories(cat, app.DB)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	if err := category.InsertProductsAndVariant(cat, app.DB); err != nil {
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	resp, err = BindResponse(cat, resp, http.StatusOK)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusInternalServerError, resp)
		if err != nil {
			return
		}
		return
	}

}

func (app *App) getCategory(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	var err error = nil
	if id == "" {
		res, err = NewMessage("id required for getting hierarchy information", http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}

	err, categories := category.GetCategory(app.DB, id, 0, nil, true)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}

	categories, err = category.GetProductAndItsVariant(categories, app.DB)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}

	res, err = BindResponse(categories, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

}

func (app *App) updateCategory(res http.ResponseWriter, req *http.Request) {
	log.Println("/category")
	id := mux.Vars(req)["id"]
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res, _ = BindResponse(nil, res, http.StatusPreconditionFailed)
	}
	var categories category.CategoryResp
	if err != json.Unmarshal(reqBody, &categories) {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
	}

	c, err := category.UpdateCategory(id, categories.Name, app.DB)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
	res, err = BindResponse(c, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
}

func (app *App) deleteCategory(res http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)
	id := params["id"]
	var err error = nil
	if id == "" {
		res, err = NewMessage("id required", http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}
	if err := category.DeleteCategory(id, app.DB); err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
	res, err = BindResponse(nil, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

}
