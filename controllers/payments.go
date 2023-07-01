package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
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
