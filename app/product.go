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

func (app *App) RegisterProductMethods() error {
	//if store does not exist then it will create a new store
	_, err := schema.NewStore(schema.Store{DB: app.DB})
	if err != nil {
		return err
	}
	app.Router.HandleFunc("/category/product", app.createProduct).Methods("POST")
	app.Router.HandleFunc("/category/product/{id}", app.getProduct).Methods("GET")
	app.Router.HandleFunc("/category/product/{id}", app.deleteProduct).Methods("DELETE")
	app.Router.HandleFunc("/category/product/{id}", app.updateProduct).Methods("PUT")
	return nil
}


//Create product. ie - category id is compulsory for creating product
//if there is nested child / variant the this endpoint will create variants also.

func (app *App) createProduct(resp http.ResponseWriter, req *http.Request) {
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
	//product nil check
	if product == nil {
		resp, err = NewMessage("product must not be empty", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	//category id validation
	if product.CategoryId == "" {
		resp, err = NewMessage("category id required for creating product", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}
	//initialising with new id
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

	//Bind Response
	resp, err = BindResponse(product, resp, http.StatusOK)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusInternalServerError, resp)
		if err != nil {
			return
		}
		return
	}

}


//Get Product - need id in request and it will return product and its variant
func (app *App) getProduct(res http.ResponseWriter, req *http.Request) {
	//getting if
	params := mux.Vars(req)
	id := params["id"]
	var err error = nil
	if id == "" {
		res, err = NewMessage("id required for getting product", http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}

	prod, err := products.GetProducts(id, app.DB, 0)
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

	if prod.Id == "" {
		res, err = NewMessage("product not found", http.StatusNotFound, res)
		if err != nil {
			return
		}
		return
	}
	//Binding Response
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
	//Before update we need to check whether id exist or not
	p, err := products.GetProducts(id, app.DB, 0)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

	if p.Id == "" {
		res, err = NewMessage("product not found", http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

	err = products.UpdateProduct(id, app.DB, &product)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

	res, err = BindResponse(product, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
}
