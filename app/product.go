package app

import (
	products "category/app/product"
	"category/helper"
	schema "category/schema"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
)

const (
	Products = "pro"
)

func (app *App) RegisterProductMethods() {
	app.Router.HandleFunc("/category/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/category/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/category/product/{id}", app.deleteProduct).Methods("DELETE")
	app.Router.HandleFunc("/category/product/{id}", app.updateProduct).Methods("PUT")
}

func (app *App) createProduct(resp http.ResponseWriter, req *http.Request) {
	//schema will create the schema if not exists
	schema.NewStore(schema.Store{DB: app.DB})
	var product *products.ProductResp
	//Decoding the request
	err := json.NewDecoder(req.Body).Decode(&product)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusInternalServerError, resp)
		if err != nil {
			return
		}
		return
	}
	//Category nil check
	if product == nil {
		resp, err = NewMessage("product must not be empty", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	if product.CategoryId == "" {
		resp, err = NewMessage("category id required for creating product", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}
	product.Id = helper.New("pro")
	//Check if category_id exist or not
	var id string
	sqlStatement := `SELECT shopalyst_category_v1.category.id FROM shopalyst_category_v1.category WHERE id = $1`
	row := app.DB.QueryRow(sqlStatement, product.CategoryId)
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			resp, err = NewMessage("invalid category id", http.StatusPreconditionFailed, resp)
			if err != nil {
				return
			}
			return
		}
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
	}

	//Now create product
	if err := products.InsertProduct(product, app.DB); err != nil {
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}
}

func (app *App) getProduct(res http.ResponseWriter, req *http.Request) {
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

	prod, err := products.GetProductAndItsVariant(id, app.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			res, err = NewMessage("Product not found", http.StatusPreconditionFailed, res)
			if err != nil {
				return
			}
			return
		}
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}

		return
	}

	res, err = BindResponse(prod, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

}

func (app *App) deleteProduct(res http.ResponseWriter, req *http.Request) {
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
	if err := products.DeleteProduct(id, app.DB); err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
}

func (app *App) updateProduct(res http.ResponseWriter, req *http.Request) {
	log.Println("/category/product")
	id := mux.Vars(req)["id"]
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

	var product products.ProductResp
	if err := json.Unmarshal(reqBody, &product); err != nil {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}
	product.Id = id
	err = products.UpdateProduct(id, app.DB, &product)
	if err != nil {
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
