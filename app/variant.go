package app

import (
	vr "category/app/variant"
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
	Var = "var"
)

//Register the Variant method
func (app *App) RegisterVariantsMethods() error {
	//schema will create the schema if not exists
	_, err := schema.NewStore(schema.Store{DB: app.DB})
	if err != nil {
		return err
	}
	app.Router.HandleFunc("/category/product/variant", app.createVariant).Methods("POST")
	app.Router.HandleFunc("/category/product/variant/{id}", app.getVariant).Methods("GET")
	app.Router.HandleFunc("/category/product/variant/{id}", app.deleteVariant).Methods("DELETE")
	app.Router.HandleFunc("/category/product/variant/{id}", app.updateVariant).Methods("PUT")
	return nil
}

//Create variant ie : product id is compulsory
func (app *App) createVariant(resp http.ResponseWriter, req *http.Request) {
	var variant *vr.Variant
	//Decoding the request
	err := json.NewDecoder(req.Body).Decode(&variant)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusInternalServerError, resp)
		if err != nil {
			return
		}
		return
	}
	//Category nil check
	if variant == nil {
		resp, err = NewMessage("variant must not be empty", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	if variant.ProductId == "" {
		resp, err = NewMessage("product id required for creating variant", http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	variant.Id = helper.New(Var)
	//Check if product_id exist or not
	var id string
	sqlStatement := `SELECT shopalyst_product_v1.product.id FROM shopalyst_product_v1.product WHERE id = $1`
	row := app.DB.QueryRow(sqlStatement, variant.ProductId)
	err = row.Scan(&id)
	if err != nil {
		if err == sql.ErrNoRows {
			resp, err = NewMessage("invalid product id", http.StatusPreconditionFailed, resp)
			if err != nil {
				return
			}
		}
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	//Now create product
	if err := vr.InsertVariant(variant, app.DB); err != nil {
		resp, err = NewMessage(err.Error(), http.StatusPreconditionFailed, resp)
		if err != nil {
			return
		}
		return
	}

	resp, err = BindResponse(variant, resp, http.StatusOK)
	if err != nil {
		resp, err = NewMessage(err.Error(), http.StatusInternalServerError, resp)
		if err != nil {
			return
		}
		return
	}
}

//get variant
func (app *App) getVariant(res http.ResponseWriter, req *http.Request) {
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
	variant, err := vr.GetVariant(id, app.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			res, err = NewMessage("variant  not found", http.StatusPreconditionFailed, res)
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

	res, err = BindResponse(variant, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
}
//Update variant
func (app *App) updateVariant(res http.ResponseWriter, req *http.Request) {
	log.Println("/category/product/variant")
	id := mux.Vars(req)["id"]
	reqBody, err := ioutil.ReadAll(req.Body)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

	var variant vr.Variant
	if err := json.Unmarshal(reqBody, &variant); err != nil {
		res, err = NewMessage(err.Error(), http.StatusPreconditionFailed, res)
		if err != nil {
			return
		}
		return
	}
	variant.Id = id

	v, err := vr.GetVariant(id, app.DB)
	if err != nil {
		if err == sql.ErrNoRows {
			res, err = NewMessage("variant not found", http.StatusPreconditionFailed, res)
			if err != nil {
				return
			}
			return
		}
		return
	}

	if v.Id == "" {
		res, err = NewMessage("variant not found", http.StatusNotFound, res)
		if err != nil {
			return
		}
		return
	}

	err = vr.UpdateVariant(id, app.DB, &variant)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}

	res, err = BindResponse(variant, res, http.StatusOK)
	if err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
}
//delete variant
func (app *App) deleteVariant(res http.ResponseWriter, req *http.Request) {
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
	if err := vr.DeleteVariant(id, app.DB); err != nil {
		res, err = NewMessage(err.Error(), http.StatusInternalServerError, res)
		if err != nil {
			return
		}
		return
	}
}


