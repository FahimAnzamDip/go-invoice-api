package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexPurposeHandler(w http.ResponseWriter, r *http.Request) {
	purpose := &models.Purpose{}

	data := purpose.Index()
	u.Respond(w, data)
}

func StorePurposeHandler(w http.ResponseWriter, r *http.Request) {
	purpose := &models.Purpose{}

	err := json.NewDecoder(r.Body).Decode(purpose)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := purpose.Store()
	u.Respond(w, res)
}

func ShowPurposeHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	purpose := &models.Purpose{}

	res := purpose.Show(uint(ID))
	u.Respond(w, res)
}

func UpdatePurposeHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	purpose := &models.Purpose{}
	err = json.NewDecoder(r.Body).Decode(purpose)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := purpose.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyPurposeHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	purpose := &models.Purpose{}

	res := purpose.Destroy(uint(ID))
	u.Respond(w, res)
}
