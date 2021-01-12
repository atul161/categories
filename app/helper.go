package app

import (
	category "category/app/category"
	"category/helper"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

//Putting the Id inside the Nested Child
/*
  c - c1 , c2 , c3 , list of p
     c1 -
*/
func MapWithId(cat *category.CategoryResp, layer int) error {
	if cat == nil {
		return nil
	}
	cat.Id = helper.New(Categories)
	for _, product := range cat.Products {
		if product == nil {
			return nil
		}
		product.Id = helper.New(Products)
		for _, variant := range product.Variants {
			if variant == nil {
				return nil
			}
			variant.Id = helper.New(Var)
		}
	}
	for _, v := range cat.ChildCategories {
		if v == nil {
			return nil
		}
		if layer == 10 {
			return errors.New("max nesting limit is upto 10")
		}
		layer = layer + 1
		err := MapWithId(v, layer)
		if err != nil {
			return err
		}
		layer--
	}

	return nil
}

type Logger struct {
	Uri    string `json:"uri"`
	Method string `json:"method"`
	Body   string `json:"body"`
}

func LogRequestHandler(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		// call the original http.Handler we're wrapping
		h.ServeHTTP(w, r)
		b, _ := ioutil.ReadAll(r.Body)
		// gather information about request and log it
		l := Logger{
			Uri:    r.URL.String(),
			Method: r.Method,
			Body:   string(b),
		}
		log.Println(l)

	}

	// http.HandlerFunc wraps a function so that it
	// implements http.Handler interface
	return http.HandlerFunc(fn)
}
