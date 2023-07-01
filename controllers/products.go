package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexProductHandler(w http.ResponseWriter, r *http.Request) {
	product := &models.Product{}

	data := product.Index()
	u.Respond(w, data)
}

func StoreProductHandler(w http.ResponseWriter, r *http.Request) {
	product := &models.Product{}

	err := json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := product.Store()
	u.Respond(w, res)
}

func ShowProductHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	product := &models.Product{}

	data := product.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateProductHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	product := &models.Product{}
	err = json.NewDecoder(r.Body).Decode(product)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := product.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyProductHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	product := &models.Product{}

	res := product.Destroy(uint(ID))
	u.Respond(w, res)
}
