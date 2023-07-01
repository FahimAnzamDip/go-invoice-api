package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	invoice := &models.Invoice{}

	data := invoice.Index()
	u.Respond(w, data)
}

func StoreInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	invoice := &models.Invoice{}

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		log.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := invoice.Store()
	u.Respond(w, res)
}

func ShowInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	invoice := &models.Invoice{}

	data := invoice.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	invoice := &models.Invoice{}
	err = json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := invoice.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	invoice := &models.Invoice{}

	data := invoice.Destroy(uint(ID))
	u.Respond(w, data)
}
