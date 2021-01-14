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

//Register Method will activate the mux
func (app *App) RegisterCategoryMethods() error {
	_, err := schema.NewStore(schema.Store{DB: app.DB})
	if err != nil {
		return err
	}
	app.Router.HandleFunc("/category/{id}", app.getCategory).Methods("GET")
	app.Router.HandleFunc("/category/{id}", app.deleteCategory).Methods("DELETE")
	app.Router.HandleFunc("/category/{id}", app.updateCategory).Methods("PUT")
	app.Router.HandleFunc("/category", app.createCategory).Methods("POST")
	return nil

}

//Create category will create the category
// if there is nested child. i.e : categories in side sub categories -> products -> variants
// It will insert categories as well as sub categories and products and its variants also
func (app *App) createCategory(resp http.ResponseWriter, req *http.Request) {
	log.Println("/category")
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

	//Insert the categories
	err = category.InsertCategories(cat, app.DB)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	//Insert Nested products and its variant if there is nesting exist
	if err := category.InsertProductsAndVariant(cat, app.DB); err != nil {
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	//Binding the response
	resp, err = BindResponse(cat, resp, http.StatusOK)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusInternalServerError, resp)
		if err != nil {
			return
		}
		return
	}

}

//get category will get the category, products and its variants.
//if there is any subcategories it will also give products and variant also
func (app *App) getCategory(res http.ResponseWriter, req *http.Request) {
	//getting id from the url
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

	//get all the categories , subcategories , products and its variants
	err, categories := category.GetCategories(id, app.DB, 0)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}

	//if category does not exist or deleted
	if categories.Id == "" {
		res, err = NewMessage("category not found", http.StatusNotFound, res)
		if err != nil {
			return
		}
		return
	}

	//Binding response
	res, err = BindResponse(categories, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

}

//Update category to update the category node.
// name of the category
func (app *App) updateCategory(res http.ResponseWriter, req *http.Request) {
	log.Println("/category")
	//getting id
	id := mux.Vars(req)["id"]
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
	}

	//marshal
	var categories category.CategoryResp
	if err != json.Unmarshal(reqBody, &categories) {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
	}

	//update the category. ie name of the category
	c, err := category.UpdateCategory(id, categories.Name, app.DB)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

	//binding the response
	res, err = BindResponse(c, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
}

//delete category will delete the category
func (app *App) deleteCategory(res http.ResponseWriter, req *http.Request) {
	//id getting
	params := mux.Vars(req)
	id := params["id"]
	var err error = nil
	_, err = schema.NewStore(schema.Store{DB: app.DB})
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
	//id validation
	if id == "" {
		res, err = NewMessage("id required", http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}
	//delete the category and its node
	if err := category.DeleteCategory(id, app.DB); err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
	//binding the response
	res, err = BindResponse(nil, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

}
