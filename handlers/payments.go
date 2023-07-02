package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexPaymentHandler(w http.ResponseWriter, r *http.Request) {
	payment := &models.Payment{}

	data := payment.Index()
	u.Respond(w, data)
}

func StorePaymentHandler(w http.ResponseWriter, r *http.Request) {
	payment := &models.Payment{}

	err := json.NewDecoder(r.Body).Decode(payment)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := payment.Store()
	u.Respond(w, res)
}

func ShowPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	payment := &models.Payment{}

	data := payment.Show(uint(ID))
	u.Respond(w, data)
}

func UpdatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	payment := &models.Payment{}
	err = json.NewDecoder(r.Body).Decode(payment)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := payment.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	payment := &models.Payment{}

	res := payment.Destroy(uint(ID))
	u.Respond(w, res)
}
