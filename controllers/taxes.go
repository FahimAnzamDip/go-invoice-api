package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexTaxHandler(w http.ResponseWriter, r *http.Request) {
	tax := &models.Tax{}

	data := tax.Index()
	u.Respond(w, data)
}

func StoreTaxHandler(w http.ResponseWriter, r *http.Request) {
	tax := &models.Tax{}

	err := json.NewDecoder(r.Body).Decode(tax)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := tax.Store()
	u.Respond(w, res)
}

func ShowTaxHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	tax := &models.Tax{}

	data := tax.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateTaxHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	tax := &models.Tax{}
	err = json.NewDecoder(r.Body).Decode(tax)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := tax.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyTaxHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	tax := &models.Tax{}

	res := tax.Destroy(uint(ID))
	u.Respond(w, res)
}
